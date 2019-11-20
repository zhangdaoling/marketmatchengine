package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/zhangdaoling/marketmatchengine/engine"
	"github.com/zhangdaoling/marketmatchengine/order"

	"github.com/Shopify/sarama"
)

type Job interface {
	Loop(shutdown chan struct{})
}

var Symbol = "A-B"
var OrderTopic = "order_" + "A-B"
var TransactionTopic = "transaction_" + "A-B"
var QuotationTopic = "quotation_" + "A-B"

func main() {
	kafkaBrokers := flag.String("order_kafka", "localhost:9092", "The Kafka brokers to connect to, as a comma separated list")
	persistPath := flag.String("persist_path", "/tmp/", "engine persist path")
	flag.Parse()
	brokers := strings.Split(*kafkaBrokers, ",")

	match, err := engine.NewEngine(Symbol, 0, 0, 0)
	if err != nil {
		log.Fatalln(err)
	}
	orderConsumer, err := NewKafkaConsumer(brokers, OrderTopic, Symbol, 0, 100, int64(match.LastIndex))
	if err != nil {
		log.Fatalln(err)
	}
	transactionProducer, err := NewKafkaTransactionProducer(brokers, TransactionTopic, Symbol, 100)
	if err != nil {
		log.Fatalln(err)
	}
	quotationProducer, err := NewKafkaQuotationProducer(brokers, QuotationTopic, Symbol, 100)
	if err != nil {
		log.Fatalln(err)
	}
	app := &App{
		PersistPath:     *persistPath,
		Engine:          match,
		OrderChan:       orderConsumer.Channel,
		TransactionChan: transactionProducer.Channel,
		QuotationChan:   quotationProducer.Channel,
	}
	doLoopJobs(app, orderConsumer, transactionProducer, quotationProducer)
}

func doLoopJobs(jobs ...Job) {
	shutdown := make(chan struct{})
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		for sig := range signals {
			if sig == os.Interrupt || sig == syscall.SIGTERM {
				log.Printf("received signal [%v], preparing to quit", sig)
				close(shutdown)
			} else if sig == syscall.SIGHUP {
				log.Printf("received signal [%v], ignored", sig)
			}
		}
	}()

	var wg sync.WaitGroup

	for index := range jobs {
		wg.Add(1)
		job := jobs[index]
		go func() {
			defer wg.Done()
			job.Loop(shutdown)
		}()
	}

	wg.Wait()
}

type App struct {
	PersistPath     string
	Engine          *engine.Engine
	OrderChan       chan *sarama.ConsumerMessage
	TransactionChan chan []*order.Transaction
	QuotationChan   chan *order.OrderBook
}

func (a *App) Loop(shutdown chan struct{}) {
	var err error
	quotationTime := time.NewTicker(5 * time.Second)
	persitTime := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-shutdown:
			fmt.Println("close app")
			_, _, err = a.Engine.Persist(a.PersistPath)
			if err != nil {
				log.Println(err)
			}
			return

		case msg := <-a.OrderChan:
			o := &order.Order{}
			err = json.Unmarshal(msg.Value, o)
			if err != nil {
				log.Printf("json marshal error: %v\n", err)
				continue
			}
			if o.Symbol != a.Engine.Symbol {
				log.Printf("symbol error: %s!=%s\n", o.Symbol, a.Engine.Symbol)
				continue
			}
			o.Index = uint64(msg.Offset)
			//log.Println(msg.Timestamp)
			//log.Println(msg.Timestamp.Unix())
			o.IndexTime = uint64(msg.Timestamp.Unix())
			result, isNext, err := a.Engine.Match(o)
			if err != nil {
				log.Println(err)
				close(shutdown)
				return
			}
			if !isNext {
				close(shutdown)
			}
			if len(result) != 0 {
				//log.Println(result)
				a.TransactionChan <- result
			}

		case <-quotationTime.C:
			info := a.Engine.Quotation()
			a.QuotationChan <- info

		case <-persitTime.C:
			fileName, size, err := a.Engine.Persist(a.PersistPath)
			if err != nil {
				log.Println(err)
			}
			log.Println("persist:", fileName, size)
		}
	}
}

type KafkaConsumer struct {
	KafkaBrokers      []string
	Topic             string
	Symbol            string
	Key               string
	Consumer          sarama.Consumer
	PartitionConsumer sarama.PartitionConsumer
	Channel           chan *sarama.ConsumerMessage
}

func NewKafkaConsumer(brokers []string, topic string, symbol string, partition int32, size int, offset int64) (consumer *KafkaConsumer, err error) {
	c, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}
	partitionConsumer, err := c.ConsumePartition(topic, partition, int64(offset))
	if err != nil {
		return nil, err
	}
	consumer = &KafkaConsumer{
		KafkaBrokers:      brokers,
		Topic:             topic,
		Symbol:            symbol,
		Key:               topic,
		Consumer:          c,
		PartitionConsumer: partitionConsumer,
		Channel:           make(chan *sarama.ConsumerMessage, size),
	}
	return
}

func (c *KafkaConsumer) Loop(shutdown chan struct{}) {
	var err error
	for {
		select {
		case <-shutdown:
			fmt.Println("close OrderConsumer")
			if err = c.Consumer.Close(); err != nil {
				log.Println(err)
			}
			if err = c.PartitionConsumer.Close(); err != nil {
				log.Println(err)
			}
			return
		case msg := <-c.PartitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)

			c.Channel <- msg
		}
	}
}

type KafkaTransactionProducer struct {
	KafkaBrokers []string
	Topic        string
	Symbol       string
	Key          string
	Producer     sarama.SyncProducer
	Channel      chan []*order.Transaction
}

func NewKafkaTransactionProducer(brokers []string, topic string, symbol string, size int) (producer *KafkaTransactionProducer, err error) {
	p, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}
	producer = &KafkaTransactionProducer{
		KafkaBrokers: brokers,
		Topic:        topic,
		Symbol:       symbol,
		Key:          topic,
		Producer:     p,
		Channel:      make(chan []*order.Transaction, size),
	}
	return
}

func (p *KafkaTransactionProducer) Loop(shutdown chan struct{}) {
	var b []byte
	var err error
	for {
		select {
		case <-shutdown:
			fmt.Println("close TransactionProducer")
			if err = p.Producer.Close(); err != nil {
				log.Println(err)
			}
			return
		case o := <-p.Channel:
			b, err = json.Marshal(o)
			if err != nil {
				log.Printf("json marshal error: %v\n", err)
				close(shutdown)
				return
			}
			msg := &sarama.ProducerMessage{
				Topic: p.Topic,
				//Key: sarama.StringEncoder(p.Key),
				Value: sarama.ByteEncoder(b),
			}
			partition, offset, err := p.Producer.SendMessage(msg)
			if err != nil {
				log.Printf("FAILED to send message: %s\n", err)
				close(shutdown)
				return
			} else {
				log.Printf("transaction message sent to partition %d at offset %d\n", partition, offset)
			}
		}
	}
}

type KafkaQuotationProducer struct {
	KafkaBrokers []string
	Topic        string
	Symbol       string
	Key          string
	Producer     sarama.SyncProducer
	Channel      chan *order.OrderBook
}

func NewKafkaQuotationProducer(brokers []string, topic string, symbol string, size int) (producer *KafkaQuotationProducer, err error) {
	p, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}
	producer = &KafkaQuotationProducer{
		KafkaBrokers: brokers,
		Topic:        topic,
		Symbol:       symbol,
		Key:          topic,
		Producer:     p,
		Channel:      make(chan *order.OrderBook, size),
	}
	return
}

func (p *KafkaQuotationProducer) Loop(shutdown chan struct{}) {
	var b []byte
	var err error
	for {
		select {
		case <-shutdown:
			fmt.Println("close QuotationProducer")
			if err = p.Producer.Close(); err != nil {
				log.Println(err)
			}
			return
		case o := <-p.Channel:
			b, err = json.Marshal(o)
			if err != nil {
				log.Printf("json marshal error: %v\n", err)
				close(shutdown)
				return
			}
			msg := &sarama.ProducerMessage{
				Topic: p.Topic,
				//Key: sarama.StringEncoder(p.Key),
				Value: sarama.ByteEncoder(b),
			}
			partition, offset, err := p.Producer.SendMessage(msg)
			if err != nil {
				log.Printf("FAILED to send message: %s\n", err)
				close(shutdown)
				return
			} else {
				log.Printf("quotation message sent to partition %d at offset %d\n", partition, offset)
			}
		}
	}
}
