package cdb

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/danielgtaylor/casing"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

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
