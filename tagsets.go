package htmlparse

import (
	"errors"
)

var (
	TagNotExist         = errors.New("tag not exist")
	SelectorSyntaxError = errors.New("the css style selector has syntax error")
	NotContentError     = errors.New("the content to be linked is not linkable")
)

type TagSets struct {
    tags []*Tag
}

//filter a tags set to another one, it can be used with a chain style
//e.g. TagSets.Find("tagName", "form").Find("method", "post").Find("tagName", "input").Find("type", "text")
func (t *TagSets) Find(filter map[string]string) *TagSets {
	tags := FilterTags(t.tags, filter)
	return &TagSets{
		tags: tags,
	}
}

func (t *TagSets) FindByName(name string) *TagSets {
	return t.Find(map[string]string{"tagName": name})
}

func (t *TagSets) FindByClass(class string) *TagSets {
	return t.Find(map[string]string{"class": class})
}

//get the tag list
func (t *TagSets) All() []*Tag {
	return t.tags
}

func (t *TagSets) GetAttributes(attrs ...string) []map[string]string {
	var ret []map[string]string
	for _, tag := range t.tags {
		values := map[string]string{}
		for _, attr := range attrs {
		    values[attr] = tag.GetAttribute(attr)
		}
		ret = append(ret, values)
	}
	return ret
}

func (t *TagSets) String() string {
	s := ""
	for _, t := range t.tags {
		s += t.String()
	}
	return s
}
