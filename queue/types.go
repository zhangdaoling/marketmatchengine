package queue

import "github.com/zhangdaoling/marketmatchengine/common"

type Item interface {
	Compare(other interface{}) int
	Key() uint32
}

type PriorityQueue interface {
	Insert(item Item) (i Item)
	Cancel(key uint32) (item Item)
	First() (item Item)
	Pop() (item Item)
	Len() (length int)
	Serialize() (zero *common.ZeroCopySink)
}
