package gen

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/danielgtaylor/casing"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type commands map[string]map[string]string

type command struct {
	Args   string
	Params string
	Flags  map[string]string
}

func ReadCommands() commands {
	cmds, err := os.ReadFile("internal/gen/calibredb_commands.yaml")
	if err != nil {
		log.Fatal(err)
	}

	meta := make(commands)
	err = yaml.Unmarshal(cmds, meta)
	if err != nil {
		log.Fatal(err)
	}

	return meta
}

func (c commands) Commands() []string {
	return maps.Keys(c)
}

func (c commands) GetCommand(name string) command {
	cmd := command{
		Flags: make(map[string]string),
	}
	for k, v := range c[name] {
		if k == "Args" {
			aa := strings.Split(v, ",")
			var params []string
			var vars []string
			for _, a := range aa {
				args := strings.Split(a, "=")
				n := strings.ToLower(args[0])
				t := args[1]
				if t == "[]string" {
					t = "...string"
				}
				params = append(params, n+" "+t)
				if t == "...string" {
					n += "..."
				}
				vars = append(vars, n)
			}
			cmd.Params = strings.Join(params, ", ")
			cmd.Args = strings.Join(vars, ", ")
			continue
		}
		cmd.Flags[k] = v
	}
	return cmd
}

func (c commands) CommandsString() string {
	fn := func(cmd string) string {
		var b strings.Builder

		b.WriteByte('"')
		if strings.HasPrefix(cmd, "SavedSearches") {
			b.WriteString("saved_searches ")
			b.WriteString(strings.ToLower(strings.TrimPrefix(cmd, "SavedSearches")))
		} else {
			b.WriteString(casing.Snake(cmd))
		}

		b.WriteByte('"')
		return b.String()
	}

	return c.CommandsSwitch("Cmd", "string", fn)
}

func (c commands) CommandsSwitch(cmd string, rt string, cr func(string) string) string {
	cmds := c.Commands()
	var b strings.Builder

	b.WriteString("func (c CaldbCmd) " + cmd + "() " + rt + " {")
	b.WriteByte('\n')
	b.WriteString("switch c {")
	b.WriteByte('\n')

	for _, cmd := range cmds {
		b.WriteString("case " + cmd + ":")
		b.WriteByte('\n')
		b.WriteString(`return ` + cr(cmd))
		b.WriteByte('\n')
	}

	b.WriteString("}\n")
	b.WriteString(`return ""`)
	b.WriteString("\n}\n")

	return b.String()
}

func (c commands) CommandFuncs() string {
	var buf bytes.Buffer
	err := cmdFunc.Execute(&buf, c)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func (c commands) CommandList() string {
	var buf bytes.Buffer
	err := cmdList.Execute(&buf, c)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func (c commands) CommandBuilder() string {
	var buf bytes.Buffer
	err := cmdBuilder.Execute(&buf, c)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

var tmplFuncs = template.FuncMap{
	"snake": lowerSnake,
}

func lowerSnake(v string) string {
	if strings.HasPrefix(v, "SavedSearches") {
		return "saved_searches " + strings.ToLower(strings.TrimPrefix(v, "SavedSearches"))
	}
	return casing.Snake(v, strings.ToLower)
}

var cmdList = template.Must(
	template.New("cmdList").
		Funcs(tmplFuncs).
		Parse(`
func ListCommands() []string {
return []string{
{{range $name := .Commands -}}
"{{snake .}}",
{{end -}}
}
}
`))

var cmdBuilder = template.Must(
	template.New("cmdBuilder").
		Funcs(tmplFuncs).
		Parse(`
{{- $commands := . -}}
{{range $name := .Commands -}}
{{$cmd := $commands.GetCommand . -}}

type {{.}} struct {
	*Command
}

{{range $flag, $val := $cmd.Flags}}
func (c *{{$name}}) {{$flag -}} 
	({{with ne $val "bool"}}v {{$val}}{{end}}) *{{$name}} {
	c.SetFlags("--{{snake $flag}}"{{with ne $val "bool"}}, v{{end}})
	return c
}	
	{{end -}}
{{end -}}
`))

var cmdFunc = template.Must(
	template.New("cmdFunc").
		Funcs(tmplFuncs).
		Parse(`
{{- $commands := . -}}
{{range $name := .Commands -}}
{{$cmd := $commands.GetCommand .}}

func (c *Command) {{.}} ({{- with $cmd.Params -}}{{.}}{{- end -}})  *{{.}} {
	{{with $cmd.Args -}}
		c.SetPositional({{.}})
	{{end -}}
	c.CdbCmd = "{{snake .}}"
	cmd := &{{.}}{
		Command: c,
	}
	return cmd
}
{{end -}}
`))

func (c commands) CommandsConst() string {
	cmds := c.Commands()
	var b strings.Builder

	b.WriteString("const (\n")
	b.WriteString(cmds[0])
	b.WriteString(" CaldbCmd = iota\n")

	for _, cmd := range cmds[1:] {
		b.WriteString(cmd)
		b.WriteByte('\n')
	}

	b.WriteString(")\n")

	return b.String()
}
