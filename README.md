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
*行情实现不太好，使用btree（google/btree）替换proprity_queue，行情改为主动推送。

*order进kafka前需要保证唯一性，撮合内部不做去重

*kafka message create timestamp没法设置（研究下kafka这个特性）

*支持取消交易功能

*log

*行情服务:数据流转:撮合行情数据->kafka->websocket。websocket根据时间去重

*清算服务:数据流:撮合交易数据->kafka->消费者->db(mysql)，用户资产按照资产类型分库分表，所有数据库操作按照uid排序再操作数据库，交易数据会被多个消费者消费，例如btc-usdt交易对，会被btc，usdt，手续费三个消费者消费。

*通知服务:数据流:撮合交易数据->kafka->通知服务

# support
*btc address：394nLFQo2XVf9ruET6JYRw4inoPF2YUaox

*eth address：0x4914eAb996a15c8b8B0896F178357838cb0aD60a
