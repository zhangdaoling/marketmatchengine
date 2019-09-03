package queue

type Item interface {
	Compare(other interface{}) int
	Key() int32
}

type PriorityQueue interface {
	Insert(item Item) (i Item)
	Cancel(key int32) (item Item)
	First() (item Item)
	Pop() (item Item)
}
