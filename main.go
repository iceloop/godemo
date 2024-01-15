package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

// Message 定义一个简单的结构体来处理JSON数据
type Message struct {
	Text string `json:"text"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "欢迎来到Go Web服务!")
	case "POST":
		var msg Message
		// 解析JSON数据
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 返回JSON响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	default:
		http.Error(w, "Method is not supported.", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", helloHandler)

	fmt.Printf("启动服务器在 http://localhost:8080/\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
