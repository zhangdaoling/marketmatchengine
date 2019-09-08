package queue

import "github.com/zhangdaoling/marketmatchengine/common"

type Item interface {
	Compare(other interface{}) int
	Key() uint32
	Serialize() (zero *common.ZeroCopySink)
}

type PriorityQueue interface {
	Insert(item Item) (i Item)
	Cancel(key uint32) (item Item)
	First() (item Item)
	Pop() (item Item)
	Len() (length uint32)
	Serialize() (zero *common.ZeroCopySink)
}
