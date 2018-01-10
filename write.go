package htmlparse

import (
    "fmt"
	"errors"
    "strings"
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
	    attrs = append(attrs, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	str := fmt.Sprintf("<%s", t.TagName)
	if len(attrs) > 0 {
	    str += " "
	}
	str += fmt.Sprintf("%s>", strings.Join(attrs, " "))
	if t.NoEnd {
	    return str
	}
	for _, seg := range t.children {
	    if seg.isText {
		    str += seg.text.Modify()
		} else {
		    str += seg.tag.Modify()
		}
	}
	str += fmt.Sprintf("</%s>", t.TagName)
	return str
}

func (t *Text)Modify() string {
    return string(t.Text)
}

func (t *Tag)WriteText(position int64, data []byte) (*Text, error) {
	if t.NoEnd {
	    return nil, NoSpaceToWrite
	}
	text := &Text{
		Text: data,
	    segment: nil,
	}
	t.writeSegment(position, text)
	return text, nil
}

//write a tag to a parent tag
func (t *Tag)WriteTag(position int64, tagname string) (*Tag, error) {
	if t.NoEnd {
	    return nil, NoSpaceToWrite
	}
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
	return tag, nil
}

//write a segment of type text or tag to a tag
func (t *Tag)writeSegment(position int64, itf interface{}) *segment {
    if t.NoEnd {
	    return nil
	}
	if len(t.children) == 0 {
	    position = 0
	}
	if position >= int64(len(t.children)) && len(t.children) > 0 {
	    position = int64(len(t.children) - 1)
	}
	seg := &segment{
	    isTag: false,
		isText: false,
		tag: nil,
		text: nil,
		parent: t,
		tree: t.segment.tree,
		offset: 0,
		limit: 0,
	}
	switch itf.(type) {
		case *Text:
			s := itf.(*Text)
			seg.isText, seg.isTag = true, false
			s.segment = seg
			seg.text = s
	    case *Tag:
			s := itf.(*Tag)
			seg.isText, seg.isTag = false, true
			s.segment = seg
			seg.tag = s
	}
	cp := make([]*segment, len(t.children[position:]))
	copy(cp, t.children[position:])
	t.children = append(append(t.children[:position], seg), cp...)
    return seg
}

//delete a tag. whether to delete its children is optional
func (t *Tag)Delete(deleteChildren int) error {
    if t.segment.tree.root == t {
	    return RootUndeletable
	}
	switch deleteChildren {
	case 1:
        for i, seg := range t.segment.parent.children {
		    if seg == t.segment {
				cp := make([]*segment, len(t.segment.parent.children))
				copy(cp, t.segment.parent.children)
			    t.segment.parent.children = append(t.segment.parent.children[:i], cp[i+1:]...)
			}
		}
	case 0:
		for _, seg := range t.children {
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
