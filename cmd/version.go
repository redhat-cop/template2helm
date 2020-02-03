package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of template2helm",
	Long:  `All software has versions. This is template2helm's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
