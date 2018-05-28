package htmlparse

type Tree struct {
	data []byte
	root *Tag
}

//find some tags from a html tree
func (t *Tree) Find(filter map[string]string) *TagSets {
	//return t.root.Find(filter)
	return (&TagSets{tags: []*Tag{t.root}}).Find(filter)
}

//find the tags by name
func (t *Tree) FindByName(name string) *TagSets {
	return (&TagSets{tags: []*Tag{t.root}}).FindByName(name)
}

//find the tags by a class
func (t *Tree) FindByClass(class string) *TagSets {
	return (&TagSets{tags: []*Tag{t.root}}).FindByClass(class)
}

//convert a tree to string
//it returns the modified document
func (t *Tree) String() string {
	return t.root.String()
}

//one should call this method to get
//the updated data of a document
func (t *Tree) Modify() string {
	return t.root.Modify()
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

//link an abstacted segment to a concrete tag which has many useful infos
func (s *segment) linkToTag(t *Tag, offset, n int64) {
	s.isText = false
	s.isTag = true
	s.text = nil
	s.tag = t
	s.offset = offset
	s.limit = offset + n
	t.segment = s
}

//link an abstact segment to a concrete text
func (s *segment) linkToText(t *Text, offset, n int64) {
	s.isText = true
	s.isTag = false
	s.text = t
	s.tag = nil
	s.offset = offset
	s.limit = offset + n
	t.segment = s
}

func (seg *segment) index() int64 {
	for i, s := range seg.parent.children {
		if s == seg {
			return int64(i)
		}
	}
	return -1
}

//get the actual bytes of a segment
func (s *segment) getContent() []byte {
	if s == nil {
		return nil
	}
	if s.parent == nil { //root tag
		return s.tree.data
	}
	return s.tree.data[s.offset:s.limit]
}

//return the previous segment of a segment
func (s *segment) prev() *segment {
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
func (s *segment) next() *segment {
	if (s == nil) || (s.parent == nil) || (len(s.parent.children) == 0) || (s.parent.children[len(s.parent.children)-1] == s) {
		return nil
	}
	for i, t := range s.parent.children {
		if t == s {
			return s.parent.children[i+1]
		}
	}
	return nil
}
