# How to generate kube-scheduler API doc?

## 1. Clone k8s.io/kube-scheduler repo

```sh
$ git clone git@github.com:kubernetes/kube-scheduler.git
```

## 2. Install `gen-crd-api-reference-docs`

Download [gen-crd-api-reference-docs](https://github.com/ahmetb/gen-crd-api-reference-docs), and then execute `go build .` to generate `gen-crd-api-reference-docs` command.

## 3. Generate API doc

```sh
$ export GEN_DOCS="step 2 download path"
```

Change to `kube-scheduler` directory and execute the follows command:
```sh
$ echo "---
title = \"Kube Scheduler API Reference\"
description = \"Reference documentation for Kube Scheduler\"
weight = 100
---" >> kube-scheduler-api.md

$ ${GEN_DOCS}/gen-crd-api-reference-docs -api-dir ./config -config ${GEN_DOCS}/example-config.json -template-dir ${GEN_DOCS}/template -out-file kube-scheduler-api.md
```

## 4. Copy API doc to website repo

Copy step 3 `kube-scheduler-api.md` to [website](https://github.com/kubernetes/website/static/docs/reference/generated/)
