# tfmigrator

[![Build Status](https://github.com/suzuki-shunsuke/tfmigrator/workflows/test/badge.svg)](https://github.com/suzuki-shunsuke/tfmigrator/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/tfmigrator)](https://goreportcard.com/report/github.com/suzuki-shunsuke/tfmigrator)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/tfmigrator.svg)](https://github.com/suzuki-shunsuke/tfmigrator)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/tfmigrator/master/LICENSE)

CLI tool to migrate Terraform configuration and State

## Overview

`tfmigrator` is a CLI tool to migrate Terraform configuration and State.

## Requirement

* Terraform
* [hcledit](https://github.com/minamijoyo/hcledit)

## Install

Download a binary from the [replease page](https://github.com/suzuki-shunsuke/tfmigrator/releases).

## How to use

```
$ terraform show -json > state.json
$ vi tfmigrator.yaml
$ cat *.tf | tfmigrator run [-skip-state] state.json
```

## Configuration file

[CONFIGURATION.md](docs/CONFIGURATION.md)

example of tfmigrator.yaml

```yaml
items:
- rule: |
    "name" not in Values
  exclude: true
- rule: |
    Values.name contains "foo"
  state_out: foo/terraform.tfstate
  resource_name: "{{.Values.name}}"
  tf_path: foo/resource.tf
- rule: |
    Values.name contains "bar"
  state_out: bar/terraform.tfstate
  resource_name: "{{.Values.name}}"
  tf_path: bar/resource.tf
```

## LICENSE

[MIT](LICENSE)
