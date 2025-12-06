package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quartzeast/rock"
)

func main() {
	engine := rock.New()

	engine.AddRoute("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s\n", "Rockman")
	})

	err := engine.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
