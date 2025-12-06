package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// HandleFunc 将路由模式注册到 DefaultServeMux
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// ResponseWriter 实现了 io.Writer 接口，Fprintf 可直接向其写入格式化数据
		// 响应会在 handler 返回时自动 flush，无需手动关闭
		fmt.Fprintf(w, "Hello, %s!", "Rockman")
	})

	// ListenAndServe 阻塞当前 goroutine，第二个参数 nil 表示使用 DefaultServeMux
	// 生产环境建议使用 http.Server 结构体以配置超时等参数
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		// ListenAndServe 总是返回非 nil error，正常关闭返回 ErrServerClosed
		log.Fatal(err)
	}
}
