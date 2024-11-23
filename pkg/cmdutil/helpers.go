package cmdutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	golog "log"
	"os"
	"path"

	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sweetloveinyourheart/planning-poker/pkg/config"
)

func generateSchema(rootCmd *cobra.Command) {
	outJSON := new(bytes.Buffer)
	err := genSchema(rootCmd, outJSON, "json")
	if err != nil {
		golog.Fatal(err)
	}

	outYAML := new(bytes.Buffer)
	err = genSchema(rootCmd, outYAML, "yaml")
	if err != nil {
		golog.Fatal(err)
	}

	var paths []string

	if _, err := os.Stat("./go.mod"); err == nil {
		paths = []string{
			path.Join(".", "cmd", "planning-pocker", "schema"),
		}
	} else {
		paths = []string{
			path.Join(".", "schema"),
		}
	}

	for _, fp := range paths {
		filePath := fmt.Sprintf("%s.%s", fp, "json")
		err = os.WriteFile(filePath, outJSON.Bytes(), 0644)
		if err != nil {
			golog.Fatal(err)
		}

		filePath = fmt.Sprintf("%s.%s", fp, "yaml")
		err = os.WriteFile(filePath, outYAML.Bytes(), 0644)
		if err != nil {
			golog.Fatal(err)
		}
	}
}
func genSchema(rootCmd *cobra.Command, w io.Writer, format string) error {
	cmds := append([]*cobra.Command{rootCmd}, rootCmd.Commands()...)

	type schemaConfig struct {
		Name         string   `json:"name" yaml:"name"`
		Usage        string   `json:"usage" yaml:"usage"`
		DefaultValue any      `json:"default" yaml:"default"`
		ValueType    string   `json:"valueType" yaml:"valueType"`
		Path         string   `json:"path" yaml:"path"`
		Env          []string `json:"env" yaml:"env"`
	}

	type schemaService struct {
		Name      string `json:"name" yaml:"name"`
		ShortName string `json:"shortName" yaml:"shortName"`
		Long      string `json:"long" yaml:"long"`
		config.Service
		Config []schemaConfig
	}

	services := make([]*schemaService, 0)
	for _, cmd := range cmds {
		cmd.InitDefaultHelpCmd()
		cmd.InitDefaultHelpFlag()

		buf := new(bytes.Buffer)

		if cmd.PreRunE != nil {
			_ = cmd.PreRunE(cmd, []string{})
		}

		service, ok := config.GetService(cmd)
		if !ok {
			continue
		}

		ss := &schemaService{
			Name:      cmd.Name(),
			ShortName: cmd.Short,
			Long:      cmd.Long,
			Service:   service,
		}
		services = append(services, ss)

		flags := cmd.NonInheritedFlags()
		flags.SetOutput(buf)
		if flags.HasAvailableFlags() {
			flags.VisitAll(func(flag *pflag.Flag) {
				binding, ok := config.Registry[flag]
				if !ok {
					return
				}
				ss.Config = append(ss.Config, schemaConfig{
					Name:         binding.Flag.Name,
					Usage:        binding.Usage,
					DefaultValue: binding.Value,
					ValueType:    fmt.Sprintf("%T", binding.Value),
					Path:         binding.Aliases[0],
					Env:          binding.Env,
				})
			})
		}

		parentFlags := cmd.InheritedFlags()
		parentFlags.SetOutput(buf)
		if parentFlags.HasAvailableFlags() {
			parentFlags.VisitAll(func(flag *pflag.Flag) {
				binding, ok := config.Registry[flag]
				if !ok {
					return
				}
				ss.Config = append(ss.Config, schemaConfig{
					Name:         binding.Flag.Name,
					Usage:        binding.Usage,
					DefaultValue: binding.Value,
					ValueType:    fmt.Sprintf("%T", binding.Value),
					Path:         binding.Aliases[0],
					Env:          binding.Env,
				})
			})
		}
	}
	switch format {
	case "yaml":
		yaml, err := yaml.Marshal(services)
		if err != nil {
			return err
		}
		if _, err := w.Write(yaml); err != nil {
			return err
		}
	case "json":
		json, err := json.MarshalIndent(services, "", "  ")
		if err != nil {
			return err
		}
		if _, err := w.Write(json); err != nil {
			return err
		}
	}
	return nil
}
