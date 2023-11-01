package main

import (
	"bufio"
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

type bufferedWriteCloser struct {
	*bufio.Writer
	io.Closer
}

// NewBufferedWriteCloser creates an io.WriteCloser from a bufio.Writer and an io.Closer
func NewBufferedWriteCloser(writer *bufio.Writer, closer io.Closer) io.WriteCloser {
	return &bufferedWriteCloser{
		Writer: writer,
		Closer: closer,
	}
}

func (h bufferedWriteCloser) Close() error {
	if err := h.Writer.Flush(); err != nil {
		return err
	}
	return h.Closer.Close()
}

type MyHandler struct {
	activeCount int64
}

var wsServer *webtransport.Server

func init() {
	// filename := "server.qlog"
	// f, err := os.Create(filename)
	// if err != nil {
	// 	log.Fatalln("create log file err:", err)
	// }
	// wc := bufio.NewWriter(f)
	quicConf := &quic.Config{
		Allow0RTT:            true,
		HandshakeIdleTimeout: time.Second * 120,
		MaxIncomingStreams:   1 << 20, // 40万个incoming
		// Tracer: func(ctx context.Context, p logging.Perspective,
		// 	ci quic.ConnectionID) *logging.ConnectionTracer {
		// 	return qlog.NewConnectionTracer(NewBufferedWriteCloser(wc, f), p, ci)
		// },
	}

	wsServer = &webtransport.Server{
		H3: http3.Server{
			Handler: &MyHandler{},
			Addr:    "localhost:4242",
			// Port:       4242,
			QuicConfig: quicConf,
			// TLSConfig:  &tls.Config{InsecureSkipVerify: true}
		},
		StreamReorderingTimeout: time.Second * 60,
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
	fmt.Println("server recv webtransport conn:", h.activeCount, "routines:", runtime.NumGoroutine()) // 最终的活跃连接数
	// Handle the connection. Here goes the application logic.
	go func() {
		for {
			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				fmt.Println("server read err:", err)
				_ = conn.CloseWithError(webtransport.SessionErrorCode(405), "read error")
				return
			}
			if _, err := io.Copy(stream, stream); err != nil {
				fmt.Println("send err:", err)
				_ = conn.CloseWithError(webtransport.SessionErrorCode(406), "write error")
				return
			}
		}
	}()
}

func RunServer() {
	log.Println("run ws server:", wsServer.ListenAndServeTLS(testdata.GetCertificatePaths()))
}

func main() {
	runtime.GOMAXPROCS(2)
	go RunServer()
	time.Sleep(time.Minute * 30)
}
