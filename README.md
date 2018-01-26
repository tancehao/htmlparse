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

---

## Install

```sh
go get -u github.com/tancehao/htmlparse
```

## api
### Parser
* #### Parse() *Tree
	The only one method needed to convert the original bytes to a tree.  
	Example:
	```go
		content, _ := ioutil.ReadFile("index.html")
		parser := htmlparse.NewParse(content)
		tree := parser.Parse()
	```

### Tree
* #### Filter(filter map[string]string) []*Tag
	Find some tags from the document with a filter, which is a key-value formated map.  
	Example:
	```go
		products := tree.Filter(map[string]string{"tagName", "div", "class": "product"})
	```

* #### Find(conditions map[string]string) *TagSets
	Similar to the Filter method, except that its return value is of TagSets type who has some useful methods.  

* #### String() string
	Return the original document.  

* #### Modify() string
	Return the modified document.  


### TagSets
* #### Find(map[string]string) *TagSets
	Return a set of tags from a set of tags or their children using a filter.  
	It can be used with a chain style.  
	Example:
	```go
		photos := tree.Find(map[string]string{
			"tagName": "div", 
			"class": "product",
		}).Find(map[string]string{
			"tagName": "img",
			"class": "photo",
		})
	```

* #### All() []*Tag
	Get all the tags in this set.  

* #### GetAttributes(attr ...string) []map[string]string
	Get some attributes from each tag in this set.  
	Example:
	```go
		inputs := form.Find(map[string]string{
			"tagName": "input"
		}).GetAttributes("type", "name", "value", "data-id")
		for _, input := range inputs {
			fmt.Printf("%s,%s,%s,%s\n", 
				input["type"], input["name"], input["value"], input["data-id"]
			)
		}
	```

* #### String() string


### Tag
* #### Find(map[string]string) *TagSets
	Find the tags from a tag's children.  

* #### GetContent() []byte
	Return the original bytes of a tag in the document, the tag's metadata is included.  
	By design, each tag or text has a pair of pointers which determined its absolute position in the document.  
	So whenever one gets the original content of a tag or text, it just fetches the subslice document[head:tail],   
	which can be no more faster.

* #### String() string
	Satisfy the Stringer.  

* #### Extract() []byte
	Filter the text from the original data of a tag. Tags wont't be included.  

* #### Index() int64
	Get the index of a tag in its among its brothers.  

* #### Prev() *Tag
	Get the previous tag of a tag under the same parent.  

* #### Next() *Tag
	Similar to Prev().  

* #### Modify() string
	Return the modified data of a tag.  
	One should call this after writing to a tag.  

* #### WriteText(position int64, data []byte) (*Text, error)
	Write text into a tag at the given index in the tag's children.  
	Example:
	```go
		names := products.Find(map[string]string{"class": "productName"})
		fmt.Println(names)
		//prints:
		//<div class="productName">Product1</div>
		//<div class="productName">Product2</div>
		//<div class="productName">Product3</div>
		//<div class="productName">Product4</div>

		for _, name := range names {
			name.WriteText(0, []byte("[ONSALE] "))
			fmt.Println(name.Modify())
		}

		//prints:
		//<div class="productName">[ONSALE] Product1</div>
		//<div class="productName">[ONSALE] Product2</div>
		//<div class="productName">[ONSALE] Product3</div>
		//<div class="productName">[ONSALE] Product4</div>
	```

* #### WriteTag(position int64, tagname string) (*Tag, error)
	Write a tag into a tag at the given index in the tag's chidren.  
	Example:
	```go
		script, _ := body.WriteTag(1000, "script")    
		//if ths position is greater than the count of the tag's children, it'll be set to the last
		script.Attributes["src"] = "http://www.foo.com"
		body.Modify()
	```

* #### Delete() error
	Delete a tag.  
	Example:
	```go
		garbage := tree.Find(map[string]string{"class":"advertisement"}).All()[0]
		garbage.Delete()
	    tree.Modify()
	```

### Text
* #### String() string
	Smilar to tag.

* #### Index() int64
	Smilar to tag.

* #### Modify() string
	Smilar to tag.

* #### Delete()
	Smilar to tag.

