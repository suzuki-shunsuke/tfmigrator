# tfmigrator-cli

[![Build Status](https://github.com/suzuki-shunsuke/tfmigrator-cli/workflows/test/badge.svg)](https://github.com/suzuki-shunsuke/tfmigrator-cli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/tfmigrator-cli)](https://goreportcard.com/report/github.com/suzuki-shunsuke/tfmigrator-cli)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/tfmigrator-cli.svg)](https://github.com/suzuki-shunsuke/tfmigrator-cli)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/tfmigrator-cli/main/LICENSE)

CLI tool to migrate Terraform configuration and State

## Note: this CLI was renamed from tfmigrator to tfmigrator-cli

This repository was renamed from tfmigrator to tfmigrator-cli, because we started to develop tfmigrator as Go package (framework).

https://github.com/suzuki-shunsuke/tfmigrator-cli/issues/25

## Blog written in Japanese

[terraformer で雑に生成した tf ファイル と state を分割したくてツールを書いた](https://techblog.szksh.cloud/tfmigrator/)

## Overview

`tfmigrator` is a CLI tool to migrate Terraform configuration and State into multiple states.
`tfmigrator` configures rules to classify resources and migrate resources to other states via `terraform state mv` and [hcledit](https://github.com/minamijoyo/hcledit).

## Requirement

* Terraform CLI
* [hcledit](https://github.com/minamijoyo/hcledit)

## Install

Download a binary from the [replease page](https://github.com/suzuki-shunsuke/tfmigrator-cli/releases).

## How to use

```
$ vi tfmigrator.yaml
$ cat *.tf | tfmigrator run [-skip-state]
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

## Restriction

tfmigrator migrates Terraform configuration with hcledit and doesn't support to expand the expression.
For example, if Terraform configuration refers local values, the migrated configuration may be broken.

## LICENSE

[MIT](LICENSE)
