module github.com/redhat-cop/template2helm

go 1.13

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

require (
	github.com/openshift/api v0.0.0-20200123162640-f263157f58d3
	github.com/redhat-cop/dash v0.0.0-20191106154059-10d43304fdc2 // indirect
	github.com/spf13/cobra v0.0.5
	gopkg.in/yaml.v2 v2.2.8
	helm.sh/helm/v3 v3.0.2
)
