package htmlparse

import (
	"bytes"
    "errors"
    "strings"
)

var (
	NilTextError     = errors.New("the html to be parsed should not be null")
)

type Parser struct {
	data     []byte //original html data
	tagStack []*Tag //to store the opened tags temporarily when being parsed

	tree *Tree //the abstract tree to be generated
}

func NewParser(data []byte) *Parser {
    return &Parser{
	    data: data,
	}
}

func (p *Parser)Parse() (*Tree, error) {
	if len(p.data) == 0 {
		return nil, NilTextError
	}
	var n, offset int64
	var seg []byte
	var err error
	var tree *Tree
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
			if parent := p.getLastTag(); parent != nil {
				s := &segment{
				    tree: p.tree,
					parent: parent,
				}
				s.LinkToTag(tag, offset, n)
				parent.addChild(s)
			} else { //it's the root tag
				tree = &Tree{
					data: p.data,
					root: tag,
				}
				s := &segment{
					tree:   tree,
					parent: nil,
				}
				p.tree = tree
				s.LinkToTag(tag, 0, int64(len(p.data)))
			}

			//push it to the tag stack, and pop it when its close tag comes
			if !tag.NoEnd {
				p.tagStack = append(p.tagStack, tag)
			}
		} else if IsCloseTag(seg) {
			tagName := ReadWord(seg[2:])
			//pop its open tag, as well as those tags enbeded within them whose close tag were missed
			leng := len(p.tagStack)
			for i := leng - 1; i >= 0; i-- {
				if string(tagName) == p.tagStack[i].TagName {
					//p.tagStack[i].setLimit(offset + n)
					p.tagStack[i].segment.limit = offset + n
					p.tagStack = p.tagStack[:i]
					break
				}
			}
			//if no tag matches this close tag, treat it as text
			if leng == len(p.tagStack) {
				p.pushText(offset, n, seg)
			}
		} else {
			p.pushText(offset, n, seg)
		}
		offset = offset + n
	}
	return tree, nil
}

func (p *Parser)pushText(offset, n int64, text []byte) error {
	txt := make([]byte, len(text))
	copy(txt, text)    //a text can be updated, so we make a copy here
	t := &Text{
		Text: txt,
	}
	if parent := p.getLastTag(); parent != nil {
		s := &segment{
		    parent: parent,
			tree: p.tree,
		}
		s.LinkToText(t, offset, n)
		parent.addChild(s)
	}
	return nil
}

func (p *Parser)parseTag(tag []byte) (*Tag, error) {
	if tag[0] != '<' || tag[1] == '/' || tag[len(tag)-1] != '>' {
		return nil, NotTagError
	}
	tagName := ReadWord(tag[1:])
	newTag := &Tag{
		TagName: string(tagName),
	    Attributes: map[string]string{},
	    Class: map[string]bool{},
	    children: []*segment{},
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
		i := bytes.Index(pair, []byte("="))
		if i == -1 {
			newTag.Attributes[string(pair)] = ""
		} else {
			attr, value := string(pair[:i]), string(pair[i+1:])
			if WrappedBy(value, "'") || WrappedBy(value, string('"')) {
			    value = value[1:len(value)-1]
			}
			newTag.Attributes[attr] = value
			if attr == "class" {
				classes := strings.Split(value, " ")
				for _, c := range classes {
					newTag.Class[bytes.TrimSpace(c)] = true
				}
			}
		}
	}

	return newTag, nil
}

func (p *Parser)getLastTag() *Tag {
	s := p.tagStack
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}
