# marketmatchengine
marketmatch engine for golang;

# dep
* need golang 1.12, go module
* git clone git@github.com:zhangdaoling/marketmatchengine.git
* cd path
* go get github.com/stretchr/testify/assert

# example
see "TestEngine" in engine/engine_test.go


0.0.0.0         account.jetbrains.com


# to do
*order按照数据库created先后，进入kafka，如果时间一样的订单进来
*order进kafka前需要保证唯一性（是进数据时保证，还是按照kafka offset保证），需不需要保证自增性，撮合暂时没有考虑去重
*test: kafka_demo simple test: pass
*kafka message create timestamp没法设置
*test:quotation array test:pass
*test:order match test: pass
*test:market order test
*test:cancel order dev and test
*more order type


*log

