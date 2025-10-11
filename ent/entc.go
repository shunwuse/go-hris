//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	if err := entc.Generate("./schema", &gen.Config{
		Target:  "./entgen",
		Package: "github.com/shunwuse/go-hris/ent/entgen",
	}); err != nil {
		log.Fatal("running ent codegen:", err)
	}
}
