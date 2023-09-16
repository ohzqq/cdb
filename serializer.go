package cdb

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type BookEncoder interface {
	Encode(v any) error
}

type BookDecoder interface {
	Decode(v any) error
}

func (b Book) Encode(enc BookEncoder) error {
	err := enc.Encode(b)
	return err
}

func (b Book) Decode(enc BookDecoder) error {
	err := enc.Decode(b)
	return err
}

func (b Book) Save(ext string) error {
	file, err := os.Create(b.Title + ext)
	if err != nil {
		return err
	}

	var enc BookEncoder
	switch ext {
	case ".yaml", ".yml":
		yenc := yaml.NewEncoder(file)
		yenc.SetIndent("  ")
		enc = yenc
	case ".toml":
		tenc := toml.NewEncoder(file)
		tenc.Indent = "  "
		enc = tenc
	case ".json":
		jenc := json.NewEncoder(file)
		jenc.SetIndent("", "  ")
		enc = jenc
	default:
		return fmt.Errorf("only yaml, toml, and json can be written\n")
	}

	err = b.Encode(enc)
	if err != nil {
		return err
	}

	return nil
}
