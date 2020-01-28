# Template2Helm

Template2Helm is a go utility that converts OpenShift templates into Helm charts.

## Usage

Install from source

```
make install
```

Run like so:

```
template2helm convert --template ~/src/openshift-templates/jobs/slack-notify-job-template.yml --chart ~/tmp/charts
```
