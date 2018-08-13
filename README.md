htmlparse
===

Htmlparse is a go tool for parsing a html document.  
  
It converts a html document into a tree. Each node in the tree is either a tag or a text. Given a tag, a programmer  
  
can easily get its original infos, including its metadata, its children, its siblings and the text wrapped in it.  
  
One can also modify a tree, by writing something into or delete a tag.  
  
It can be used in web crawlers, analysis, batch formating and etc.  

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
```go
    type Parser struct {
        //unexported fields
    }
```

* (p *Parser)Parse() *Tag

	The only one method needed to convert the original bytes to a tree. The root tag was returned. 
	#### Example:
	```go
		import (
			"github.com/tancehao/htmlparse"		
		)
        
		//...
		
	    content, _ := ioutil.ReadFile("index.html")
		parser := htmlparse.NewParser(content)
		root, err := parser.Parse()
	```


### Tag
```go
    type Tag struct {
        TagName    string
        Attributes map[string]string
        Class      []string
        NoEnd      bool    //whether it's a single tag
        
        //unexported fields
    }
```

* (t *Tag)HasClass(class string) bool

    Whether a given tag has a class.
    #### Example
    ```go
        div, _ := htmlparse.NewParser(`<div class="product box">...</div>`).Parse()
        fmt.Println(div.HasClass("box"))  // true
    ```

* (t *Tag)Find(map[string]string) *TagSets
	
    Find the tags from a tag's children by a set of conditions. 
    #### Example
    ```go
        formBytes := []byte(`<form method="post" action"test">
          <input type="hidden" name="foo1" value="bar1" />
          <input type="hidden" name="foo2" value="bar2" />
          <input type="text" name="foo3" value="bar3" />
          <textarea name="foo4" value="bar4" />
          </form>
        `)
        formTag, _ := htmlparse.NewParser(formBytes).Parse()
        hiddenInputes := formTag.Find(map[string]string{
            "tagName": "input",
            "type": "hidden",
        })
        fmt.Println(hiddenInputes)
        //<input type="hidden" name="foo1" value="bar1" />
        //<input type="hidden" name="foo2" value="bar2" />
    ```


* (t *Tag)FindByName(name string) *TagSets

    Find a set of tags by their names.


* (t *Tag)FindByClass(class string) *TagSets
    
    Find a set of tags by their classes.


* (t *Tag)FindByFunc(f func(*Tag) bool) *TagSets

    Find a set of tags by a function which tells whether a tag should by returned.
    #### Example:
    ```go
        root.FindByFunc(func(t *htmlparse.Tag) bool {
            cond1 := t.TagName == "img"
            cond2 := t.HasClass("product")
            cond3 := strings.Contains(t.Attributes["Src"], ".png")
            return cond1 && cond2 && cond3
        })
    ```


* (t *Tag)FindByCssSelector(path string) *TagSets

    Find a set of tags by a css selector.
    #### Example:
    ```go
        productBytes := []byte(`<div class="product">
            <div class="pictures"></div>
            <div class="infos">
                <div class="brand">BRAND</brand>
                <div class="name">NAME</div>
                <div class="prices">
                    $1000.00 <span><s>$1200.00</s></span>
                </div>
            </div>
        </div>`)
        productTag, _ := htmlparse.Next(productBytes).Parse()
        originalPrices := productTag.FindByCssSelector(" .infos .prices span s")
        fmt.Println(originalPrices)
        //$1200.00
    ```


* (t *Tag)GetContent() []byte

	Return the original bytes of a tag in the document, the tag's metadata is included.  
	By design, each tag or text has a pair of pointers pointing to its head and tail in the original document.
    This function searchs for the contents using the pair of indexes, which guarantees the speed when it works.


* (t *Tag)String() string
	
    Implements the fmt.Stringer.  


* (t *Tag)Extract() []byte
	
    Filter the text from the original data of a tag. Tags wont't be included.  
    #### Example
    ```go
        navifatorBytes := []byte(`<div class="navs">
            <a href="/category/society">Society</a>
            <a href="/category/sports">Sports</a>
            <a href="/category/mil">Mil</a>
            <a href="/category/cars">Cars</a>
            <a href="/category/tech">Tech</a>
        </div>`)
        navigatorTag, _ := htmlparse.NewParser(navifatorBytes).Parse()
        fmt.Println(string(navigatorTag.Extract()))
        //Society
        //Sports
        //Mil
        //Cars
        //Tech
    ```


* (t *Tag)Index() int

	Get the index of a tag among its siblings.  


* (t *Tag)Prev() *Tag
	
    Get the previous tag of an existing one under the same parent.  


* (t *Tag)Next() *Tag
	
    Similar to Prev().  


* (t *Tag)Modify() string
	
    Return the modified data of a tag.
    When reading something from a tag, one gets the original data from the document. However, when writing a tag, one is not writing to the original document but the tag struct, so this method gets the original data plus the differences one has made to a tag.


* (t *Tag)WriteText(position int, data []byte) (*Text, error)

	Write text into a tag at the given index in the tag's children.  
	Example:
	```go
		names := products.Find(map[string]string{"class": "productName"})
		fmt.Println(names)
		//prints:
		//<div class="productName">Product1</div>
		//<div class="productName">Product2</div>
		//<div class="productName">Product3</div>

		for _, name := range names.All() {
			name.WriteText(0, []byte("[ONSALE] "))
			fmt.Println(name.Modify())
		}

		//prints:
		//<div class="productName">[ONSALE] Product1</div>
		//<div class="productName">[ONSALE] Product2</div>
		//<div class="productName">[ONSALE] Product3</div>
	```

* (t *Tag)WriteTag(position int, tagname string) (*Tag, error)

	Write a tag into a tag at the given index in the tag's chidren. 
	Example:
	```go
		script, _ := body.WriteTag(1000, "script")
		//if ths position is greater than the count of the tag's children, it'll be set to the last
		script.Attributes["src"] = "http://www.foo.com"
		fmt.Println(body.Modify())

        //a script tag was inserted into the html document.
	```

* (t *Tag)Delete() error
	
    Delete a tag.  
	Example:
	```go
		garbages := html.FindByClass("ad").All()
        for _, g := range garbages {
            g.Delete()
        }
	    fmt.Println(html.Modify())
	```


### TagSets
```go
    type TagSets struct {
        //unexported fields
    }
```

* (ts *TagSets)Find(map[string]string) *TagSets
	
    Return a set of tags from a set of tags or their children using a filter.  
	It can be used in a chain style. 
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


* (ts *TagSets)All() []*Tag
	
    Get all the tags in this set.  


* (ts *TagSets)GetAttributes(attr ...string) []map[string]string
	
    Get some attributes from each tag in this set.  
	Example:
	```go
        
        formBytes := []byte(`<form method="post" action"test">
          <input type="hidden" name="foo1" value="bar1" />
          <input type="hidden" name="foo2" value="bar2" />
          <input type="text" name="foo3" value="bar3" />
          <textarea name="foo4" value="bar4" />
          </form>
        `)
        formTag, _ := htmlparse.NewParser(formBytes).Parse()
		inputs := formTag.Find(map[string]string{
			"tagName": "input"
		}).GetAttributes("type", "name", "value")
		for _, input := range inputs {
			fmt.Printf("%s,%s,%s\n", 
				input["type"], input["name"], input["value"]
			)
		}
	    //hidden,foo1,bar1
        //hidden,foo2,bar2
        //text,foo3,bar3
        
    ```

* (ts *TagSets)String() string


### Text
```go
    type Text struct {
        Text []byte
        
        //unexported fields
    }
```

* (ts *Text)String() string
	
    Smilar to tag.


* (ts *Text)Index() int
	
    Smilar to tag.

* (t *Text)Modify() string
	
    Smilar to tag.

* (t *Text)Delete()
	
    Smilar to tag.

