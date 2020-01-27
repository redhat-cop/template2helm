module github.com/redhat-cop/template2helm

go 1.13

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

require (
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/labstack/gommon v0.3.0
	github.com/neelance/parallel v0.0.0-20160708114440-4de9ce63d14c // indirect
	github.com/openshift/api v0.0.0-20200123162640-f263157f58d3
	github.com/opentracing/basictracer-go v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/prometheus/client_golang v1.3.0 // indirect
	github.com/prometheus/common v0.9.1 // indirect
	github.com/redhat-cop/dash v0.0.0-20191106154059-10d43304fdc2 // indirect
	github.com/slimsag/godocmd v0.0.0-20161025000126-a1005ad29fe3 // indirect
	github.com/sourcegraph/ctxvfs v0.0.0-20180418081416-2b65f1b1ea81 // indirect
	github.com/sourcegraph/go-langserver v2.0.0+incompatible // indirect
	github.com/sourcegraph/jsonrpc2 v0.0.0-20191222043438-96c4efab7ee2 // indirect
	github.com/spf13/cobra v0.0.5
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200122134326-e047566fdf82 // indirect
	golang.org/x/tools v0.0.0-20200124200720-1b668f209185 // indirect
	gopkg.in/yaml.v2 v2.2.8
	helm.sh/helm/v3 v3.0.2
	k8s.io/apimachinery v0.17.1
	sigs.k8s.io/yaml v1.1.0
)
