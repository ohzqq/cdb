package calibredb

import (
	"github.com/ohzqq/cdb"
)

// SetMetadata sets book metadata
func SetMetadata(lib, pos string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("set_metadata"), PositionalArgs(pos))
	return cmd
}

func Fields(val map[string]string) Opt {
	var fields []string
	for _, k := range allowedFields {
		if v, ok := val[k]; ok {
			switch k {
			case cdb.Cover:
				continue
			case cdb.Identifiers:
				//id = v
				continue
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

var allowedFields = []string{
	"authors",
	"author_sort",
	"comments",
	"identifiers",
	"languages",
	"pubdate",
	"publisher",
	"rating",
	"series",
	"series_index",
	"sort",
	"tags",
	"timestamp",
	"title",
	"#narrators",
	"#duration",
}
