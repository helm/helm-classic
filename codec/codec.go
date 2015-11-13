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
)

var JSON jsonCodec
var YAML yamlCodec

type Encoder interface {
	// Write one manifest to one file
	One(interface{}) error
	// Write all objects to one file
	All(...interface{}) error
}

type Decoder interface {
	One() (*Object, error)
	All() ([]*Object, error)
}

type Codec interface {
	Encode(io.Writer) Encoder
	Decode([]byte) Decoder
}
