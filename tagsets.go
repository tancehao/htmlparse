package htmlparse

import (
	"errors"
	"strings"
)

var CssSelecotrSeparators = []byte{
	' ', '.', '#', '>', '+', '[', ']',
}

var (
	TagNotExist         = errors.New("tag not exist")
	SelectorSyntaxError = errors.New("the css style selector has syntax error")
	NotContentError     = errors.New("the content to be linked is not linkable")
)

type TagSets struct {
	tags []*Tag
}

//find a set of tags from the children of tags in a set with some conditions
func (t *TagSets) Find(conds map[string]string) *TagSets {
	return t.FindWithFunc(func(tag *Tag) bool {
		return tag.checkByConditions(conds)
	})
}

//find a set of tags from the children of tags in a set with a function
func (t *TagSets) FindWithFunc(f func(*Tag) bool) *TagSets {
    result := &TagSets{}
    //result := []*Tag{}
	for _, tag := range t.tags {
		result.merge(tag.FindWithFunc(f))
	}
	return result
}

//convert a tag set to another one using a function
func (t *TagSets) Map(f func(*Tag) *Tag) *TagSets {
	result := []*Tag{}
	for _, tag := range t.tags {
		result = append(result, f(tag))
	}
	return &TagSets{tags: result}
}

func (t *TagSets) FindByName(name string) *TagSets {
	return t.Find(map[string]string{"tagName": name})
}

func (t *TagSets) FindByClass(class string) *TagSets {
	return t.FindWithFunc(func(tag *Tag) bool {
		return tag.HasClass(class)
	})
}

//get the tag list
func (t *TagSets) All() []*Tag {
	return t.tags
}

func (t *TagSets) push (tag *Tag) {
    t.tags = append(t.tags, tag)
}

func (t *TagSets) merge (src *TagSets) {
    t.tags = append(t.tags, src.tags...)
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

func (t *TagSets) HasTag(tag *Tag) bool {
	for _, tg := range t.tags {
		if tg == tag {
			return true
		}
	}
	return false
}

func (t *TagSets) String() string {
	s := ""
	for _, t := range t.tags {
		s += t.String()
	}
	return s
}

//find a set of tags using a css selector
func (t *TagSets) FindByCssSelector(path string) (ret *TagSets) {
	result := []*Tag{}
	paths := strings.Split(path, ",")
	for _, p := range paths {
		subset := t.findByPath(p)
		if subset == nil {
			break
		}
		for _, tag := range subset.tags {
			if t.HasTag(tag) {
				continue
			}
			result = append(result, tag)
		}
	}
	return &TagSets{tags: result}
}

func (t *TagSets) findByPath(path string) *TagSets {
	var subset []*Tag
	var selector string
	selector = readSelector(path)
	if selector == "" {
		return t
	}
	/* Note that there should not simply call Find, because most selectors have specified meanings on relationship between tags */
	switch selector[0] {
	case '#':
		for _, tag := range t.tags {
			if tag.checkByCondition("id", selector[1:]) {
				subset = append(subset, tag)
			}
		}
	case '.':
		for _, tag := range t.tags {
			if tag.HasClass(selector[1:]) {
				subset = append(subset, tag)
			}
		}
	case ' ': //all the descendant
		for _, tag := range t.tags {
			ss := tag.FindWithFunc(func(tg *Tag) bool {
				return true
			})
			/* the original tag must be the first element */
			subset = append(subset, ss.tags[1:]...)
		}
	case '>':
		for _, tag := range t.tags {
			for _, child := range tag.children {
				if child.isTag {
					subset = append(subset, child.tag)
				}
			}
		}
	case '+':
		for _, tag := range t.tags {
			t := tag
			for t.Next() != nil {
				subset = append(subset, t)
				t = t.Next()
			}
		}

	case '[':

	default: //tag name
	    for _, tag := range t.tags {
            if tag.TagName == selector {
                subset = append(subset, tag)
            }
        }
    }
	return (&TagSets{tags: subset}).findByPath(path[len(selector):])
}

func readSelector(path string) (selector string) {
	for i := 1; i < len(path); i++ {
		for _, s := range CssSelecotrSeparators {
			if s == path[i] {
				return path[:i]
			}
		}
	}
	return path
}
