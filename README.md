htmlparse
===

Htmlparse is a go tool for parsing a html document into a tree.

The basic structure of this package is the Tag struct. Given a tag, you can easily do these things:

* Get its attributes.
* Extract its text contents as fast as possible.
* Find the tags you want from its children with some filters.

---

* [Install](#install)
* [Api](#api)
* [Examples](#examples)

---

## Install

```sh
go get -u github.com/tancehao/htmlparse
```

## api
### Tree
* #### Filter(filter map[string]string) []*Tag
* #### Find(attr, value string) *TagSets
* #### String() string
* #### Modify() string


### Tag
* #### Find(attr, value string) *TagSets
* #### GetContent() []byte
* #### String() string
* #### Extract() []byte
* #### Index() int64
* #### Prev() *Tag
* #### Next() *Tag
* #### Modify() string
* #### WriteText(position int64, data []byte) (*Text, error)
* #### WriteTag(position int64, tagname string) (*Tag, error)
* #### Delete() error


### TagSets
* #### Find(attr, value string) *TagSets
* #### All() []*Tag
* #### GetAttributes(attr string) []string
* #### String() string


### Text
* #### String() string
* #### Index() int64
* #### Modify() string
* #### Delete()

## Examples

```go
func main() {
    content, _ := ioutil.ReadFile("baidu.html")
	parser := htmlparse.NewParser(content)
	tree, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	
	//find some tags with a filter
    filter := map[string]string{
	    "tagName": "input",
		"type": "hidden",
	}
    for _, tag := range tree.Filter(filter) {
	    fmt.Println(tag)
	}
}
```

the codes above prints:
```go
<input type=hidden name=bdorz_come value=1> 
<input type=hidden name=ie value=utf-8> 
<input type=hidden name=f value=8> 
<input type=hidden name=rsv_bp value=1> 
<input type=hidden name=rsv_idx value=1> 
<input type=hidden name=tn value=baidu>
```

Or you can use a chain:
```go
	imgTags := tree.Find("id", "products").Find("class", "product").Find("tagName", "img").Find("class", "product_photo")
	photos := imgTags.GetAttributes("src")
	fmt.Println(photos)
```

Prints:
```go
img13.360buyimg.com/n5/s54x54_jfs/t10675/253/1344769770/66891/92d54ca4/59df2e7fN86c99a27.jpg
img13.360buyimg.com/n5/s54x54_jfs/t10408/221/1296471743/54709/b8bf69a6/59df2e82N1f855465.jpg
img13.360buyimg.com/n5/s54x54_jfs/t10438/209/1299858339/22067/94f4941f/59df2e82N11980eca.jpg
img13.360buyimg.com/n5/s54x54_jfs/t10357/189/1296983301/32053/a0e9283e/59df2e82N3e8f5183.jpg
img13.360buyimg.com/n5/s54x54_jfs/t10450/186/1312686657/70920/1167f96b/59df2e83Nc6f15397.jpg
```
