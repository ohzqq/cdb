package calibredb

import (
	"github.com/ohzqq/cdb"
	"github.com/spf13/viper"
)

// SetMetadata sets book metadata
func SetMetadata(lib, pos string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("set_metadata"), PositionalArgs(pos))
	return cmd
}

func Fields(val map[string]string) Opt {
	var fields []string

	models := cdb.DefaultModels()
	if cdb.GetLib(viper.GetString("lib")).Audiobooks {
		models = cdb.AudiobookModels()
	}

	for _, k := range models.Editable() {
		if v, ok := val[k]; ok {
			switch k {
			case cdb.SeriesIndex:
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
