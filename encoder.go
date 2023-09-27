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

type EncoderInit func(*EncoderConfig) BookEncoder
type EncoderOpt func(*EncoderConfig)
type BookEncoder func(w io.Writer) Encoder

type EncoderConfig struct {
	indent   int
	Format   string
	editable bool
	book     *Book
	encoder  BookEncoder
	decoder  BookDecoder
}

func (s *EncoderConfig) Decoder(init DecoderInit) *EncoderConfig {
	s.decoder = init(s)
	return s
}

func (s *EncoderConfig) Encoder(init EncoderInit) *EncoderConfig {
	s.encoder = init(s)
	return s
}

func (s *EncoderConfig) ReadFrom(r io.Reader) error {
	d := s.decoder(r)
	return d.Decode(s.book)
}

func NewEncoder(b *Book, opts ...EncoderOpt) *EncoderConfig {
	enc := &EncoderConfig{
		book:   b,
		indent: 2,
		Format: ".txt",
	}

	for _, opt := range opts {
		opt(enc)
	}

	return enc
}

func WithIndent(n int) EncoderOpt {
	return func(enc *EncoderConfig) {
		enc.indent = n
	}
}

func EditableOnly() EncoderOpt {
	return func(enc *EncoderConfig) {
		enc.editable = true
	}
}

func WithFormat(ext string) EncoderOpt {
	return func(enc *EncoderConfig) {
		enc.Format = ext
	}
}

func EncodeYAML(e *EncoderConfig) BookEncoder {
	e.Format = ".yaml"
	return func(w io.Writer) Encoder {
		enc := yaml.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndent(e.indent)
		}
		return enc
	}
}

func EncodeJSON(e *EncoderConfig) BookEncoder {
	e.Format = ".json"
	return func(w io.Writer) Encoder {
		enc := json.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndent("", strings.Repeat(" ", e.indent))
		}
		return enc
	}
}

func EncodeTOML(e *EncoderConfig) BookEncoder {
	e.Format = ".toml"
	return func(w io.Writer) Encoder {
		enc := toml.NewEncoder(w)
		if e.indent > 0 {
			enc.SetIndentSymbol(strings.Repeat(" ", e.indent))
		}
		return enc
	}
}

func (e *EncoderConfig) WriteFile(name string) error {
	file, err := os.Create(name + e.Format)
	if err != nil {
		return err
	}
	defer file.Close()

	return e.WriteTo(file)
}

func (e *EncoderConfig) WriteTo(w io.Writer) error {
	enc := e.encoder(w)

	if e.editable {
		return enc.Encode(e.book.EditableFields)
	}

	return enc.Encode(e.book)
}
