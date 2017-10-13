# kbn

learn bounded tree-width probabilistic models with latent variables via k-tree sampling

[![Build Status](https://travis-ci.org/britojr/kbn.svg?branch=master)](https://travis-ci.org/britojr/kbn)
[![Coverage Status](https://coveralls.io/repos/github/britojr/kbn/badge.svg?branch=master)](https://coveralls.io/github/britojr/kbn?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/britojr/kbn)](https://goreportcard.com/report/github.com/britojr/kbn)
[![GoDoc](https://godoc.org/github.com/britojr/kbn?status.svg)](http://godoc.org/github.com/britojr/kbn)

___

## Installation and usage

### Get, install and test:

		go get -u github.com/britojr/kbn...
		go install github.com/britojr/kbn...
		go test github.com/britojr/kbn... -cover

### Usage:

		kbn --help
		Usage: kbn <command> [options]

		Commands:

			struct		sample bounded tree-width structure
			param		parameter learning using expectation-maximization
			marginals	compute marginal distribution for each variable in the model

		To see details of each command:

			kbn <command> --help

### Examples:

		# move to examples directory
		cd examples/

		# sample a structure with tree-width 4 and 10 latent variables
		kbn struct -d example.csv -cs example-struct.ct0 -k 4 -h 10

		# learn parameters
		kbn param -dist rand -mode indep -cl example-struct.ct0 -cs example-struct.ctp -d example.csv

		# compute marginals
		kbn marginals -c example-struct.ctp -m example.mar
