package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keptn/lifecycle-toolkit/cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const outputFileName = "output"
const providerName = "provider"

func init() {
	flags := convertSLICmd.Flags()

	flags.StringP(outputFileName, "o", "manifests.yaml", "Set the output file name")
	flags.StringP(providerName, "p", "dynatrace", "Set the provider")

	_ = viper.BindPFlag(outputFileName, flags.Lookup(outputFileName))
	_ = viper.BindPFlag(providerName, flags.Lookup(providerName))
}

type SLIContent struct {
	Indicators map[string]string `yaml:"indicators"`
}

// startCmd represents the start command
var convertSLICmd = &cobra.Command{
	Use:   "convert-sli",
	Short: "Convert SLI to AnalysisTemplate",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("input file must be set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Converting...")

		// Open the file for reading
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println("Error opening the file:", err)
			return
		}
		defer file.Close() // Close the file when we're done with it

		// Read the file content
		content, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading the file:", err)
			return
		}

		// Convert the content to a string and print it
		//fmt.Println(string(content))

		var sliContent SLIContent

		err = yaml.Unmarshal([]byte(content), &sliContent)
		if err != nil {
			fmt.Println("Error while unmarshalling YAML:", err)
			return
		}

		converter := pkg.NewSLIConverter()
		templates := converter.Convert(sliContent.Indicators, viper.GetString(providerName))
		for _, t := range templates {
			fmt.Printf("---------------------------------------------\n")
			// fmt.Printf("%v\n", t)
			yamlBytes, err := yaml.Marshal(t)
			if err != nil {
				fmt.Println("Error while marshalling to YAML:", err)
				return
			}

			yamlString := string(yamlBytes)
			fmt.Println(yamlString)
		}

		// TODO put it into the file available in --output or write to stdout
	},
}
