package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Revi-Studios/project/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	targetDir := "./md"
	templatePath := "./docs.tmpl"

	if err := os.RemoveAll(targetDir); err != nil {
		log.Fatalf("Failed to delete directory: %v", err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	tmpl, err := template.New("docs.tmpl").Funcs(template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}).ParseFiles(templatePath)

	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	data := cobraToCmdDataTree(cmd.Root, nil)

	if err := generateBranchTree(data, tmpl, targetDir); err != nil {
		log.Fatalf("Failed to generate documentation: %v", err)
	}
}

type CmdData struct {
	Name                    string
	Short                   string
	Long                    string
	MdFile                  string
	Example                 string
	CommandPath             string
	Parent                  *CmdData
	Use                     string
	Aliases                 []string
	LocalFlags              []FlagData
	SubCommands             []*CmdData
	HasLong                 bool
	HasExample              bool
	HasParent               bool
	HasAvailableLocalFlags  bool
	HasAvailableSubCommands bool
}

type FlagData struct {
	Name      string
	Shorthand string
	Usage     string
	DefValue  string
}

func generateBranchTree(cmd *CmdData, tmpl *template.Template, dir string) error {
	filename := filepath.Join(dir, cmd.MdFile)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, cmd); err != nil {
		return err
	}

	for _, sub := range cmd.SubCommands {
		if err := generateBranchTree(sub, tmpl, dir); err != nil {
			return err
		}
	}

	return nil
}

func cobraToCmdData(cmd *cobra.Command) *CmdData {
	data := &CmdData{
		Name:                    cmd.Name(),
		Short:                   cmd.Short,
		Long:                    cmd.Long,
		MdFile:                  toMdFileName(cmd.CommandPath()),
		Example:                 cmd.Example,
		CommandPath:             cmd.CommandPath(),
		Use:                     cmd.UseLine(),
		Aliases:                 cmd.Aliases,
		HasLong:                 cmd.Long != "",
		HasExample:              cmd.Example != "",
		HasParent:               false,
		HasAvailableLocalFlags:  cmd.HasAvailableLocalFlags(),
		HasAvailableSubCommands: cmd.HasAvailableSubCommands(),
	}

	return data
}

func cobraToCmdDataTree(cmd *cobra.Command, parent *CmdData) *CmdData {
	data := &CmdData{
		Name:                    cmd.Name(),
		Short:                   cmd.Short,
		Long:                    cmd.Long,
		MdFile:                  toMdFileName(cmd.CommandPath()),
		Example:                 cmd.Example,
		CommandPath:             cmd.CommandPath(),
		Parent:                  parent,
		Use:                     cmd.UseLine(),
		Aliases:                 cmd.Aliases,
		HasLong:                 cmd.Long != "",
		HasExample:              cmd.Example != "",
		HasParent:               parent != nil,
		HasAvailableLocalFlags:  cmd.HasAvailableLocalFlags(),
		HasAvailableSubCommands: cmd.HasAvailableSubCommands(),
	}

	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		data.LocalFlags = append(data.LocalFlags, FlagData{
			Name:      f.Name,
			Shorthand: f.Shorthand,
			Usage:     f.Usage,
			DefValue:  f.DefValue,
		})
	})

	for _, sub := range cmd.Commands() {
		if !sub.IsAvailableCommand() || sub.Name() == "help" {
			continue
		}
		data.SubCommands = append(data.SubCommands, cobraToCmdDataTree(sub, data))
	}

	return data
}

func toMdFileName(str string) string {
	str = strings.ReplaceAll(str, " ", "-")
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	str += ".md"

	return str
}
