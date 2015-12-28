package main

// import (
//          "fmt"
// )
// TODO: Not the best implementation, it leaks memory and should be changed. Only for proof of concept
// TODO: This needs mutex

type nutList struct {
	Elements map[string]nutElement
	Cap      int
	Count    int
}

type nutElement struct {
	PubKey       [32]byte
	Autenticated bool
	UnknownUser  bool
	Nr           int
}

func initnutList(cap int) *nutList {
	n := nutList{Cap: cap, Count: 0}
	n.Elements = make(map[string]nutElement)

	return &n
}

func (n *nutList) add(str string) {
	// if (len(n.Elements) >= n.Cap) {
	// 	n.Elements = n.Elements[1:] // Discard top element
	// }
	n.Count = n.Count + 1
	n.Elements[str] = nutElement{Nr: n.Count}
}

func (n *nutList) delete(str string) {
	delete(n.Elements, str)
}

func (n *nutList) get(key string) (nutElement, bool) {
	data, ok := n.Elements[key]
	if ok {
		return data, true
	}
	return nutElement{}, false
}

func (n *nutList) update(key string, ne nutElement) {
	elem, ok := n.Elements[key]
	if ok {
		elem = ne
		n.Elements[key] = elem
	}
}
