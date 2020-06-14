# gomodcheck

[![ci](https://github.com/johejo/gomodcheck/workflows/ci/badge.svg?branch=master)](https://github.com/johejo/gomodcheck/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/johejo/gomodcheck/branch/master/graph/badge.svg)](https://codecov.io/gh/johejo/gomodcheck)
[![Go Report Card](https://goreportcard.com/badge/github.com/johejo/gomodcheck)](https://goreportcard.com/report/github.com/johejo/gomodcheck)

## Description

gomodcheck is a tool that reads go.mod and checks if a module needs updating.

## Install

```
go get github.com/johejo/gomodcheck
```

## Example

```
$ cat go.mod
module github.com/johejo/gomodcheck

go 1.14

require (
        github.com/PuerkitoBio/goquery v1.5.0
        github.com/google/go-github/v32 v32.0.0
        github.com/hashicorp/go-version v1.2.0
        golang.org/x/mod v0.3.0
        golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
```

```
$ gomodcheck
2020/06/14 22:24:28 github.com/PuerkitoBio/goquery is behind, latest=v1.5.1
2020/06/14 22:24:28 github.com/hashicorp/go-version v1.2.0 is latest
2020/06/14 22:24:28 github.com/google/go-github/v32 v32.0.0 is latest
2020/06/14 22:24:28 golang.org/x/mod v0.3.0 is latest
2020/06/14 22:24:28 golang/oauth2 has no tags
```

## Warning

gomodcheck uses GitHub API, so be careful of rate limit.<br>
If you set the environment variable `GITHUB_TOKEN`, gomodcheck will use it.
