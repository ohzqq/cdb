package cdb

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// Serializer holds the options for encoding or decoding a book.
type Serializer struct {
	Indent   int
	Format   string
	Editable bool
	Book     *Book
	encoder  BookEncoder
	decoder  BookDecoder
}

// SerializerOpt sets the options for a serializer.
type SerializerOpt func(*Serializer)

// Encoder is an interface for types that can Encode a book.
type Encoder interface {
	Encode(v any) error
}

// EncoderInit takes the serializer options and returns a function that
// initializes an Encoder.
type EncoderInit func(*Serializer) BookEncoder

// BookEncoder takes a io.Writer and returns a type that can Encode a book.
type BookEncoder func(w io.Writer) Encoder

// Decoder is an interface for types that can decode a book.
type Decoder interface {
	Decode(v any) error
}

// DecoderInit takes the serializer options and returns a function that
// initializes a Decoder.
type DecoderInit func(*Serializer) BookDecoder

// BookDecoder takes a io.Writer and returns a type that can Decode a book.
type BookDecoder func(r io.Reader) Decoder

// NewSerializer constructs a serializer for a book with default options.
func NewSerializer(b *Book, opts ...SerializerOpt) *Serializer {
	s := &Serializer{
		Book:   b,
		Indent: 2,
		Format: ".txt",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Decoder takes a DecoderInit and sets the BookDecoder.
func (s *Serializer) Decoder(init DecoderInit) *Serializer {
	s.decoder = init(s)
	return s
}

// Encoder takes a EncoderInit and sets the BookEncoder.
func (s *Serializer) Encoder(init EncoderInit) *Serializer {
	s.encoder = init(s)
	return s
}

// WithIndent sets the number of spaces for indenting.
func WithIndent(n int) SerializerOpt {
	return func(s *Serializer) {
		s.Indent = n
	}
}

// EditableOnly sets the option for serializing only editable book fields.
func EditableOnly() SerializerOpt {
	return func(s *Serializer) {
		s.Editable = true
	}
}

// WithFormat sets the encoding format extension.
func WithFormat(ext string) SerializerOpt {
	return func(s *Serializer) {
		s.Format = ext
	}
}

// DecodeYAML configures a YAML BookDecoder.
func DecodeYAML(s *Serializer) BookDecoder {
	s.Format = ".yaml"
	return func(r io.Reader) Decoder {
		dec := yaml.NewDecoder(r)
		return dec
	}
}

// DecodeJSON configures a JSON BookDecoder.
func DecodeJSON(s *Serializer) BookDecoder {
	s.Format = ".json"
	return func(r io.Reader) Decoder {
		dec := json.NewDecoder(r)
		return dec
	}
}

// DecodeTOML configures a TOML BookDecoder.
func DecodeTOML(s *Serializer) BookDecoder {
	s.Format = ".toml"
	return func(r io.Reader) Decoder {
		dec := toml.NewDecoder(r)
		return dec
	}
}

// ReadFile reads a file for decoding.
func (s *Serializer) ReadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.ReadFrom(file)
}

// ReadFrom reads from io.Reader for decoding.
func (s *Serializer) ReadFrom(r io.Reader) error {
	d := s.decoder(r)
	return d.Decode(s.Book)
}

// EncodeYAML configures a YAML BookEncoder.
func EncodeYAML(s *Serializer) BookEncoder {
	s.Format = ".yaml"
	return func(w io.Writer) Encoder {
		enc := yaml.NewEncoder(w)
		if s.Indent > 0 {
			enc.SetIndent(s.Indent)
		}
		return enc
	}
}

// EncodeJSON configures a JSON BookEncoder.
func EncodeJSON(s *Serializer) BookEncoder {
	s.Format = ".json"
	return func(w io.Writer) Encoder {
		enc := json.NewEncoder(w)
		if s.Indent > 0 {
			enc.SetIndent("", strings.Repeat(" ", s.Indent))
		}
		return enc
	}
}

// EncodeTOML configures a TOML BookEncoder.
func EncodeTOML(s *Serializer) BookEncoder {
	s.Format = ".toml"
	return func(w io.Writer) Encoder {
		enc := toml.NewEncoder(w)
		if s.Indent > 0 {
			enc.SetIndentSymbol(strings.Repeat(" ", s.Indent))
		}
		return enc
	}
}

// WriteFile writes an encoded book to disk.
func (s *Serializer) WriteFile(name string) error {
	file, err := os.Create(name + s.Format)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.WriteTo(file)
}

// WriteTo writes an encoded book to io.Writer.
func (s *Serializer) WriteTo(w io.Writer) error {
	enc := s.encoder(w)

	if s.Editable {
		return enc.Encode(s.Book.EditableFields)
	}

	return enc.Encode(s.Book)
}
