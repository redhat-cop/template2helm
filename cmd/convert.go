package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

			// Convert myTemplate.Objects into individual files
			var templates []*chart.File
			err = objectToTemplate(&myTemplate.Objects, &myTemplate.ObjectLabels, &templates)
			checkErr(err, "Failed object to template conversion")

			// Convert myTemplate.Parameters into a yaml string map
			values := make(map[string]interface{})
			err = paramsToValues(&myTemplate.Parameters, &values, &templates)
			checkErr(err, "Failed parameter to value conversion")

			myChart := chart.Chart{
				Metadata: &chart.Metadata{
					Name:        myTemplate.ObjectMeta.Name,
					Version:     "v0.0.1",
					Description: myTemplate.ObjectMeta.Annotations["description"],
					Tags:        myTemplate.ObjectMeta.Annotations["tags"],
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
func objectToTemplate(objects *[]runtime.RawExtension, templateLabels *map[string]string, templates *[]*chart.File) error {
	o := *objects

	for _, v := range o {
		var k8sR unstructured.Unstructured
		err := json.Unmarshal(v.Raw, &k8sR)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to unmarshal Raw resource\n%v\n", v.Raw) + err.Error())
		}
		name := "templates/" + strings.ToLower(k8sR.GetKind()+".yaml")

		labels := k8sR.GetLabels()
		if labels == nil {
			k8sR.SetLabels(*templateLabels)
		} else {
			for key, value := range *templateLabels {
				labels[key] = value
			}
			k8sR.SetLabels(labels)
		}

		updatedJSON, err := k8sR.MarshalJSON()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to marshal Unstructured record to JSON\n%v\n", k8sR) + err.Error())
		}

		log.Printf("Creating a template for object %s", name)
		data, err := yaml.JSONToYAML(updatedJSON)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to marshal Raw resource back to YAML\n%v\n", updatedJSON) + err.Error())
		}

		tf := chart.File{
			Name: name,
			Data: data,
		}
		*templates = append(*templates, &tf)
	}

	return nil
}

func paramsToValues(param *[]template.Parameter, values *map[string]interface{}, templates *[]*chart.File) error {

	p := *param
	t := *templates
	v := *values

	for _, pm := range p {
		name := strings.ToLower(pm.Name)
		log.Printf("Convert parameter %s to value .%s", pm.Name, name)

		for i, tf := range t {
			// Search and replace ${PARAM} with {{ .Values.param }}
			raw := tf.Data
			// Handle string format parameters
			ns := strings.ReplaceAll(string(raw), fmt.Sprintf("${%s}", pm.Name), fmt.Sprintf("{{ .Values.%s }}", name))
			// TODO Handle binary formatted data differently
			ns = strings.ReplaceAll(ns, fmt.Sprintf("${{%s}}", pm.Name), fmt.Sprintf("{{ .Values.%s }}", name))
			ntf := chart.File{
				Name: tf.Name,
				Data: []byte(ns),
			}

			t[i] = &ntf
		}

		if pm.Value != "" {
			v[name] = pm.Value
		} else {
			v[name] = "# TODO: must define a default value for ." + name
		}
	}

	*templates = t
	*values = v

	return nil
}
