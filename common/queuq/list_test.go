package queuq

import "testing"

type element struct {
	Data int
	Key  int
}

func (e *element) GetElementKey() interface{} {
	return e.Key
}

func TestListQueue_Put(t *testing.T) {
	queue := NewListQueue()
	queue.Put(&element{1, 1})
	queue.Put(&element{2, 2})
	queue.Put(&element{3, 3})

	e1 := queue.Pop()
	if e1.GetElementKey() != 1 {
		t.Fatalf("error")
	}
	e2 := queue.Pop()
	if e2.GetElementKey() != 2 {
		t.Fatalf("error")
	}
	e3 := queue.Pop()
	if e3.GetElementKey() != 3 {
		t.Fatalf("error")
	}
}

func TestListQueue_Get(t *testing.T) {
	queue := NewListQueue()
	queue.Put(&element{1, 1})
	queue.Put(&element{2, 2})
	queue.Put(&element{3, 3})

	e := queue.Get(&element{3, 3})
	if e != nil {
		t.Log(e)
	} else {
		t.Fatalf("err %v ", e)
	}
}
