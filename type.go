package htmlparse

import (
    "errors"
)

var (
	TagNotExist         = errors.New("tag not exist")
	SelectorSyntaxError = errors.New("the css style selector has syntax error")
    NotContentError = errors.New("the content to be linked is not linkable")
)

type Tree struct {
	data []byte
	root *Tag
}

//a general model of a tag or text, which can be found with its absolute position
//a comment is treated as a text as well
type Segment struct {
	IsText bool
	IsTag  bool
	tag    *Tag
	text   *Text
	tree   *Tree
	Parent *Tag

	//offset and limit determins the absolute position 
	//of a segment in the document
	offset int64
	limit  int64
}

//an abstact of a html tag
type Tag struct {
	TagName    string
	Attributes map[string]string
	Class      map[string]bool
	NoEnd      bool //whether it's a single tag
	children   []*Segment
	segment    *Segment
}

type TagSets struct {
    tags []*Tag
}

type Text struct {
	text    []byte
	segment *Segment
}

//link an abstacted segment to a concrete tag which has many useful infos
func (s *Segment)LinkToTag(t *Tag, offset, n int64) {
	s.IsText = false
	s.IsTag = true
	s.text = nil
	s.tag = t
	s.offset = offset
	s.limit = offset + n
	t.segment = s
}

//link an abstact segment to a concrete text
func (s *Segment)LinkToText(t *Text, offset, n int64) {
	s.IsText = true
	s.IsTag = false
	s.text = t
	s.tag = nil
	s.offset = offset
	s.limit = offset + n
	t.segment = s
}

func (t *Tag)AddChild(s *Segment) {
    t.children = append(t.children, s)
}

func (t *Tag)SetLimit(n int64) {
	t.segment.limit = n
}

//create a new text. It is used when write data to a tag
func newText(parent *Tag, data []byte) *Text {
	segment := &Segment{
		IsText: true,
		IsTag:  false,
		tag:    nil,
		text:   nil,
		tree:   parent.segment.tree,
		Parent: parent,
		offset: 0,
		limit: 0,
	}
	text := &Text{
	    data: data,
	    segment: segment,
	}
	segment.text = text
	return text
}
