// Package codec provides a JSON/YAML codec for Manifests.
//
// Usage:
//
// 	// Decode one manifest from a JSON file.
// 	man, err := JSON.Decode(b).One()
// 	// Decode all of the manifests out of this file.
// 	manifests, err := YAML.Decode(b).All()
// 	err := YAML.Encode(filename).One("hello")
//	// Encode multiple objects to one file (as separate docs).
// 	err := YAML.Encode(filename).All("one", "two", "three")
package codec

import (
	"io"
	"io/ioutil"
)

// JSON is the default JSON encoder/decoder.
var JSON jsonCodec

// YAML is the default YAML encoder/decoder.
var YAML yamlCodec

// Encoder describes something capable of encoding to a given format.
//
// An Encoder should be able to encode one object to an output stream, or
// many objects to an output stream.
//
// For example, a single YAML file can contain multiple YAML objects, and
// a single JSONList file can contain many JSON objects.
type Encoder interface {
	// Write one object to one file
	One(interface{}) error
	// Write all objects to one file
	All(...interface{}) error
}

// Decoder decodes an encoded representation into one or many objects.
type Decoder interface {
	// Get one object from a file.
	One() (*Object, error)
	// Get all objects from a file.
	All() ([]*Object, error)
}

// Codec has an encoder and a decoder for a particular encoding.
type Codec interface {
	Encode(io.Writer) Encoder
	Decode([]byte) Decoder
}

// DecodeFile returns a decoder pre-populated with the file contents.
func DecodeFile(filename string, c Codec) (Decoder, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return c.Decode(data), nil
}
