# s3discover
[![GitHub release](http://img.shields.io/github/release/jakewarren/s3discover.svg?style=flat-square)](https://github.com/jakewarren/s3discover/releases])
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/s3discover/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/s3discover)](https://goreportcard.com/report/github.com/jakewarren/s3discover)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)

> scrape a website to discover s3 buckets

## Install
### Option 1: Binary

Download the latest release from [https://github.com/jakewarren/s3discover/releases/latest](https://github.com/jakewarren/s3discover/releases/latest)

### Option 2: From source

```
go get github.com/jakewarren/s3discover
```
## Usage

```
❯ s3discover -h
Usage: s3discover [<flags>] <domain>

Example: s3discover github.com

Optional flags:

  -d, --debug     enable debug logging
  -h, --help      display help
  -v, --verbose   enable verbose output
  -V, --version   display version

```
## Example

```
❯ s3discover github.com
shopifyorderlimits.s3.amazonaws.com
github-cloud.s3.amazonaws.com
```

## Similar Projects

https://github.com/random-robbie/AWS-Scanner/

## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## License

MIT © 2018 Jake Warren
