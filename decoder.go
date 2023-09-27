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

type Decoder struct {
	book   Book
	reader io.Reader
	BookDecoder
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
