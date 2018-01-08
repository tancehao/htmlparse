package htmlparse

import (
	"strings"
	"fmt"
)

/*
*
******************************************************
                  Tree's Methods
******************************************************
*
*/

func (t *Tree)Filter(filter map[string]string) []*Tag {
	tags := []*Tag{t.root}
	return FilterTags(tags, filter)
}

//find some tags from a html tree
func (t *Tree)Find(attr, value string) *TagSets {
	return t.root.Find(attr, value)
}

//convert a tree to string
//it returns the modified document
func (t *Tree)String() string {
    return t.root.String()
}

/*
*
*******************************************************
                  segment's Methods
*******************************************************
*/

//get the actual bytes of a segment
func (s *Segment)getContent() []byte {
	if s == nil {
		return nil
	}
	if s.Parent == nil { //root tag
		return s.tree.data
	}
	return s.tree.data[s.offset:s.limit]
}

//filter a tag sets to another with conditions
func FilterTags(originTags []*Tag, filter map[string]string) []*Tag {
	result := []*Tag{}
	for _, tag := range originTags {
		if len(tag.children) > 0 { //if a tag has children, process them recursively
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

/*
*
****************************************************
                 Tag's Methods
****************************************************
*
*/

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

//if only returns the original data
//one should call Modify() to get the modified data
func (t *Tag)String() string {
	if t == nil {
	    return ""
	}
	return string(t.GetContent())
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

/*
*
****************************************************
               TagSets's Methods
****************************************************
*
*/

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


/*
*
*************************************************
                  Text's Methods
*************************************************
*
*/
func (t *Text)String() string {
	if t == nil {
	    return ""
	}
    return string(t.segment.getContent())
}
