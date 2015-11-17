package codec

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type jsonCodec struct{}

func (c jsonCodec) Encode(dest io.Writer) Encoder {
	return &jsonEncoder{
		out: dest,
	}
}

func (c jsonCodec) Decode(b []byte) Decoder {
	return &jsonDecoder{data: b}
}

type jsonEncoder struct {
	out io.Writer
}

func (e *jsonEncoder) One(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	e.out.Write(data)

	return nil
}

// All() encodes multiple JSON objects into one file.
// Right now this encodes to JSONLines.
func (e *jsonEncoder) All(vs ...interface{}) error {
	for _, v := range vs {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		e.out.Write(data)
		e.out.Write([]byte("\n"))
	}
	return nil
}

type jsonDecoder struct {
	data []byte
}

// All returns all documents in the original.
// JSON does not really support multi-doc, so we try JSON decoding, then fall
// back on JSONList decoding.
//
// Decoding using All will make up to two decoding passes on your data before
// it determines which decoder to use (so it can easily be more than 2).
func (d jsonDecoder) All() ([]*Object, error) {
	var phony interface{}
	if err := json.Unmarshal(d.data, &phony); err == nil {
		// We have a single JSON document.
		return []*Object{{data: d.data, dec: jdec}}, nil
	}

	lines := bytes.Split(d.data, []byte("\n"))

	// If it is 0, it's not JSON. If it's 1, it should have parsed.
	if len(lines) < 2 {
		return []*Object{}, errors.New("Data is neither JSON nor JSONL. (linecount)")
	}

	println(string(lines[0]))
	if err := json.Unmarshal(lines[0], &phony); err != nil {
		// Whoops.. failed again.
		return []*Object{}, errors.New("Data is neither JSON nor JSONL: " + err.Error())
	}

	buf := make([]*Object, len(lines))
	for i, l := range lines {
		buf[i] = &Object{data: l, dec: jdec}
	}
	return buf, nil
}

func (d jsonDecoder) One() (*Object, error) {
	return &Object{
		data: d.data,
		dec:  jdec,
	}, nil
}

func jdec(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
