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

type Serialize struct {
	ext      string
	editable bool
	rw       io.ReadWriter
}

func NewSerializer(rw io.ReadWriter, ext string) *Serialize {
	return &Serialize{
		ext: ext,
		rw:  rw,
	}
}

func (e *Serialize) Editable() *Serialize {
	e.editable = true
	return e
}

func (e *Serialize) EncodeBook(enc BookEncoder, b Book) error {
	if e.editable {
		return enc.Encode(b.editableStringMap())
	}
	return enc.Encode(b)
}

func (e *Serialize) DecodeBook(enc BookDecoder, b *Book) error {
	return enc.Decode(b)
}

func (e *Serialize) Encode(b Book) error {
	var enc BookEncoder
	switch e.ext {
	case ".yaml", ".yml":
		yenc := yaml.NewEncoder(e.rw)
		yenc.SetIndent(2)
		enc = yenc
	case ".toml":
		tenc := toml.NewEncoder(e.rw)
		tenc.Indent = "  "
		enc = tenc
	case ".json":
		jenc := json.NewEncoder(e.rw)
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

func (e *Serialize) Decode(b *Book) error {
	var enc BookDecoder
	switch e.ext {
	case ".yaml", ".yml":
		enc = yaml.NewDecoder(e.rw)
	case ".toml":
		buf := new(bytes.Buffer)
		buf.ReadFrom(e.rw)
		return toml.Unmarshal(buf.Bytes(), b)
	case ".json":
		enc = json.NewDecoder(e.rw)
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
