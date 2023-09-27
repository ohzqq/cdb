package cdb

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type Decoder interface {
	Decode(v any) error
}

type BookDecoder func(r io.Reader) Decoder
type DecoderInit func(*Serializer) BookDecoder

type DecoderConfig struct {
	book    *Book
	Ext     string
	decoder BookDecoder
	Decoder
}

func DecodeYAML(e *Serializer) BookDecoder {
	e.Format = ".yaml"
	return func(r io.Reader) Decoder {
		dec := yaml.NewDecoder(r)
		return dec
	}
}

func DecodeJSON(e *Serializer) BookDecoder {
	e.Format = ".json"
	return func(r io.Reader) Decoder {
		dec := json.NewDecoder(r)
		return dec
	}
}

func DecodeTOML(e *Serializer) BookDecoder {
	e.Format = ".toml"
	return func(r io.Reader) Decoder {
		dec := toml.NewDecoder(r)
		return dec
	}
}

func (s *Serializer) ReadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.ReadFrom(file)
}
