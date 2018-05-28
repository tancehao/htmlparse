package htmlparse

import (
	"errors"
	"fmt"
	"strings"
)

var (
	NoSpaceToWrite  = errors.New("single tag has no space to write in")
	RootUndeletable = errors.New("can't delete the root tag")
)

//an abstact of a html tag
type Tag struct {
	TagName    string
	Attributes map[string]string
	Class      []string
	NoEnd      bool //whether it's a single tag
	children   []*segment
	segment    *segment
}

//whether a tag has a class
func (t *Tag) HasClass(class string) bool {
	for _, c := range t.Class {
		if c == class {
			return true
		}
	}
	return false
}

//filter a set of tags from a tag and it's children by tag name
func (t *Tag) FindByName(name string) []*Tag {
	return t.Find(map[string]string{"tagName": name})
}

//filter a set of tags from a tag and it's children by class
func (t *Tag) FindByClass(class string) []*Tag {
	return t.FindWithFunc(func(t *Tag) bool {
		return t.HasClass(class)
	})
}

//filter a subset of tags from a tag's children
func (t *Tag) Find(conds map[string]string) []*Tag {
	return t.FindWithFunc(func(tag *Tag) bool {
		return tag.checkByConditions(conds)
	})
}

func (t *Tag) FindWithFunc(f func(*Tag) bool) []*Tag {
	result := []*Tag{}
	if f(t) {
		result = append(result, t)
	}
	for _, seg := range t.children {
		if seg.isTag && f(seg.tag) {
			result = append(result, seg.tag.FindWithFunc(f)...)
			//result = append(result, seg.tag)
		}
	}
	return result
}

//check if a tag satisfies a set of conditions
func (t *Tag) checkByConditions(conds map[string]string) bool {
	for k, v := range conds {
		if t.checkByCondition(k, v) == false {
			return false
		}
	}
	return true
}

//check if a tag satisfies a condition
func (t *Tag) checkByCondition(attr, value string) bool {
	switch attr {
	case "tagName":
		if t.TagName == value {
			return true
		} else {
			return false
		}
	case "class":
		return t.HasClass(value)
	default:
		v, ok := t.Attributes[attr]
		if ok && (value == v || value == "") {
			return true
		}
	}
	return false
}

//get the contents contained by a tag, including its tagname and attributes
func (t *Tag) GetContent() []byte {
	if t == nil {
		return nil
	}
	return t.segment.getContent()
}

func (t *Tag) GetAttribute(attr string) string {
	switch attr {
	case "tagName":
		return t.TagName
	case "content":
		return string(t.Extract())
	default:
		v, ok := t.Attributes[attr]
		if ok {
			return v
		} else {
			return ""
		}
	}
}

//if only returns the original data
//one should call Modify() to get the modified data
func (t *Tag) String() string {
	if t == nil {
		return ""
	}
	return string(t.GetContent())
}

//remove the metadata of a tag
//all of its children were returned
func (t *Tag) Unwrap() []byte {
	if len(t.children) == 0 || t.NoEnd {
		return []byte{}
	}
	leng := len(t.children)
	return t.segment.tree.data[t.children[0].offset:t.children[leng-1].limit]
}

//return the text wrapped in a tag with all metadata of tags removed
func (t *Tag) Extract() []byte {
	if len(t.children) == 0 || t.NoEnd {
		return []byte{}
	}
	text := []byte{}
	for _, c := range t.children {
		if c.isText {
			text = append(text, ' ') //split the bytes from diffrent places with a blank
			text = append(text, c.getContent()...)
		} else {
			text = append(text, c.tag.Extract()...)
		}
	}
	return text
}

//return the previous tag of a tag
func (t *Tag) Prev() *Tag {
	for seg := t.segment.prev(); seg != nil; seg = seg.prev() {
		if seg.isTag == true {
			return seg.tag
		}
	}
	return nil
}

//return the next tag of a tag
func (t *Tag) Next() *Tag {
	for seg := t.segment.next(); seg != nil; seg = seg.next() {
		if seg.isTag == true {
			return seg.tag
		}
	}
	return nil
}

//return the index of a tag in its parent
func (t *Tag) Index() int64 {
	return t.segment.index()
}

//return the modified data of a tag
func (t *Tag) Modify() string {
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

func (t *Tag) WriteText(position int64, data []byte) (*Text, error) {
	if t.NoEnd {
		return nil, NoSpaceToWrite
	}
	text := &Text{
		Text:    data,
		segment: nil,
	}
	t.writeSegment(position, text)
	return text, nil
}

//write a tag to a parent tag
func (t *Tag) WriteTag(position int64, tagname string) (*Tag, error) {
	if t.NoEnd {
		return nil, NoSpaceToWrite
	}
	noend := IsSingleTag(tagname)
	tag := &Tag{
		TagName:    tagname,
		Attributes: map[string]string{},
		Class:      []string{},
		NoEnd:      noend,
		children:   []*segment{},
		segment:    nil,
	}
	t.writeSegment(position, tag)
	return tag, nil
}

//write a segment of type text or tag to a tag
func (t *Tag) writeSegment(position int64, itf interface{}) *segment {
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
		isTag:  false,
		isText: false,
		tag:    nil,
		text:   nil,
		parent: t,
		tree:   t.segment.tree,
		offset: 0,
		limit:  0,
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
func (t *Tag) Delete(deleteChildren int) error {
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

func (t *Tag) addChild(s *segment) {
	t.children = append(t.children, s)
}
