package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"

	"github.com/someview/transport-benchmark/testdata"
)

type MyHandler struct {
	activeCount int64
}

var wsServer *webtransport.Server

func init() {
	wsServer = &webtransport.Server{
		H3: http3.Server{
			Handler: &MyHandler{},
			Addr:    "localhost:4242",
			// Port:       4242,
			QuicConfig: &quic.Config{Allow0RTT: true},
			// TLSConfig:  &tls.Config{InsecureSkipVerify: true}
		},
	}
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsServer.Upgrade(w, r)
	if err != nil {
		log.Printf("upgrading failed: %s", err)
		w.WriteHeader(500)
		return
	}
	atomic.AddInt64(&h.activeCount, 1)
	fmt.Println("server recv webtransport conn:", h.activeCount) // 最终的活跃连接数
	// Handle the connection. Here goes the application logic.
	go func() {
		for {
			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				fmt.Println("server read err:", err)
				return
			}
			if _, err := io.Copy(stream, stream); err != nil {
				fmt.Println("server write err:", err)
			}
		}
	}()
}

func RunServer() {
	log.Println("hello world:", wsServer.ListenAndServeTLS(testdata.GetCertificatePaths()))
}

func main() {
	// 设置最大核心数量为2
	runtime.GOMAXPROCS(2)
	go RunServer()
	time.Sleep(time.Minute * 30)
}
