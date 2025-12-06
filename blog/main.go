package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quartzeast/rock"
)

func main() {
	engine := rock.New()

	userGroup := engine.Group("/api/user")
	userGroup.AddRoute("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s\n", "Rockman")
	})
	userGroup.AddRoute("/profile", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is %s's profile\n", "Rockman")
	})

	orderGroup := engine.Group("/api/order")
	orderGroup.AddRoute("/list", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Order list for %s\n", "Rockman")
	})

	err := engine.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
