<div align="center">
  <h1>
    summer
    <br />
    <a href="https://github.com/utilyre/bevy_prank/blob/main/LICENSE">
      <img alt="license" src="https://img.shields.io/github/v/tag/utilyre/summer?label=version" />
    </a>
    <a href="https://go.dev/doc/go1.23">
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
    A high-performance utility for recursively generating checksums, designed
    to efficiently process large directories and files, offers a comprehensive
    API for seamless integration into other applications, ensuring reliable
    data integrity across various platforms.
  </p>
</div>

## Installation

```bash
go install github.com/utilyre/summer/cmd/summer@latest
```

## Usage

```
Usage of summer:
  -algo value
        sum using cryptographic hash function VALUE (default md5)
  -digest-workers int
        run N digest workers in parallel (default 1)
  -read-workers int
        run N read workers in parallel (default 1)
```
