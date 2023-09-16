package cdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type BookEncoder interface {
	Encode(v any) error
}

type BookDecoder interface {
	Decode(v any) error
}

type Encode struct {
	ext      string
	editable bool
	enc      BookEncoder
	book     Book
	writer   io.ReadWriter
}

func NewSerializer(w io.ReadWriter, ext string) *Encode {
	return &Encode{
		ext:    ext,
		writer: w,
	}
}

func (e *Encode) Editable() *Encode {
	e.editable = true
	return e
}

func (e *Encode) Encode(b Book) error {
	var enc BookEncoder
	switch e.ext {
	case ".yaml", ".yml":
		yenc := yaml.NewEncoder(e.writer)
		yenc.SetIndent(2)
		enc = yenc
	case ".toml":
		tenc := toml.NewEncoder(e.writer)
		tenc.Indent = "  "
		enc = tenc
	case ".json":
		jenc := json.NewEncoder(e.writer)
		jenc.SetIndent("", "  ")
		enc = jenc
	default:
		return fmt.Errorf("only yaml, toml, and json can be written\n")
	}

	if e.editable {
		return enc.Encode(b.editableStringMap())
	}

	return enc.Encode(b)
}

func (e *Encode) Decode(b *Book) error {
	var enc BookDecoder
	switch e.ext {
	case ".yaml", ".yml":
		enc = yaml.NewDecoder(e.writer)
	case ".toml":
		buf := new(bytes.Buffer)
		buf.ReadFrom(e.writer)
		return toml.Unmarshal(buf.Bytes(), b)
	//enc = toml.NewDecoder(e.writer)
	case ".json":
		enc = json.NewDecoder(e.writer)
	default:
		return fmt.Errorf("only yaml, toml, and json can be written\n")
	}

	return enc.Decode(b)
}

func (b *Book) ReadFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	ser := NewSerializer(file, filepath.Ext(name))
	return ser.Decode(b)
}

func (b Book) Print(ext string, editable bool) error {
	enc := NewSerializer(os.Stdout, ext)
	if editable {
		enc.Editable()
	}
	return enc.Encode(b)
}

func (b Book) Save(ext string, editable bool) error {
	file, err := os.Create(b.Title + ext)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := NewSerializer(file, ext)
	if editable {
		enc.Editable()
	}
	return enc.Encode(b)
}
