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
type segment struct {
	isText bool
	isTag  bool
	tag    *Tag
	text   *Text
	tree   *Tree
	parent *Tag

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
	children   []*segment
	segment    *segment
}

type TagSets struct {
    tags []*Tag
}

type Text struct {
	Text    []byte
	segment *segment
}

//link an abstacted segment to a concrete tag which has many useful infos
func (s *segment)LinkToTag(t *Tag, offset, n int64) {
	s.isText = false
	s.isTag = true
	s.text = nil
	s.tag = t
	s.offset = offset
	s.limit = offset + n
	t.segment = s
}

//link an abstact segment to a concrete text
func (s *segment)LinkToText(t *Text, offset, n int64) {
	s.isText = true
	s.isTag = false
	s.text = t
	s.tag = nil
	s.offset = offset
	s.limit = offset + n
	t.segment = s
}

func (seg *segment)index() int64 {
    for i, s := range seg.parent.children {
	    if s == seg {
		    return int64(i)
		}
	}
	return -1
}

func (t *Tag)addChild(s *segment) {
    t.children = append(t.children, s)
}

