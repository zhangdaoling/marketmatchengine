# marketmatchengine
market match engine， golang;

# dep
* need golang 1.12, go module，kafka
* git clone git@github.com:zhangdaoling/marketmatchengine.git
* cd path
* go get github.com/stretchr/testify/assert

# example
see "TestEngine" in engine/engine_test.go

0.0.0.0         account.jetbrains.com


# to do
*order按照数据库created先后，进入kafka，kafka充当定序功能

*order进kafka前需要保证唯一性（是进数据时保证，还是按照kafka offset保证），撮合内部没有考虑去重

*test: kafka_demo simple test: pass

*kafka message create timestamp没法设置（研究下kafka这个特性）

*test:quotation array test:pass

*test:order match test: pass

*test:market order test

*test:cancel order dev and test

*more order type: cancel type

*log

*行情服务。目前行情直接推送进入kafka

*清算服务。撮合结果根据资产类型进入kafka patatition。消费者按照手续费、资产类型分别清算，减少锁数据库

*通知服务。
