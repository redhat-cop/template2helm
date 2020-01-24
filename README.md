# Template2Helm

Template2Helm is a go utility that converts OpenShift templates into Helm charts.

## Usage

Install deps

```
go mod vendor
```

Run like so:

```
go run main.go convert -template ~/src/openshift-templates/jobs/slack-notify-job-template.yml -chart ~/tmp/charts
```
