package cdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type BookEncoder interface {
	Encode(v any) error
}

type BookDecoder interface {
	Decode(v any) error
}

type Encoder struct {
	writer io.Writer
	indent bool
	BookEncoder
}

type Decoder struct {
	reader io.Reader
	BookDecoder
}

func NewEncoder(w io.Writer) *Encoder {
	enc := &Encoder{
		writer: w,
	}
	enc.Format(".json")
	return enc
}

func (e *Encoder) Format(ext string) *Encoder {
	switch ext {
	case ".yaml", ".yml":
		yenc := yaml.NewEncoder(e.writer)
		if e.indent {
			yenc.SetIndent(2)
		}
		e.SetEncoder(yenc)

	case ".toml":
		tenc := toml.NewEncoder(e.writer)
		e.SetEncoder(tenc)

	case ".json":
		jenc := json.NewEncoder(e.writer)
		if e.indent {
			jenc.SetIndent("", "  ")
		}
		e.SetEncoder(jenc)
	}

	return e
}

func (enc *Encoder) SetEncoder(b BookEncoder) *Encoder {
	enc.BookEncoder = b
	return enc
}

func NewDecoder(r io.Reader) *Decoder {
	dec := &Decoder{
		reader: r,
	}
	return dec
}

func ReadMetadataFile(path string) (*Decoder, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return &Decoder{}, err
	}
	buf := bytes.NewBuffer(data)

	d := NewDecoder(buf)
	ext := filepath.Ext(path)
	if ext != ".yaml" && ext != ".yml" && ext != ".json" && ext != ".toml" {
		return d, fmt.Errorf("set custom decoder")
	}
	d.Format(ext)

	return d, nil
}

func (e *Decoder) Format(ext string) *Decoder {
	switch ext {
	case ".yaml", ".yml":
		return e.SetDecoder(yaml.NewDecoder(e.reader))
	case ".toml":
		return e.SetDecoder(toml.NewDecoder(e.reader))
	case ".json":
		return e.SetDecoder(json.NewDecoder(e.reader))
	default:
		return e
	}
}

func (d *Decoder) SetDecoder(bd BookDecoder) *Decoder {
	d.BookDecoder = bd
	return d
}
