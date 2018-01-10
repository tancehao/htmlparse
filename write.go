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
	    if seg.isText {
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

//write a []byte to a tag
func (t *Tag)Write(position int64, data []byte) (*Text, error) {
    if t.NoEnd {
	    return nil, NoSpaceToWrite
	}
	if len(t.children) == 0 {
	    position = 0
	}
	if position >= len(t.children) {
	    position = len(t.children) - 1
	}
	seg := &segment{
	    isTag: false,
		isText: true,
		tag: nil,
		text: nil,
		parent: t,
		tree, t.segment.tree,
		offset: 0,
		limit: 0,
	}
	text := &Text{
	    Text: data,
		segment: seg,
	}
	seg.text = text
	t.children = append(append(t.children[:position], seg), parent.children[position:]...)
    return text, nil
}

func (t *Tag)WriteTag(position int64, tagname string) *Tag {
	noend := IsSingleTag(tagname)
	tag := &Tag{
	    TagName: tagname,
		Attributes: map[string]string{},
		Class: map[string]bool{},
		NoEnd: noend,
		children: []*segment{},
		segment: nil,
	}
	t.writeSegment(position, tag)
	return tag
}

//write a segment of type text or tag to a tag
func (t *Tag)writeSegment(position int64, segment interface{}) *segment {
    if t.NoEnd {
	    return nil, NoSpaceToWrite
	}
	if len(t.children) == 0 {
	    position = 0
	}
	if position >= len(t.children) {
	    position = len(t.children) - 1
	}
	seg := &segment{
	    isTag: false,
		isText: false,
		tag: nil,
		text: nil,
		parent: t,
		tree, t.segment.tree,
		offset: 0,
		limit: 0,
	}
	switch segment.(type) {
		case *Text:
			s := segment.(*Text)
			seg.isText, seg.isTag = true, false
			s.segment = seg
			seg.text = s
	    case *Tag:
			s := segment.(*Tag)
			seg.isText, seg.isTag = false, true
			s.segment = seg
			seg.tag = s
	}
	t.children = append(append(t.children[:position], seg), parent.children[position:]...)
    return seg

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
        for i, seg := range t.segment.parent.children {
		    if seg == t.segment {
			    t.segment.parent.children = append(t.segment.parent.children[:i], t.segment.parent.children[i+1:]...)
			}
		}
	case 0:
		for i, seg := range t.children {
		    seg.parent = t.segment.parent
			t.segment.parent.children = append(t.segment.parent.children, seg)
		}
	}
	return nil
}

//delete a text
func (t *Text)Delete() error {
    for i, seg := range t.segment.parent.children {
	    if seg == t.segment {
		    t.segment.parent.children = append(t.segment.parent.children[:i], t.segment.parent.children[i+1:]...)
		}
	}
	return nil
}
