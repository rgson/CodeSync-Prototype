package queue

type FIFO struct {
	Head *Element
}

type Element struct {
	Value interface{}
	Next *Element
}

func (q *FIFO) Enqueue(e interface{}) {

	element := &Element{Value: e}

	if q.Head != nil {
		element.Next = q.Head.Next
		q.Head.Next = element
	}
	
	q.Head = element

}

func (q *FIFO) Dequeue() (value interface{}) {

	if q.Head != nil {	
		value = q.Head.Next.Value
		q.Head.Next = q.Head.Next.Next
	}
	
	return

}

func (q *FIFO) Peek() (value interface{}) {

	if q.Head != nil {
		value = q.Head.Next.Value
	}
	
	return
}

func (q *FIFO) Count() (count int) {

	if q.Head == nil {
		return 0
	}

	element := q.Head.Next
	for element != q.Head {
		count++
		element = element.Next
	}
	return

}

func (q *FIFO) ToArray() (elements []interface{}) {

	count := q.Count()
	elements = make([]interface{}, count)
	
	element := q.Head
	for i := 0; i < count; i++ {
		element = element.Next
		elements[i] = *element
	}
	
	return
}
