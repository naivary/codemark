# //+codemark

**Codemark** is a **comment-based annotation library** that lets you define
**custom annotations** (called _markers_) and map them to reusable definitions.
With Codemark, you can automate the generation of boilerplate files and
specifications tailored to your project needs.

The generation of boilerplate files is done by `generators` which can be
extended by you.

By embedding lightweight markers directly in your comments, you can keep
specifications synchronized with your source while reducing repetitive manual
work.

### Key features

- Define your own **markers** to represent custom annotations.
- Map markers to **custom definitions** that generators can interpret.
- Automatically generate **boilerplate files**, such as JSON Schema or OpenAPI
  specs.
- Keep your documentation and specifications **in sync** with your codebase.
- Implement your own outputer, generators and converter to fulfill your needs.

## Install

```go
go install github.com/naivary/codemark@latest
```

## Getting started

To get started run the following command:

```bash
# NOTE: --kind=gen can be left out becuase its the default. It's only included
# for better understanding
codemark explain --kind=gen all
```

This gives you a list of generators available and a summary. If you want to know
more about a generator, their available resources and options you can use the
following commands:

```bash
# Get a list of all resources
codemark explain openapi

# Get a list of all options of a resource
codemark explain openapi:schema

# Get a detail description of the option and how to use it.
codemark explain openapi:schema:minItems
```

## OpenAPI generator

One of the builtin generators is the OpenAPI generator. To generate a OpenAPI
compatible JSON Schema you can use the following example:

```go
//go:generate codemark gen ./... -o openapi:fs

type AuthRequest struct {
    // +openapi:schema:format="email"
    // +openapi:schema:required
    Email string

    // +openapi:schema:required
    Password string
}

func main() {}
```

This will generate JSON Schemas and write them to the local file system.
