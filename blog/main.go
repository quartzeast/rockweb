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
	userGroup.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s\n", "Rockman")
	})
	userGroup.POST("/profile", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is %s's profile\n", "Rockman")
	})
	userGroup.ANY("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User status endpoint\n")
	})

	err := engine.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
