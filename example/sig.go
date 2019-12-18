package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//优雅的退出  信号
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello,world!"))
	})
	server := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		fmt.Println(server.ListenAndServe())
	}()
	sigs := make(chan os.Signal)
	signal.Notify(sigs)
	c := <-sigs
	fmt.Println(c.String())
	fmt.Println(server.Close())
	time.Sleep(20 * time.Second)
	fmt.Println("退出")
}
