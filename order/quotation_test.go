package order

import (
	"sort"
	"testing"
)

func TestInsert(t *testing.T) {
}

func TestSubAmount(t *testing.T) {
}

func TestBinarySerach(t *testing.T) {
	a := []int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55}
	x := 54
	i := search(a, x)
	t.Logf("i=%d", i)
	i = sort.Search(len(a), func(i int) bool { return a[i] >= x })
	if i < len(a) && a[i] == x {
		t.Logf("found %d at index %d in %v\n", x, i, a)
	} else {
		t.Logf("%d not found in %v\n", x, a)
	}
}

func search(arr []int, value int) int {
	i, j := 0, len(arr)
	for i < j {
		h := int(uint(i+j) >> 1)
		// i ≤ h < j
		if !(arr[h] >= value) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	return i
}
