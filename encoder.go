package cdb

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type Encoder interface {
	Encode(v any) error
}

type EncoderInit func(*Serializer) BookEncoder
type SerializerOpt func(*Serializer)
type BookEncoder func(w io.Writer) Encoder

type Serializer struct {
	indent   int
	Format   string
	editable bool
	book     *Book
	encoder  BookEncoder
	decoder  BookDecoder
}

func (s *Serializer) Decoder(init DecoderInit) *Serializer {
	s.decoder = init(s)
	return s
}

func (s *Serializer) Encoder(init EncoderInit) *Serializer {
	s.encoder = init(s)
	return s
}

func (s *Serializer) ReadFrom(r io.Reader) error {
	d := s.decoder(r)
	return d.Decode(s.book)
}

func NewEncoder(b *Book, opts ...SerializerOpt) *Serializer {
	enc := &Serializer{
		book:   b,
		indent: 2,
		Format: ".txt",
	}

	for _, opt := range opts {
		opt(enc)
	}

	return enc
}

func WithIndent(n int) SerializerOpt {
	return func(enc *Serializer) {
		enc.indent = n
	}
}

func EditableOnly() SerializerOpt {
	return func(enc *Serializer) {
		enc.editable = true
	}
}

func WithFormat(ext string) SerializerOpt {
	return func(enc *Serializer) {
		enc.Format = ext
	}
}

func EncodeYAML(e *Serializer) BookEncoder {
	e.Format = ".yaml"
	return func(w io.Writer) Encoder {
		enc := yaml.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndent(e.indent)
		}
		return enc
	}
}

func EncodeJSON(e *Serializer) BookEncoder {
	e.Format = ".json"
	return func(w io.Writer) Encoder {
		enc := json.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndent("", strings.Repeat(" ", e.indent))
		}
		return enc
	}
}

func EncodeTOML(e *Serializer) BookEncoder {
	e.Format = ".toml"
	return func(w io.Writer) Encoder {
		enc := toml.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndentSymbol(strings.Repeat(" ", e.indent))
		}
		return enc
	}
}

func (e *Serializer) WriteFile(name string) error {
	file, err := os.Create(name + e.Format)
	if err != nil {
		return err
	}
	defer file.Close()

	return e.WriteTo(file)
}

func (e *Serializer) WriteTo(w io.Writer) error {
	enc := e.encoder(w)

	if e.editable {
		return enc.Encode(e.book.EditableFields)
	}

	return enc.Encode(e.book)
}
