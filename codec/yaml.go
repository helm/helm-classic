package codec

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/ghodss/yaml"
)

type yamlCodec struct{}

func (c yamlCodec) Decode(b []byte) Decoder {
	return &yamlDecoder{data: b}
}

func (c yamlCodec) Encode(out io.Writer) Encoder {
	return &yamlEncoder{out: out}
}

type yamlEncoder struct {
	out io.Writer
}

func (e *yamlEncoder) One(v interface{}) error {
	buf, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	e.out.Write(buf)
	return nil
}

func (e *yamlEncoder) All(vs ...interface{}) error {
	c := len(vs) - 1
	for i, v := range vs {
		if err := e.One(v); err != nil {
			return err
		}
		if i < c {
			e.out.Write([]byte(yamlSeparator))
			e.out.Write([]byte("\n"))
		}
	}
	return nil
}

type yamlDecoder struct {
	data []byte
}

// All returns all documents in a single YAML file.
func (d *yamlDecoder) All() ([]*Object, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(d.data))
	scanner.Split(SplitYAMLDocument)

	ms := []*Object{}
	for scanner.Scan() {
		m := &Object{
			data: append([]byte(nil), scanner.Bytes()...),
			dec: func(b []byte, v interface{}) error {
				return yaml.Unmarshal(b, v)
			},
		}
		ms = append(ms, m)
	}

	return ms, scanner.Err()
}

// One returns no more than one YAML doc, even if the file contains more.
func (d *yamlDecoder) One() (*Object, error) {
	ms, err := d.All()
	if err != nil {
		return nil, err
	}
	if len(ms) == 0 {
		return nil, errors.New("No document")
	}
	return ms[0], nil
}

const yamlSeparator = "\n---"

// SplitYAMLDocument is a bufio.SplitFunc for splitting a YAML document into individual documents.
//
// This is from Kubernetes' 'pkg/util/yaml'.splitYAMLDocument, which is unfortunately
// not exported.
func SplitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
