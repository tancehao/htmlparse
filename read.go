package htmlparse

import (
    "bytes"
)

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
				if c.isTag {
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
func (s *segment)getContent() []byte {
	if s == nil {
		return nil
	}
	if s.parent == nil { //root tag
		return s.tree.data
	}
	return s.tree.data[s.offset:s.limit]
}

//return the previous segment of a segment
func (s *segment)prev() *segment {
    if s == nil || s.parent == nil || len(s.parent.children) == 0 || s.parent.children[0] == s {
	    return nil
	}
    for i, t := range s.parent.children {
	    if t == s {
			return s.parent.children[i-1]
		}
	}
	return nil
}

//return the next segment of a segment
func (s *segment)next() *segment {
    if (s == nil) || (s.parent == nil) || (len(s.parent.children) == 0) || (s.parent.children[len(s.parent.children) - 1] == s) {
	    return nil
	}
    for i, t := range s.parent.children {
	    if t == s {
			return s.parent.children[i+1]
		}
	}
	return nil
}

/*
*
****************************************************
                 Tag's Methods
****************************************************
*
*/

//find a set of tags wrapped within a tag with a condition
func (t *Tag)Find(attr, value string) *TagSets {
	sets := &TagSets{
	    tags: []*Tag{t},
	}
	return sets.Find(attr, value)
}

//get the contents contained by a tag, including its tagname and attributes
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

//return the original data wrapped within a tag
func (t *Tag)Extract() []byte {
	if len(t.children) == 0 || t.NoEnd {
	    return []byte{}
	}
	leng := len(t.children)
	return t.segment.tree.data[t.children[0].offset:t.children[leng-1].limit]
}

//return the previous tag of a tag
func (t *Tag)Prev() *Tag {
	for seg := t.segment.prev(); seg != nil; seg = seg.prev(){
	    if seg.isTag == true {
		    return seg.tag
		}
	}
	return nil
}

//return the next tag of a tag
func (t *Tag)Next() *Tag {
    for seg := t.segment.next(); seg != nil; seg = seg.next() {
	    if seg.isTag == true {
		    return seg.tag
		}
	}
	return nil
}

//return the index of a tag in its parent
func (t *Tag)Index() int64 {
    return t.segment.index()
}

//check whether a tag is what we want using a filter
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
		value = bytes.TrimSpace(value)
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

func (s *TagSets)Select(chainSelector string) *TagSets {
    selectors, err := ParseChainSelector(chainSelector)
	if err != nil {
	    return &TagSets{
			tags: []*Tag{},
		}
	}
	ret := &TagSets{
	    tags: s.tags,
	}
	for slt := range selectors {
	    ret = ret.selectByOne(slt)
	}
	return ret
}

//select with css style selector
//one can't use some complex selectors except by tagName, by id or by class
func (s *TagSets)selectByOne(selector string) *TagSets {
    var filters map[string]string
	last := 0
	for i, s := range selector {
	    if i == 0 || (s != '.' && s != '#') {
		    continue
		}
		switch selector[last] {
		case '.':
			filters["class"] += fmt.Sprintf("%s ", selector[last+1:i])
		case '#':
			filters["id"] += selector[last+1:i]
		default:
			filters["tagName"] = selector[last+1:i]
		}
	}
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

//return the index of a text in its parent
func (t *Text)Index() int64 {
    return t.segment.index()
}

