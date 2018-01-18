package htmlparse

import (
	"strings"
)

type Text struct {
	Text    []byte
	segment *segment
}

func (t *Text) String() string {
	if t == nil {
		return ""
	}
	return string(t.segment.getContent())
}

//return the index of a text in its parent
func (t *Text) Index() int64 {
	return t.segment.index()
}

//delete a text
func (t *Text) Delete() error {
	for i, seg := range t.segment.parent.children {
		if seg == t.segment {
			t.segment.parent.children = append(t.segment.parent.children[:i], t.segment.parent.children[i+1:]...)
		}
	}
	return nil
}

func (t *Text) Modify() string {
	return strings.TrimSpace(string(t.Text))
}
