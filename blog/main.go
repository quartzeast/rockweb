package main

import (
	"fmt"
	"log"

	"github.com/quartzeast/rock"
)

func main() {
	engine := rock.New()

	userGroup := engine.Group("/api/user")
	userGroup.GET("/hello", func(ctx *rock.Context) {
		fmt.Fprintf(ctx.Writer, "Hello, Rockman!\n")
	})

	userGroup.POST("/hello", func(ctx *rock.Context) {
		fmt.Fprintf(ctx.Writer, "Posted to hello endpoint\n")
	})

	userGroup.POST("/profile", func(ctx *rock.Context) {
		fmt.Fprintf(ctx.Writer, "This is %s's profile\n", "Rockman")
	})

	userGroup.ANY("/status", func(ctx *rock.Context) {
		fmt.Fprintf(ctx.Writer, "User status endpoint\n")
	})

	err := engine.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
