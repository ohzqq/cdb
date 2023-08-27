package cdb

import (
	"log"
)

// SetMetadata sets book metadata
func SetMetadata(lib, pos string, meta map[string]string, args ...Opt) *Command {
	cmd, err := NewCommand(lib, args...)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Opt(Cmd("set_metadata"), PositionalArgs(pos), Fields(meta))
	return cmd
}

func Fields(val map[string]string) Opt {
	var fields []string
	for _, k := range AllModels().Editable() {
		if v, ok := val[k]; ok {
			switch k {
			case SeriesIndex:
				if v != "0" || v != "0.0" {
					fields = append(fields, "--field", k+":"+v)
				}
				continue
			}
			fields = append(fields, "--field", k+":"+v)
		}
	}
	return Flags(fields...)
}
