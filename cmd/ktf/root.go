package ktf

import (
	"fmt"
	"io/ioutil"
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
		// YAML string stored in a variable
		yamlString, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		// Map to store the parsed YAML data
		var data map[string]interface{}

		// Unmarshal the YAML string into the data map
		marshallErr := yaml.Unmarshal([]byte(yamlString), &data)
		if marshallErr != nil {
			fmt.Println(err)
		}

		code, err := ktf.BuildManifest(data)

		fmt.Printf("%s", code)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
