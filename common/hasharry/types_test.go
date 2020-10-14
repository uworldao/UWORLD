package hasharry

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestStringToHash(t *testing.T) {
	str := []byte("ABC")
	hash := sha256.Sum256(str)
	hashObj := BytesToHash(hash[:])
	fmt.Println(hashObj.String())
	hashObj2, _ := StringToHash(hashObj.String())
	fmt.Println(hashObj2.String())
	if !hashObj2.IsEqual(hash) {
		t.Errorf("error function")
	}
	t.Log("success")
}
