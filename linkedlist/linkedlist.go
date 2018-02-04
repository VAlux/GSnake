package linkedlist

import "fmt"

// Node is a single list node containing pointer to the next node and data payload
type Node struct {
	next *Node
	Data interface{}
}

// LinkedList is a list of nodes connected by 'next' pointers in one direction
type LinkedList struct {
	head *Node
	size int
}

func (node *Node) String() string {
	return fmt.Sprintf("data: %s", node.Data)
}

// Init returns initialized linked list
func (list *LinkedList) Init() *LinkedList {
	list.size = 0
	return list
}

// New creates empty linked list
func New() *LinkedList {
	list := new(LinkedList).Init()
	return list
}

// Size returns the actual amount of elements in the linked list
func (list *LinkedList) Size() int {
	return list.size
}

// Next returns the next node relative to the current node
func (n *Node) Next() *Node {
	return n.next
}

// Head returns first node of the list
func (list *LinkedList) Head() *Node {
	return list.head
}

// Back returns last node of the list
func (list *LinkedList) Back() *Node {
	current := list.head
	for current != nil && current.next != nil {
		current = current.next
	}
	return current
}

// Append adds node to the end of the list
func (list *LinkedList) Append(node *Node) {
	if list.head == nil {
		list.head = node
	} else {
		list.Back().next = node
	}
	list.size++
}

// AppendList appends whole linked list to the end of existing one
func (list *LinkedList) AppendList(newList *LinkedList) {
	if list.head == nil {
		list = newList
	} else {
		list.Back().next = newList.head
	}
	list.size += newList.size
}

// Prepend adds node before the first elemnt of the list
func (list *LinkedList) Prepend(node *Node) {
	if list.head == nil {
		list.head = node
	} else {
		node.next = list.head
		list.head = node
	}
	list.size++
}

// PrependList appends whole lined list before the head of existing one
func (list *LinkedList) PrependList(newList *LinkedList) {
	newList.Back().next = list.head
	list.head = newList.head
	list.size += newList.size
}

// RemoveLast remove last node from the list and return it
func (list *LinkedList) RemoveLast() *Node {
	current := list.head
	for current.next.next != nil {
		current = current.next
	}
	last := current.next
	current.next = nil
	list.size--

	return last
}

// GetAt returns node of the list at the specified index beginning from 0
func (list *LinkedList) GetAt(index int) *Node {
	current := list.head
	for i := 0; i < index && current.next != nil; i++ {
		current = current.next
	}
	return current
}
