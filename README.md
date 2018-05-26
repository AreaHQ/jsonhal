## jsonhal

A simple Go package to make custom structs marshal into [HAL](http://stateless.co/hal_specification.html) compatible JSON responses.

[![Travis Status for RichardKnop/jsonhal](https://travis-ci.org/RichardKnop/jsonhal.svg?branch=master&label=linux+build)](https://travis-ci.org/RichardKnop/jsonhal)
[![godoc for RichardKnop/jsonhal](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/RichardKnop/jsonhal)
[![codecov for RichardKnop/jsonhal](https://codecov.io/gh/RichardKnop/jsonhal/branch/master/graph/badge.svg)](https://codecov.io/gh/RichardKnop/jsonhal)

---


Just add `jsonhal.Hal` as anonymous field to your structs and use `SetLink` to set hyperlinks and optionally `SetEmbedded` to set embedded resources.

Example:

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/RichardKnop/jsonhal"
)

// HelloWorld ...
type HelloWorld struct {
	jsonhal.Hal
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Foobar ...
type Foobar struct {
	jsonhal.Hal
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func main() {
	var (
		helloWorld   *HelloWorld
		jsonResponse []byte
		err          error
	)

	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink("self", "/v1/hello/world/1", "")

	jsonResponse, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(jsonResponse))
	// {
	// 	"_links": {
	// 		"self": {
	// 			"href": "/v1/hello/world/1"
	// 		}
	// 	},
	// 	"id": 1,
	// 	"name": "Hello World"
	// }

	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self", // name
		"/v1/hello/world?offset=2&limit=2", // href
		"", // title
	)
	helloWorld.SetLink(
		"next", // name
		"/v1/hello/world?offset=4&limit=2", // href
		"", // title
	)
	helloWorld.SetLink(
		"previous",                         // name
		"/v1/hello/world?offset=0&limit=2", // href
		"", // title
	)
	jsonResponse, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(jsonResponse))
	// {
	// 	"_links": {
	// 		"next": {
	// 			"href": "/v1/hello/world?offset=4\u0026limit=2"
	// 		},
	// 		"previous": {
	// 			"href": "/v1/hello/world?offset=0\u0026limit=2"
	// 		},
	// 		"self": {
	// 			"href": "/v1/hello/world?offset=2\u0026limit=2"
	// 		}
	// 	},
	// 	"id": 1,
	// 	"name": "Hello World"
	// }

	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink("self", "/v1/hello/world/1", "")

	// Add embedded resources
	foobars := []*Foobar{
		&Foobar{
			Hal: jsonhal.Hal{
				Links: map[string]*jsonhal.Link{
					"self": &jsonhal.Link{Href: "/v1/foo/bar/1"},
				},
			},
			ID:   1,
			Name: "Foo bar 1",
		},
		&Foobar{
			Hal: jsonhal.Hal{
				Links: map[string]*jsonhal.Link{
					"self": &jsonhal.Link{Href: "/v1/foo/bar/2"},
				},
			},
			ID:   2,
			Name: "Foo bar 2",
		},
	}
	helloWorld.SetEmbedded("foobars", Embedded(foobars))

	jsonResponse, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(jsonResponse))
	// {
	// 	"_links": {
	// 		"self": {
	// 			"href": "/v1/hello/world/1"
	// 		}
	// 	},
	// 	"_embedded": {
	// 		"foobars": [
	// 			{
	// 				"_links": {
	// 					"self": {
	// 						"href": "/v1/foo/bar/1"
	// 					}
	// 				},
	// 				"id": 1,
	// 				"name": "Foo bar 1"
	// 			},
	// 			{
	// 				"_links": {
	// 					"self": {
	// 						"href": "/v1/foo/bar/2"
	// 					}
	// 				},
	// 				"id": 2,
	// 				"name": "Foo bar 2"
	// 			}
	// 		]
	// 	},
	// 	"id": 1,
	// 	"name": "Hello World"
	// }
}
```
