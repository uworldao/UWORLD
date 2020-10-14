package queuq

type Element interface {
	GetElementKey() interface{}
}

type Queue interface {
	Put(e Element)
	Pop() Element
	Get(e Element) Element
}
