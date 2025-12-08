package main

import (
	"fmt"
	"log"

	"github.com/quartzeast/rock"
)

func main() {
	engine := rock.New()

	// Test static routes
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

	// Test named parameters
	userGroup.GET("/:id", func(ctx *rock.Context) {
		id := ctx.Param("id")
		fmt.Fprintf(ctx.Writer, "User ID: %s\n", id)
	})

	userGroup.GET("/:id/posts", func(ctx *rock.Context) {
		id := ctx.Param("id")
		fmt.Fprintf(ctx.Writer, "Posts for user %s\n", id)
	})

	userGroup.GET("/:id/posts/:postId", func(ctx *rock.Context) {
		id := ctx.Param("id")
		postId := ctx.Param("postId")
		fmt.Fprintf(ctx.Writer, "User %s, Post %s\n", id, postId)
	})

	// Test catch-all wildcard
	staticGroup := engine.Group("/static")
	staticGroup.GET("/*filepath", func(ctx *rock.Context) {
		filepath := ctx.Param("filepath")
		fmt.Fprintf(ctx.Writer, "Requesting static file: %s\n", filepath)
	})

	// Test root routes
	rootGroup := engine.Group("")
	rootGroup.GET("/", func(ctx *rock.Context) {
		fmt.Fprintf(ctx.Writer, "Welcome to Rock Web Framework!\n")
	})

	rootGroup.GET("/about", func(ctx *rock.Context) {
		fmt.Fprintf(ctx.Writer, "About Rock Web Framework\n")
	})

	fmt.Println("Server starting on :8081")
	fmt.Println("\nTest URLs:")
	fmt.Println("  GET  http://localhost:8081/")
	fmt.Println("  GET  http://localhost:8081/about")
	fmt.Println("  GET  http://localhost:8081/api/user/hello")
	fmt.Println("  POST http://localhost:8081/api/user/hello")
	fmt.Println("  GET  http://localhost:8081/api/user/123")
	fmt.Println("  GET  http://localhost:8081/api/user/123/posts")
	fmt.Println("  GET  http://localhost:8081/api/user/123/posts/456")
	fmt.Println("  GET  http://localhost:8081/static/css/style.css")
	fmt.Println("  ANY  http://localhost:8081/api/user/status")

	err := engine.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
