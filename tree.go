package htmlparse

//a general model of a tag or text, which can be found with its absolute position
//a comment is treated as a text as well
type segment struct {
	isText bool
	isTag  bool
	tag    *Tag
	text   *Text
	parent *Tag

	//offset and limit determins the absolute position
	//of a segment in the document
	offset int64
	limit  int64

	data []byte //the source bytes, empty if it's not the root tag
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
	tmp := s
	for tmp.parent != nil {
		tmp = tmp.parent
	}
	return tmp.data[s.offset:s.limit]
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
