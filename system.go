package system

import (
	"fmt"

	"github.com/goptos/runtime"
	"honnef.co/go/js/dom/v2"
)

type Scope = runtime.Scope
type Elem struct {
	dom.Element
}
type Event = dom.Event

func (self *Elem) New(tag_name string) *Elem {
	var window = dom.GetWindow()
	var document = window.Document()
	var el = document.CreateElement(tag_name)
	return &Elem{el}
}

func (self *Elem) On(event_name string, cb func(dom.Event)) *Elem {
	self.AddEventListener(event_name, true, cb)
	return self
}

func (self *Elem) Attr(attr_name string, attr_value string) *Elem {
	self.SetAttribute(attr_name, attr_value)
	return self
}

func (self *Elem) Text(data string) *Elem {
	var window = dom.GetWindow()
	var document = window.Document()
	var node = document.CreateTextNode(data)
	self.AppendChild(node)
	return self
}

func (self *Elem) DynText(cx *Scope, f func() string) *Elem {
	var window = dom.GetWindow()
	var document = window.Document()
	var node = document.CreateTextNode("")
	self.AppendChild(node)
	cx.CreateEffect(func() {
		var value = f()
		node.SetNodeValue(value)
	})
	return self
}

func (self *Elem) Child(node *Elem) *Elem {
	self.AppendChild(node)
	return self
}

func (self *Elem) DynChild(cx *Scope, f func() bool, node *Elem) *Elem {
	cx.CreateEffect(func() {
		var contains = false
		for _, child := range self.ChildNodes() {
			if node.IsEqualNode(child) {
				contains = true
				break
			}
		}
		if f() {
			if !contains {
				self.AppendChild(node)
			}
		} else {
			if contains {
				self.RemoveChild(node)
			}
		}
	})
	return self
}

func Each[T any](elem *Elem,
	cx *Scope,
	cF func() []T,
	kF func(T) int,
	vF func(*Scope, T) *Elem) *Elem {
	var window = dom.GetWindow()
	var document = window.Document()
	var m = make(map[string]T)
	cx.CreateEffect(func() {
		var items = cF()
		var mK = make(map[string]struct{})
		// add
		for _, item := range items {
			var key = fmt.Sprintf("%d", kF(item))
			_, ok := m[key]
			if !ok {
				elem.Child(vF(cx, item).Attr("id", key))
				m[key] = item
			}
			mK[key] = struct{}{}
		}
		// remove
		for key := range m {
			_, ok := mK[key]
			if !ok {
				var child = document.GetElementByID(key)
				if child != nil {
					child.Remove()
				}
				delete(m, key)
			}
		}
	})
	return elem
}

func Mount(view func(cx *Scope) *Elem) {
	var cx = (*Scope).New(nil)
	var window = dom.GetWindow()
	var document = window.Document()
	var body = document.GetElementsByTagName("body")[0]
	var root = view(cx)
	body.AppendChild(root)
}

func Run(view func(cx *Scope) *Elem) {
	c := make(chan struct{}, 0)
	Mount(view)
	<-c
}

// Todo DynAttr (etc.)
