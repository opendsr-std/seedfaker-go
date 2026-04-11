# seedfaker-go

Go bindings for [seedfaker](https://github.com/opendsr-std/seedfaker) — deterministic synthetic data generator.

## Install

```bash
go get github.com/opendsr-std/seedfaker-go
```

Pre-built `libseedfaker_ffi` binaries are included for linux-x64, linux-arm64, darwin-x64, darwin-arm64.

## Usage

```go
package main

import (
    "fmt"
    seedfaker "github.com/opendsr-std/seedfaker-go"
)

func main() {
    f, _ := seedfaker.New(seedfaker.Options{Seed: "demo", Locale: "en"})
    defer f.Close()

    name, _ := f.Generate("name")
    email, _ := f.Generate("email")
    fmt.Println(name, email)

    records, _ := f.GenerateRecords(seedfaker.RecordOptions{
        Fields: []string{"name", "email", "phone"},
        N:      5,
        Ctx:    "strict",
    })
    for _, r := range records {
        fmt.Printf("%s\t%s\t%s\n", r["name"], r["email"], r["phone"])
    }

    fp, _ := seedfaker.Fingerprint()
    fmt.Println("fingerprint:", fp)
}
```

## Determinism

Same seed + same arguments = identical output across CLI, Python, Node.js, Go, PHP, Ruby.

## Requirements

- Go >= 1.21
- CGO enabled (default on most platforms)

Full documentation: [github.com/opendsr-std/seedfaker](https://github.com/opendsr-std/seedfaker)
