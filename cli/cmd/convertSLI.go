package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const outputFileName = "output"

func init() {
	flags := convertSLICmd.Flags()

	flags.StringP(outputFileName, "o", "manifests.yaml", "Set the output file name")

	_ = viper.BindPFlag(outputFileName, flags.Lookup(outputFileName))
}

// startCmd represents the start command
var convertSLICmd = &cobra.Command{
	Use:   "convert-sli",
	Short: "Convert SLI to AnalysisTemplate",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Converting...")
	},
}
