package htmlarse

import (
    "fmt"
	"errors"
)

var (
    NoSpaceToWrite = errors.New("single tag has no space to write in")
    RootUndeletable = errors.New("can't delete the root tag")
)

//one should call this method to get 
//the updated data of a document
func (t *Tree)Modify() string {
    return t.root.Modify();
}

//return the modified data of a tag
func (t *Tag)Modify() string {
	attrs := []string{}
	for k, v := range t.Attributes {
	    attrs = append(attrs, fmt.Fprintf("%s=\"%s\"", k, v))
	}
	str := fmt.Fprintf("<%s %s>", t.TagName, strings.Join(attrs, " ")
	if t.NoEnd {
	    return str
	}
	for seg := range t.children {
	    if seg.IsText {
		    str += seg.text.Modify()
		} else {
		    str += seg.tag.Modify()
		}
	}
	str += fmt.Fprintf("</%s>", t.TagName)
	return str
}

func (t *Text)Modify() string {
    return string(t.Text)
}


func (t *Tag)Write(position int64, data []byte) (n, error) {
    if t.NoEnd {
	    return 0, NoSpaceToWrite
	}
	if len(t.children) == 0 {
		text := newText(t, data)
	}
}

func (t *Tag)WriteAfter(data []byte) (n, error) {}

func (t *Tag)WriteBefore(data []byte) (n, error) {}

//delete a tag. whether to delete its children is optional
func (t *Tag)Delete(deleteChildren int) error {
    if t.segment.tree.root == t {
	    return RootUndeletable
	}
	switch deleteChildren {
	case 1:
        for i, seg := range t.segment.Parent.children {
		    if seg == t.segment {
			    t.segment.Parent.children = append(t.segment.Parent.children[:i], t.segment.Parent.children[i+1:]...)
			}
		}
	case 0:
		for i, seg := range t.children {
		    seg.Parent = t.segment.Parent
			t.segment.Parent.children = append(t.segment.Parent.children, seg)
		}
	}
	return nil
}

//delete a text
func (t *Text)Delete() error {
    for i, seg := range t.segment.Parent.children {
	    if seg == t.segment {
		    t.segment.Parent.children = append(t.segment.Parent.children[:i], t.segment.Parent.children[i+1:]...)
		}
	}
	return nil
}
