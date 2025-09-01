# //+codemark

**Codemark** is a **comment-based annotation library** that lets you define
**custom annotations** (called _markers_) and map them to reusable definitions.
With Codemark, you can automate the generation of boilerplate files and
specifications tailored to your project needs.

A primary use case is the automatic generation of **JSON Schema** and **OpenAPI
specifications** from annotated code. By embedding lightweight markers directly
in your comments, you can keep specifications synchronized with your source
while reducing repetitive manual work.

### Key Features

- Define your own **markers** to represent custom annotations.
- Map markers to **custom definitions** that Codemark can interpret.
- Automatically generate **boilerplate files**, such as JSON Schema or OpenAPI
  specs.
- Keep your documentation and specifications **in sync** with your codebase.

## Install

```go
go install github.com/naivary/codemark@latest
```

## Simple Usage

Codemark integrates seamlessly with Go’s go generate directive, allowing you to
embed markers directly in your code comments. These markers are interpreted by
Codemark’s generators to produce artifacts such as OpenAPI-compatible JSON
Schemas.

```go
//go:generate codemark gen ./...
package main

type AuthRequest struct {
    // +openapi:schema:required
    // +openapi:schema:format="email"
    Email string

    // +openapi:schema:required
    Password string
}
```

### Generating the artifacts

To process the markers and generate the artifacts, run:

```go
go generate ./...
```

Codemark will scan your Go source files, interpret the markers, and produce the
appropriate boilerplate files (such as a JSON Schema in this case).
