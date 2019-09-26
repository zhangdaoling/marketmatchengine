#!/usr/bin/env bash

/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic order_usdt-btc
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic transaction_usdt-btc
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic quotation_usdt-btc
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic cancel_usdt-btc
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --create --topic "order_usdt-btc";
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --create --topic "transaction_usdt-btc";
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --create --topic "quotation_usdt-btc";
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --create --topic "cancel_usdt-btc";
