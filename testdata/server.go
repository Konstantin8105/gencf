package main

import (
	"fmt"
	"log"
	"net/http"
)

//go:generate gensf -struct=M -o=struct_gen.go -i=server.go

// M is some struct
type M struct {
	// parameter a
	a int

	// parameter b
	b uint8

	// parameter c
	// with multiline comments
	c float32

	// d is anonymous struct
	d struct {
		// internal value d.e
		e uint16

		// internal value d.f
		f float64
	}

	// h with slice
	h []string
}

func handler(w http.ResponseWriter, r *http.Request) {
	var m M
	m.h = append(m.h, "--some--")
	fmt.Fprintf(w, m.FormDefault("/resultOfM"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
