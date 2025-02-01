#!/usr/bin/env bash

hyperfine -w 3 --export-markdown=results.md --export-json=results.json \
	'./summer-v0.1.0 database' \
	'./summer-v0.8.0 generate -r database' \
	'./summer-v0.8.0 generate -r --open-file-jobs=8 --digest-jobs=8 database'
