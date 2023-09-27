package cdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielgtaylor/casing"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type BookEncoder interface {
	Encode(v any) error
}

type EncoderConfig func(*Encoder) EncoderFunc
type EncoderOpt func(*Encoder)
type EncoderFunc func(w io.Writer) BookEncoder

type Encoder struct {
	indent   int
	format   string
	editable bool
	book     Book
	init     EncoderConfig
}

type BookDecoder interface {
	Decode(v any) error
}

type Decoder struct {
	book   Book
	reader io.Reader
	BookDecoder
}

func NewEncoder(b Book, init EncoderConfig, opts ...EncoderOpt) *Encoder {
	enc := &Encoder{
		book:   b,
		indent: 2,
		init:   init,
		format: ".txt",
	}

	for _, opt := range opts {
		opt(enc)
	}

	return enc
}

func WithIndent(n int) EncoderOpt {
	return func(enc *Encoder) {
		enc.indent = n
	}
}

func EditableOnly() EncoderOpt {
	return func(enc *Encoder) {
		enc.editable = true
	}
}

func WithFormat(ext string) EncoderOpt {
	return func(enc *Encoder) {
		enc.format = ext
	}
}

func EncodeYAML(e *Encoder) EncoderFunc {
	return func(w io.Writer) BookEncoder {
		e.format = ".yaml"
		enc := yaml.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndent(e.indent)
		}
		return enc
	}
}

func EncodeJSON(e *Encoder) EncoderFunc {
	return func(w io.Writer) BookEncoder {
		e.format = ".json"
		enc := json.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndent("", strings.Repeat(" ", e.indent))
		}
		return enc
	}
}

func EncodeTOML(e *Encoder) EncoderFunc {
	return func(w io.Writer) BookEncoder {
		enc := toml.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndentSymbol(strings.Repeat(" ", e.indent))
		}
		e.format = ".toml"
		return enc
	}
}

func (e *Encoder) WriteTo(w io.Writer) error {
	fn := e.init(e)
	enc := fn(w)

	if e.editable {
		return enc.Encode(e.book.EditableFields)
	}

	return enc.Encode(e.book)
}

func (e *Encoder) Save(name string) error {
	if name == "" {
		name = casing.Snake(e.book.Title)
	}

	file, err := os.Create(name + e.format)
	if err != nil {
		return err
	}
	defer file.Close()

	return e.WriteTo(file)
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
