# Progotgen DUMP

The `protoc` compiler plugin which dumps the generation request details
in `google.golang.org/protobuf/compiler/protogen` format to `stderr`. No files written to disk during the dump.

Quick crafted for personal use to see `protogen` data internals while building `protoc` plugin on it's code base.

## Install

```shell
go install github.com/albenik/protoc-gen-dump/cmd/protoc-gen-dump@latest
```

## Usage

```shell
protoc -protogen-dump_out=. path/to/file.proto
```
