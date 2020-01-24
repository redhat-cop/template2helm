package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	template "github.com/openshift/api/template/v1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

const (
	tplPathDefault   = "."
	tplPathUsage     = "Path to an OpenShift Template, relative or absolute"
	chartPathDefault = "."
	chartPathUsage   = "Destination directory of the Chart."
)

var (
	tplPath   string
	chartPath string

	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Given the path to an OpenShift template file, spit out a Helm chart.",
		Long:  `Long version...`,
		RunE: func(cmd *cobra.Command, args []string) error {

			var myTemplate template.Template

			yamlFile, err := ioutil.ReadFile(filepath.Clean(tplPath))
			if err != nil {
				return fmt.Errorf("Couldn't load template: %v", err)
			}

			err = yaml.Unmarshal(yamlFile, &myTemplate)
			checkErr(err, "Unable to marshal template")

			// TODO: Convert myTemplate.Objects into individual files
			var templates []*chart.File
			// TODO: Convert myTemplate.Parameters into a yaml string map
			var values map[string]interface{}

			myChart := chart.Chart{
				Metadata: &chart.Metadata{
					Name: myTemplate.ObjectMeta.Name,
					// TODO: add description, labels, etc.
				},
				Templates: templates,
				Values:    values,
			}

			if myChart.Metadata.Name == "" {
				ext := filepath.Ext(tplPath)
				name := filepath.Base(string(tplPath))[0 : len(filepath.Base(string(tplPath)))-len(ext)]
				myChart.Metadata.Name = name
			}

			err = chartutil.SaveDir(&myChart, chartPath)
			checkErr(err, fmt.Sprintf("Failed to save chart %s", myChart.Metadata.Name))

			return nil
		},
	}
)

func init() {
	convertCmd.Flags().StringVarP(&tplPath, "template", "t", tplPathDefault, tplPathUsage)
	convertCmd.Flags().StringVarP(&chartPath, "chart", "c", chartPathDefault, chartPathUsage)
	rootCmd.AddCommand(convertCmd)
}

func checkErr(err error, msg string) {
	if err != nil {
		fmt.Print(fmt.Errorf(msg + err.Error()))
		os.Exit(1)
	}
	return
}
