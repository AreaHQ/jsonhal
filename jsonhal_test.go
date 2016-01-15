package jsonhal

import (
	"bytes"
	"log"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

// HelloWorld is a simple test struct
type HelloWorld struct {
	Hal
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Foobar is a simple test struct
type Foobar struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Qux is a simple test struct
type Qux struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

var expectedJSON = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON2 = []byte(`{
	"_links": {
		"next": {
			"href": "/v1/hello/world?offset=4\u0026limit=2"
		},
		"previous": {
			"href": "/v1/hello/world?offset=0\u0026limit=2"
		},
		"self": {
			"href": "/v1/hello/world?offset=2\u0026limit=2"
		}
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON3 = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"_embedded": {
		"foobars": [{
			"id": 1,
			"name": "Foo bar 1"
		}, {
			"id": 2,
			"name": "Foo bar 2"
		}]
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON4 = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"_embedded": {
		"foobars": [{
			"id": 1,
			"name": "Foo bar 1"
		}, {
			"id": 2,
			"name": "Foo bar 2"
		}],
		"quxes": [{
			"id": 1,
			"name": "Qux 1"
		}, {
			"id": 2,
			"name": "Qux 2"
		}]
	},
	"id": 1,
	"name": "Hello World"
}`)

func TestHal(t *testing.T) {
	var (
		helloWorld        *HelloWorld
		expected          *bytes.Buffer
		actual            []byte
		err               error
		foobars           []*Foobar
		quxes             []*Qux
		embeddedResources []Embedded
	)

	// Let's test a simple scenario with just a self link
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's text more links
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

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON2)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's test embedded resources
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Add embedded foobars
	foobars = []*Foobar{
		&Foobar{ID: 1, Name: "Foo bar 1"},
		&Foobar{ID: 2, Name: "Foo bar 2"},
	}
	embeddedResources = make([]Embedded, len(foobars))
	for i, foobar := range foobars {
		embeddedResources[i] = Embedded(foobar)
	}
	helloWorld.SetEmbedded("foobars", embeddedResources)

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON3)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's test multiple embedded resources
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Add embedded foobars
	foobars = []*Foobar{
		&Foobar{ID: 1, Name: "Foo bar 1"},
		&Foobar{ID: 2, Name: "Foo bar 2"},
	}
	embeddedResources = make([]Embedded, len(foobars))
	for i, foobar := range foobars {
		embeddedResources[i] = Embedded(foobar)
	}
	helloWorld.SetEmbedded("foobars", embeddedResources)

	// Add embedded quxes
	quxes = []*Qux{
		&Qux{ID: 1, Name: "Qux 1"},
		&Qux{ID: 2, Name: "Qux 2"},
	}
	embeddedResources = make([]Embedded, len(quxes))
	for i, qux := range quxes {
		embeddedResources[i] = Embedded(qux)
	}
	helloWorld.SetEmbedded("quxes", embeddedResources)

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON4)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))
}

func TestGetLink(t *testing.T) {
	helloWorld := new(HelloWorld)

	var (
		link *Link
		err  error
	)

	// Test when object has no links
	link, err = helloWorld.GetLink("self")
	assert.Nil(t, link)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Link \"self\" not found", err.Error())
	}

	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Test getting a bogus link
	link, err = helloWorld.GetLink("bogus")
	assert.Nil(t, link)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Link \"bogus\" not found", err.Error())
	}

	// Test getting a valid link
	link, err = helloWorld.GetLink("self")
	assert.Nil(t, err)
	if assert.NotNil(t, link) {
		assert.Equal(t, "/v1/hello/world/1", link.Href)
		assert.Equal(t, "", link.Title)
	}
}

func TestGetEmbedded(t *testing.T) {
	helloWorld := new(HelloWorld)

	var (
		embedded          []Embedded
		err               error
		foobars           []*Foobar
		embeddedResources []Embedded
	)

	// Test when object has no embedded resources
	embedded, err = helloWorld.GetEmbedded("foobars")
	assert.Nil(t, embedded)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Embedded \"foobars\" not found", err.Error())
	}

	// Add embedded foobars
	foobars = []*Foobar{
		&Foobar{ID: 1, Name: "Foo bar 1"},
		&Foobar{ID: 2, Name: "Foo bar 2"},
	}
	embeddedResources = make([]Embedded, len(foobars))
	for i, foobar := range foobars {
		embeddedResources[i] = Embedded(foobar)
	}
	helloWorld.SetEmbedded("foobars", embeddedResources)

	// Test geting bogus embedded resources
	embedded, err = helloWorld.GetEmbedded("bogus")
	assert.Nil(t, embedded)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Embedded \"bogus\" not found", err.Error())
	}

	// Test geting valid embedded resources
	embedded, err = helloWorld.GetEmbedded("foobars")
	assert.Nil(t, err)
	if assert.NotNil(t, embedded) {
		assert.Equal(t, 2, len(embedded))
	}
}
