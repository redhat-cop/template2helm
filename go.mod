module github.com/redhat-cop/template2helm

go 1.13

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

require (
	github.com/openshift/api v0.0.0-20200123162640-f263157f58d3
	github.com/spf13/cobra v1.1.3
	helm.sh/helm/v3 v3.6.1
	k8s.io/apimachinery v0.21.0
	sigs.k8s.io/yaml v1.2.0
)
