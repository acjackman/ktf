package ktf

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/acjackman/ktf/pkg/ktf"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "ktf",
	Short: "ktf - create terraform manifest resources from yaml",
	Long: `ktf helps generate terraform resources for kubernetes manifests

Pass a terraform manifest`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// this does the trick
		var inputReader io.Reader = cmd.InOrStdin()

		// the argument received looks like a file, we try to open it
		if len(args) > 0 && args[0] != "-" {
			file, err := os.Open(args[0])
			if err != nil {
				log.Printf("failed open file: %v", err)
				os.Exit(1)
			}
			inputReader = file
		}

		yamlString, err := io.ReadAll(inputReader)
		if err != nil {
			log.Printf("failed to read: %v", err)
			os.Exit(1)
		}

		// Map to store the parsed YAML data
		var data map[string]interface{}

		// Unmarshal the YAML string into the data map
		marshallErr := yaml.Unmarshal([]byte(yamlString), &data)
		if marshallErr != nil {
			log.Printf("Unable to unmarshal yaml: %v", err)
			os.Exit(1)
		}

		code, err := ktf.BuildManifest(data)
		if err != nil {
			log.Printf("Error building HCL: %v", err)
			os.Exit(1)
		}

		fmt.Printf("%s", code)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
