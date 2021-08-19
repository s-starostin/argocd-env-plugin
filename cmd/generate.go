package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yaml "sigs.k8s.io/yaml"
)

type Config struct {
	includePaths []string
	path         string
}

func getConfig(cmd *cobra.Command) *Config {
	const delimiter = ":"
	path, _ := cmd.Flags().GetString("path")
	includePaths, _ := cmd.Flags().GetString("include-paths")

	return &Config{
		includePaths: func() []string {
			if len(includePaths) > 0 {
				return strings.Split(includePaths, delimiter)
			} else {
				return getEnvAsSlice("AEP_INCLUDE_PATHS", []string{}, delimiter)
			}
		}(),
		path: func() string {
			if len(path) > 0 {
				return path
			} else {
				return getEnv("AEP_PATH", "-")
			}
		}(),
	}
}

// NewGenerateCommand initializes the generate command
func NewGenerateCommand() *cobra.Command {
	const StdIn = "-"

	var command = &cobra.Command{
		Use:   "generate --path=<path> --include-paths=<additional path>:<additional path>",
		RunE: func(cmd *cobra.Command, args []string) error {
			var manifests []unstructured.Unstructured
			var err error
			var errs []error

			conf := getConfig(cmd)

			if conf.path == StdIn {
				manifests, err = readManifestData(cmd.InOrStdin())
				if err != nil {
					return err
				}
			} else {
				files, err := listYamlFiles(conf.path)
				if len(files) < 1 {
					return fmt.Errorf("no YAML files were found in %s", conf.path)
				}
				if err != nil {
					return err
				}

				manifests, errs = readFilesAsManifests(files)
				if len(errs) != 0 {
					// TODO: handle multiple errors nicely
					return fmt.Errorf("could not read YAML files: %s", errs)
				}
			}

			if len(conf.includePaths) > 0 {
				var addManifests []unstructured.Unstructured

				for _, path := range conf.includePaths {
					files, err := listYamlFiles(path)
					if len(files) < 1 {
						continue
					}
					if err != nil {
						return err
					}
					addManifests, errs = readFilesAsManifests(files)
					if len(errs) != 0 {
						return fmt.Errorf("could not read YAML files: %s", errs)
					}
					manifests = append(manifests, addManifests...)
				}
			}

			for _, manifest := range manifests {

				if len(manifest.Object) == 0 {
					continue
				}

				if err != nil {
					return err
				}

				output, err := yaml.Marshal(manifest.Object)
				if err != nil {
					fmt.Errorf("could not export %s into YAML: %s", err)
				}

				fmt.Fprintf(cmd.OutOrStdout(), "%s---\n", string(output))
			}

			return nil
		},
	}

	return command
}
