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

This will give you a list of generators available and a short description of
their usage. If you want to know more about a generator and the available
resources and options you can use the following command:

```bash
# To get a list of all resources
codemark explain openapi

# To get a list of all options of a resource
codemark explain openapi:schema

# To get a detail description of the option and how to use it.
codemark explain openapi:schema:minItems
```

## Architecture

The internal architecture of codemark is composed of the following key
components:

1. Loader: Read a Go project and extract key information.
2. Converter: Convert the value type of a marker to the user-defined Option type
   allowing users to associate methods with options (extensible)
3. Generator: Generate artifacts based on the information by the loader
   (extensible)
4. Outputer: Write the generated artifacts to a target (extensible)
