package types

import (
	"fmt"
	"testing"
)

func TestCalCoinBase(t *testing.T) {
	var i uint64
	var sum uint64
	for i = 0; i <= 10368000; i++ {
		x := CalCoinBase(i, 1)
		if x != 0 {
			fmt.Println(i, " = ", x)
		}
		sum += x
	}
	fmt.Println(sum)
}
