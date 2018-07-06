#!/bin/bash

PKG=i18n_test go run ../binfs/cmd/binfs/main.go testdata > i18ndata_test.go
