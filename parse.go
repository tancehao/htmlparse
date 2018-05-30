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

func (p *Parser) Parse() (*Tree, error) {
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
					parent: parent,
				}
				s.linkToTag(tag, offset, n)
				parent.addChild(s)
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
			tagName := ReadWord(seg[2:])
			for p.currentContainer.TagName != tagName && p.currentContainer != p.root {
				//the close tag of the current container was misseed
				p.currentContainer.segment.limit = offset
				p.currentContainer = p.currentContainer.segment.parent.tag
			}
			p.currentContainer.segment.limit = offset + n
			if p.currentContainer.TagName == tagName {
				p.currentContainer = p.currentContainer.segment.parent.tag
			} else {
				break
			}
		} else {
			p.pushText(offset, n, seg)
		}
		offset = offset + n
	}
	return p.root, nil
}

func (p *Parser) pushText(offset, n int64, text []byte) error {
	txt := make([]byte, len(text))
	copy(txt, text) //a text can be updated, so we make a copy here
	t := &Text{
		Text: txt,
	}
	if parent := p.getLastTag(); parent != nil {
		s := &segment{
			parent: parent,
		}
		s.linkToText(t, offset, n)
		parent.addChild(s)
	}
	return nil
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

func (p *Parser) getLastTag() *Tag {
	s := p.tagStack
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}
