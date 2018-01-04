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

func (t *Tree)Filter(filter map[string]string) []*Tag {
	tags := []*Tag{t.root}
	return FilterTags(tags, filter)
}

//find some tags from a html tree
func (t *Tree)Find(attr, value string) *TagSets {
	return t.root.Find(attr, value)
}

//a general model of a tag or text, which can be found with its absolute position
//a comment is treated as a text as well
type Segment struct {
	IsText bool
	IsTag  bool
	tree   *Tree
	Parent *Tag
	tag    *Tag
	text   *Text
	offset int64 //在父元素内的起始位置
	limit  int64 //在父元素内的结束位置
}

//get the actual bytes of a segment
func (s *Segment)getContent() []byte {
	if s == nil {
		return nil
	}
	if s.Parent == nil { //根元素
		return s.tree.data
	}
	return s.tree.data[s.offset:s.limit]
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

//an abstact of a html tag
type Tag struct {
	TagName    string
	Attributes map[string]string
	Class      map[string]bool
	NoEnd      bool //没有关闭标签？
	children   []*Segment
	segment    *Segment
}

//filter a tag sets to another with conditions
func FilterTags(originTags []*Tag, filter map[string]string) []*Tag {
	result := []*Tag{}
	for _, tag := range originTags {
		if len(tag.children) > 0 { //有下级标签，递归地查询
			subTags := []*Tag{}
			if tag.checkByFilter(filter) {
			    result = append(result, tag)
			}
			for _, c := range tag.children {
				if c.IsTag {
					subTags = append(subTags, c.tag)
				}
			}
			if len(subTags) > 0 {
				ts := FilterTags(subTags, filter)
				for _, t := range ts {
				    result = append(result, t)
				}
			}
		} else {
			if tag.checkByFilter(filter) {
			    result = append(result, tag)
			}
		}
	}
	return result
}

func (t *Tag)Find(attr, value string) *TagSets {
	sets := &TagSets{
	    tags: []*Tag{t},
	}
	return sets.Find(attr, value)
}

//get the contents contained by a tag, except its tagname and attributes
func (t *Tag) GetContent() []byte {
	if t == nil {
		return nil
	}
	return t.segment.getContent()
}

func (t *Tag)String() string {
    return string(t.GetContent())
}

func (t *Tag)AddChild(s *Segment) {
    t.children = append(t.children, s)
}

func (t *Tag)checkByFilter(filter map[string]string) bool {
    for k, v := range filter {
	    if !t.checkByCondition(k, v) {
		    return false
		}
	}
	return true
}

//check if one tag satisfies the condition
func (t *Tag)checkByCondition(attr, value string) bool {
	switch attr {
	case "tagName":
		if t.TagName == value {
			return true
		} else {
			return false
		}
	case "class":
		if _, ok := t.Class[value]; ok {
			return true
		}
	default:
		v, ok := t.Attributes[attr]
		if ok && (value == v || value == "") {
			return true
		}
	}
	return false
}

func (t *Tag)Extract() []byte {
	if len(t.children) == 0 || t.NoEnd {
	    return []byte{}
	}
	leng := len(t.children)
	return t.segment.tree.data[t.children[0].offset:t.children[leng-1].limit]
}

func (t *Tag)SetLimit(n int64) {
	t.segment.limit = n
}

type TagSets struct {
    tags []*Tag
}

//filter a tags set to another one, it can be used with a chain style
//e.g. TagSets.Find("tagName", "form").Find("method", "post").Find("tagName", "input").Find("type", "text")
func (t *TagSets)Find(attr, value string) *TagSets {
	var filter = map[string]string{
	    attr: value,
	}
	tags := FilterTags(t.tags, filter)
	return &TagSets{
	    tags: tags,
	}
}

func (t *TagSets)All() []*Tag {
    return t.tags;
}

func (t *TagSets)One() *Tag {
    if len(t.tags) == 0 {
	    return nil
	}
	return t.tags[0]
}

func (t *TagSets)GetAttributes(attr string) []string {
	attrs := []string{}
	for _, tag := range t.tags {
		if v, ok := tag.Attributes[attr]; ok {
		    attrs = append(attrs, v)
		}
	}
	return attrs
}

func (t *TagSets)String() string {
	s := ""
	for _, t := range t.tags {
	    s += t.String()
	}
	return s
}

type Text struct {
	text    []byte
	segment *Segment
}

