package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	template "github.com/openshift/api/template/v1"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
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

			// Convert to json first
			jsonB, err := yaml.YAMLToJSON(yamlFile)
			checkErr(err, fmt.Sprintf("Error trasnfoming yaml to json: \n%s", string(yamlFile)))

			err = json.Unmarshal(jsonB, &myTemplate)
			checkErr(err, "Unable to marshal template")

			//fmt.Println("Printing Template")
			//fmt.Println(myTemplate)

			// TODO: Convert myTemplate.Objects into individual files
			var templates []*chart.File
			err = objectToTemplate(myTemplate.Objects, &templates)
			checkErr(err, "Failed object to template conversion")

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

// Convert the object list in the openshift template to a set of template files in the chart
func objectToTemplate(objects []runtime.RawExtension, templates *[]*chart.File) error {
	//fmt.Println(objects)

	for _, v := range objects {
		var k8sR unstructured.Unstructured
		err := json.Unmarshal(v.Raw, &k8sR)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to unmarshal Raw resource\n%v\n", v.Raw) + err.Error())
		}
		log.Printf("Creating a template for object %s", k8sR.GetName())
		tf := chart.File{
			Name: k8sR.GetKind() + k8sR.GetName(),
			Data: v.Raw,
		}
		*templates = append(*templates, &tf)
	}

	return nil
}
