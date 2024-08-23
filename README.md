<div align="center">
  <h1>
    summer
    <br />
    <a href="httpshttps://github.com/utilyre/summer/releases/latest">
      <img alt="license" src="https://img.shields.io/github/v/tag/utilyre/summer?label=version" />
    </a>
    <a href="https://go.dev">
      <img alt="downloads" src="https://img.shields.io/github/go-mod/go-version/utilyre/summer?label=go" />
    </a>
    <a href="https://github.com/utilyre/summer/issues">
      <img alt="issues" src="https://img.shields.io/github/issues/utilyre/bevy_prank?label=issues" />
    </a>
    <a href="https://github.com/utilyre/summer/actions/workflows/ci.yaml">
      <img alt="version" src="https://img.shields.io/github/actions/workflow/status/utilyre/summer/ci.yaml?label=ci" />
    </a>
  </h1>
  <p>
    High-performance utility for generating checksums in parallel.
  </p>
</div>

## Installation

- [Latest Release](https://github.com/utilyre/summer/releases/latest)

- Manual

  ```bash
  go install github.com/utilyre/summer/cmd/summer@latest
  ```

## Usage

For starters, if you just run the generate command (without any arguments), it
should give you the checksum of every file in your current directory.

```
$ summer generate
081ecc5e6dd6ba0d150fc4bc0e62ec50  bar
764efa883dda1e11db47671c4a3bbd9e  foo
168065a0236e2e64c9c6cdd086c55f63  nested/baz
```

In case you only care about certain files/directories, list those as arguments.

```
$ summer generate bar nested
081ecc5e6dd6ba0d150fc4bc0e62ec50  bar
168065a0236e2e64c9c6cdd086c55f63  nested/baz
```

To utilize more cores of your CPU, pass `--read-workers=n` and
`--digest-workers=m` flags, where `n` and `m` are the number of workers used
for each task respectively.

Run `summer help generate` to learn more about different flags.

## API

It is possible to call the API of this utility directly in your own
application. Here's an example:

```go
package main

import (
	"context"
	"log"

	"github.com/utilyre/summer/pkg/summer"
)

func main() {
	results, err := summer.SumTree(context.TODO(), []string{"."})
	if err != nil {
		log.Fatal(err)
	}

	for result := range results {
		if result.Err != nil {
			log.Println(result.Err)
			continue
		}

		// TODO
	}
}
```
