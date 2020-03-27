#!/bin/bash

#### WORK-IN-PROGRESS

# git clone https://github.com/kubernetes/kubernetes.git

# copied staging, src/k8s.io/kube-scheduler
# /tmp/tmpdc3iihix/src/k8s.io/kubernetes/staging/src/k8s.io/kube-scheduler

# Change this to installation directory of gen-crd-api-reference-docs
GEN_DOCS=/tmp/gen-crd-api-reference-docs

# Change this to where the website repository is cloned.
K8S_WEBSITE=/website

K8S_SCHEDULER=k8s.io/kube-scheduler

# table style substitutions
TABLE_SUB='<div class=\"table-responsive\"><table class=\"table table-bordered\">'
THEAD_SUB='<thead class=\"thead-light\">'

# FIX, generate in reference docs
CONTENT_DIR=kubernetes-sigs/reference-docs

echo "+++
title = \"Kube Scheduler API Reference\"
description = \"Reference documentation for Kube Scheduler\"
weight = 100
+++" > ${CONTENT_DIR}/common.md

${GEN_DOCS}/gen-crd-api-reference-docs -api-dir ./config -config ${GEN_DOCS}/example-config.json -template-dir ${GEN_DOCS}/template -out-file temp_common.md
cat temp_common.md >> ${CONTENT_DIR}/common.md
rm temp_common.md

sed 's/<table>/'"$TABLE_SUB"'/g' ${CONTENT_DIR}/common.md > temp_common.md
sed 's/<thead>/'"$THEAD_SUB"'/g' temp_common.md > ${CONTENT_DIR}/common.md
rm temp_common.md