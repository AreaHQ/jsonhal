// Package jsonhal provides structs and methods to easily wrap your own data
// in a HAL compatible struct with support for hyperlinks and embedded resources
// HAL specification: http://stateless.co/hal_specification.html
package jsonhal

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

// Link represents a link in "_links" object
type Link struct {
	Href  string `json:"href" mapstructure:"href"`
	Title string `json:"title,omitempty" mapstructure:"title"`
}

// Embedded represents a resource in "_embedded" object
type Embedded interface{}

// EmbedSetter is the interface that wraps the basic setEmbedded method.
//
// SetEmbedded adds a slice of objects under a named key in the embedded map
type EmbedSetter interface {
	SetEmbedded(name string, embedded Embedded)
}

// EmbedGetter is the interface that wraps the basic getEmbedded method.
//
// GetEmbedded returns a slice of embedded resources by name or error
type EmbedGetter interface {
	GetEmbedded(name string) (Embedded, error)
}

// Embedder is the interface that wraps the basic setEmbedded and getEmbedded methods.
type Embedder interface {
	EmbedSetter
	EmbedGetter
}

// Hal is used for composition, include it as anonymous field in your structs
type Hal struct {
	Links    map[string]*Link    `json:"_links,omitempty" mapstructure:"_links"`
	Embedded map[string]Embedded `json:"_embedded,omitempty" mapstructure:"_embedded"`
	decoder  *mapstructure.Decoder
}

// SetLink sets a link (self, next, etc). Title argument is optional
func (h *Hal) SetLink(name, href, title string) {
	if h.Links == nil {
		h.Links = make(map[string]*Link, 0)
	}
	h.Links[name] = &Link{Href: href, Title: title}
}

// DeleteLink removes a link named name if it is found
func (h *Hal) DeleteLink(name string) {
	if h.Links != nil {
		delete(h.Links, name)
	}
}

// GetLink returns a link by name or error
func (h *Hal) GetLink(name string) (*Link, error) {
	if h.Links == nil {
		return nil, fmt.Errorf("Link \"%s\" not found", name)
	}
	link, ok := h.Links[name]
	if !ok {
		return nil, fmt.Errorf("Link \"%s\" not found", name)
	}
	return link, nil
}

// SetEmbedded adds a slice of objects under a named key in the embedded map
func (h *Hal) SetEmbedded(name string, embedded Embedded) {
	if h.Embedded == nil {
		h.Embedded = make(map[string]Embedded, 0)
	}
	h.Embedded[name] = embedded
}

// GetEmbedded returns a slice of embedded resources by name or error
func (h *Hal) GetEmbedded(name string) (Embedded, error) {
	if h.Embedded == nil {
		return nil, fmt.Errorf("Embedded \"%s\" not found", name)
	}
	embedded, ok := h.Embedded[name]
	if !ok {
		return nil, fmt.Errorf("Embedded \"%s\" not found", name)
	}
	return embedded, nil
}

// CountEmbedded counts number of embedded items
func (h *Hal) CountEmbedded(name string) (int, error) {
	e, err := h.GetEmbedded(name)
	if err != nil {
		return 0, err
	}
	if reflect.TypeOf(interface{}(e)).Kind() != reflect.Slice && reflect.TypeOf(interface{}(e)).Kind() != reflect.Map {
		return 0, errors.New("Embedded object is not a slice or a map")
	}
	return reflect.ValueOf(interface{}(e)).Len(), nil
}

// decodeHook is used to support datatypes that mapstructure does not support native
func (h *Hal) decodeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {

	// only if target datatype is time.Time  and if source datatype is string
	if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
		return time.Parse(time.RFC3339, data.(string))
	}

	//everything else would not be handled for now
	return data, nil
}

// DecodeEmbedded decodes embedded objects into a struct
func (h *Hal) DecodeEmbedded(name string, result interface{}) (err error) {
	var dec *mapstructure.Decoder
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)

		}
	}()

	e, err := h.GetEmbedded(name)
	if err != nil {
		panic(err)
	}
	//setup a new decoder if not already present
	if h.decoder == nil {
		dec, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: result, DecodeHook: h.decodeHook})
		if err != nil {
			panic(err)
		}
		h.decoder = dec
	}

	err = h.decoder.Decode(e)
	if err != nil {
		panic(err)
	}
	return nil
}

// DeleteEmbedded removes an embedded resource named name if it is found
func (h *Hal) DeleteEmbedded(name string) {
	if h.Embedded != nil {
		delete(h.Embedded, name)
	}
}
