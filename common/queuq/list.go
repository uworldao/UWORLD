package queuq

type ListQueue struct {
	e    Element
	next *ListQueue
	head *ListQueue
}

func NewListQueue() *ListQueue {
	el := &ListQueue{}
	return &ListQueue{head: el}
}

func (l *ListQueue) Put(e Element) {
	n1 := l.head.next
	n2 := l.head
	for n1 != nil {
		n2 = n1
		n1 = n1.next
	}
	n1 = &ListQueue{e: e}
	n2.next = n1
}

func (l *ListQueue) Pop() Element {
	e := l.head.next.e
	l.head.next = l.head.next.next
	return e
}

func (l *ListQueue) Get(e Element) Element {
	n := l.head.next
	for n != nil {
		if n.e.GetElementKey() == e.GetElementKey() {
			return n.e
		}
		n = n.next
	}
	return nil
}
