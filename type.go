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
	Text    []byte
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

//add a child of type tag to a tag
func (parent *Tag)AddChildTag(child *Tag, position int) *Tag {
    if len(parent) == 0 {
	    position = 0
	}
	seg := &Segment{
	    IsTag: true,
		IsText: false,
		tag: child,
		text: nil,
		parent: parent,
		tree: parent.segment.tree,
	}
	child.segment = seg
    parent.children = append(append(parent.children[:position], seg), parent.children[position:])
    return parent
}

//add a child of type text to a tag
func (parent *Tag)AddTextChild(child *Text, position int) *Tag {
    if len(parent) == 0 {
	    position = 0
	}
	seg := &Segment{
	    IsTag: false,
		IsText: true,
		tag: nil,
		text: child,
		parent: parent,
		tree: parent.segment.tree,
	}
	child.segment = seg
    parent.children = append(append(parent.children[:position], seg), parent.children[position:])
    return parent
}

