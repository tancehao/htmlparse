package htmlparse

import (
	"bytes"
	"errors"
	"strings"
)

var (
	NilTextError = errors.New("the html to be parsed should not be null")
)

type Parser struct {
	data []byte //original html data

	currentContainer *Tag
	root             *Tag //the root tag of the dom tree, which contains the source bytes
}

func NewParser(data []byte) *Parser {
	return &Parser{
		data: data,
	}
}

func (p *Parser) Parse() (*Tag, error) {
    if len(p.data) == 0 {
		return nil, NilTextError
	}
	var n, offset int64
	var seg []byte
	var err error
	var tag *Tag
	for {
		n, seg, err = ReadSegment(p.data, offset)
        if err != nil {
			break
		}
		if IsOpenTag(seg) {
			tag, err = p.parseTag(seg)
			if err != nil {
				continue
			}

			//build relationship with its parent when it's not the root tag
			if p.root != nil {
				s := &segment{
					parent: p.currentContainer,
				}
				s.linkToTag(tag, offset, n)
				s.parent.addChild(s)
            } else { //it's the root tag
				s := &segment{
					parent: nil,
					data:   p.data,
				}
				s.linkToTag(tag, 0, int64(len(p.data)))
				p.root = tag
			}
			if tag.NoEnd == false {
				p.currentContainer = tag
			}
        } else if IsCloseTag(seg) {
			tagName := string(ReadWord(seg[2:]))
			for p.currentContainer.TagName != tagName && p.currentContainer != p.root {
				//the close tag of the current container was misseed
				p.currentContainer.segment.limit = offset
				p.currentContainer = p.currentContainer.segment.parent
			}
			p.currentContainer.segment.limit = offset + n
			if p.currentContainer.TagName == tagName {
				p.currentContainer = p.currentContainer.segment.parent
			} else {
				break
			}
		} else {
			if p.currentContainer == nil {
                break
            }
            p.pushText(offset, n, seg)
		}
		offset = offset + n
	}
	return p.root, nil
}

func (p *Parser) parseTag(tag []byte) (*Tag, error) {
	if tag[0] != '<' || tag[1] == '/' || tag[len(tag)-1] != '>' {
		return nil, NotTagError
	}
	tagName := ReadWord(tag[1:])
	newTag := &Tag{
		TagName:    string(tagName),
		Attributes: map[string]string{},
		Class:      []string{},
		children:   []*segment{},
	}

	if IsSingleTag(string(tagName)) {
		newTag.NoEnd = true
	}

	tag = tag[len(tagName)+1 : len(tag)-1]
	if len(tag) == 0 {
		return newTag, nil
	}
	//parse the attributes
	attrValues := bytes.Split(tag, []byte(" "))
	for _, pair := range attrValues {
		if len(pair) == 0 {
			continue
		}
		i := bytes.Index(pair, []byte("="))
		if i == -1 {
			newTag.Attributes[string(pair)] = ""
		} else {
			attr, value := string(pair[:i]), string(pair[i+1:])
			if WrappedBy(value, "'") || WrappedBy(value, string('"')) {
				value = value[1 : len(value)-1]
			}
			newTag.Attributes[attr] = value
			if attr == "class" {
				classes := strings.Split(value, " ")
				for _, c := range classes {
					newTag.Class = append(newTag.Class, strings.TrimSpace(c))
				}
			}
		}
	}

	return newTag, nil
}

func (p *Parser) pushText(offset, n int64, text []byte) error {
	txt := make([]byte, len(text))
	copy(txt, text) //a text can be updated, so we make a copy here
	t := &Text{
		Text: txt,
	}
    if p.root != nil {
		s := &segment{
			parent: p.currentContainer,
		}
		s.linkToText(t, offset, n)
        p.currentContainer.addChild(s)
	}
	return nil
}

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
		tmp = tmp.parent.segment
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
