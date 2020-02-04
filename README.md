# Template2Helm

Template2Helm is a go utility that converts OpenShift templates into Helm charts.

## Installation

Installing is very simple. Simply download the proper binary from our latest [release](https://github.com/redhat-cop/template2helm/releases), and put it on your `$PATH`.

## Usage

template2helm has one primary function, `convert`. It can be used like so to convert an OpenShift template to a Helm chart.

```
template2helm convert --template ./examples/slack-notify-job-template.yml --chart ~/tmp/charts
```

We have several [example templates](./examples/) you can use to get started.

## Contribution

Please open issues and pull requests! Check out our [development guide](./docs/dev_guide.md) for more info on how to get started. We also follow the general [contribution guidelines](https://redhat-cop.github.io/contrib/) for pull requests outlined on the [Red Hat Community of Practice](https://redhat-cop.github.io) website.
