package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/hsts"
)

func main() {
	fmt.Println("Hello, World!")
	http.Handle("/whats_time", hsts.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().Format(time.RFC3339)))
	})))
	http.ListenAndServeTLS(":80", "localhost.crt", "localhost.key", nil)
}

// openssl req -new -subj "/C=RU/ST=Msk/CN=localhost" -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr
// openssl x509 -req -sha256 -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt
//
// cd https
// go run main.go
// sudo tcpdump -A -i lo0 'port 80'
// curl -k https://localhost:80/whats_time -v
