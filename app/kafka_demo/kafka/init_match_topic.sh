#!/usr/bin/env bash

/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic order_A-B
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic transaction_A-B
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic quotation_A-B
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic cancel_A-B
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --create --topic "order_A-B";
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --create --topic "transaction_A-B";
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --create --topic "quotation_A-B";
/Users/zhangdaoling/work/soft/kafka_2.12-2.3.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --create --topic "cancel_A-B";
