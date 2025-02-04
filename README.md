# kragin :rocket:

> **Version**: `v0.3.1`  
> **Description**: The package provides a set of tools for developing [Krakend](https://www.krakend.io/) plugins.

Welcome to **kragin**! This repository contains various helper packages that simplify plugin creation for the [Krakend API Gateway](https://www.krakend.io/community/).

## Features :sparkles:
- **Log**: Minimalistic logging interface + a no-op implementation.
- **Populator**: Struct field populator with tag-based configuration loading.
- **Plugin**: Easy-to-use “modifier” helpers for requests/responses.
- **Registerer**: Tools to register custom logic or transformations in Krakend.
- **Wrapper**: Utilities for request manipulation (body reading, form parsing, headers, etc.).

## Installation :wrench:

To use kragin in your own Go module, run:

```bash
go get github.com/dumb-tech/kragin@v0.3.1
```

Then import the packages you need:


```go
package yourplugin

import (
    "github.com/dumb-tech/kragin/log"
    "github.com/dumb-tech/kragin/populator"
    "github.com/dumb-tech/kragin/plugin"
    "github.com/dumb-tech/kragin/registerer"
    "github.com/dumb-tech/kragin/wrapper"
)
```

## Packages Overview :bookmark_tabs:

### 1. `log`

Defines a Logger interface with typical logging methods:

```go
type Logger interface {
    Debug(v ...any)
    Info(v ...any)
    Warning(v ...any)
    Error(v ...any)
    Critical(v ...any)
    Fatal(v ...any)
}
```

Includes a NoopLogger that does nothing (useful for testing or quiet modes). 
If you'll implement `RegisterLogger`method then no-op logger is useful for wait before this method will be called by
KrakenD.

### 2. `populator`

Allows you to populate struct fields from a map[string]any using struct tags:

```go
// Example struct
type Config struct {
    Host string `krakend:"host,required=true"`
    Port int    `krakend:"port,default=8080"`
}

// Usage
cfg := &Config{}
err := populator.PopulatePluginConfig(myConfigStore, "myconfig", cfg)
// Now cfg.Host, cfg.Port, etc. are populated!
```
 + Supports default values (default=...)
 + Supports required fields (required=true)
 + Supports nested structs


### 3. `plugin`

Helps create plugin modifiers for request/response transformations:

```go
p := plugin.NewModifier("my-awesome-plugin", "v1.0.0")

func (app *Application) requestHandlerFunc(config configuration.Configuration, deps *dependencies) func(any) (any, error) {
    return func(request any) (any, error) {
        r, ok := request.(wrapper.Request)
        if !ok {
            return request, wrapper.ErrUnknownRequestDataType(request)
        }
        req := wrapper.Modifier(r)
        
        req.SetHeader("X-Custom-Secret-Key", "my-secret-key")
        
        return req, nil
    }
}
```

RequestHandler and ResponseHandler let you add request/response modifiers.
Under the hood, uses registerer.ModifierRegisterer.

### 4. `registerer`

Core interfaces and helper types for plugin registration in Krakend:

 + `Registerer` interface for modifiers, handlers, clients.
 + `ModifierRegisterer` for adding request/response modifiers or custom ones.

### 5. `wrapper`

Utilities for working with incoming requests:

 + `Request` interface to get context, headers, params, body, etc.
 + `RequestWrapper` for reading/changing the body, manipulating query params, and more.

## Contributing :handshake:
We :heart: contributions! Feel free to open issues or submit pull requests.

## License :page_facing_up:

This project is licensed under the 
[GNU General Public License v3.0 (GPLv3)](https://www.gnu.org/licenses/gpl-3.0.txt).

#### This Markdown powered by ChatGPT o1
