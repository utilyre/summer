#!/usr/bin/env bash

hyperfine -w 3 --export-markdown=results.md --export-json=results.json \
	'./summer1 database' \
	'./summer2 generate -r database' \
	'./summer2 generate -r --read-jobs=8 --digest-jobs=8 database'
