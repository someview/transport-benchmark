package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"

	"github.com/someview/transport-benchmark/testdata"
)

var sendCount = 0
var recvCount = 0

var multiMode = 0  // 大量客户端，均发送消息
var singleMode = 1 // 大量客户端，只有一个客户端在发送消息
var slientMode = 2 // 大量客户端，不发送消息

func RunClient(runMode int) {

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	testdata.AddRootCA(pool)

	var d = &webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: true,
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	_, conn, err := d.Dial(ctx, "https://localhost:4242/", nil)
	// overUDP
	if err != nil {
		log.Fatalln("dial err:", err)
	}
	if runMode != multiMode {
		return
	}
	stream, err := conn.OpenStream()
	if err != nil {
		log.Fatalln("open stream err:", err)
		return
	}

	go func() {
		maxData := make([]byte, 4096)
		for {
			_, err := stream.Read(maxData)
			if err != nil {
				_ = stream.Close()
				return
			}

		}
	}()

	go func() {
		for {
			write, err := stream.Write([]byte("hello"))
			if err != nil {
				fmt.Println("send err:", write)
				return
			}
		}
	}()
}

func ReportView() {
	for range time.NewTicker(time.Second * 20).C {
		fmt.Println("时间:", time.Now(), "发送消息数量:", sendCount, "recvCount:", recvCount)
	}
}

func main() {
	var mode = flag.Int("mode", 0, "运行模式")
	flag.Parse()
	switch *mode {
	case slientMode:
		for i := 0; i < 1e4; i++ {
			go RunClient(slientMode)
		}
	case singleMode:
		for i := 0; i < 1e4-1; i++ {
			go RunClient(slientMode)
		}
		go RunClient(multiMode)
	case multiMode:
		for i := 0; i < 1e4; i++ {
			go RunClient(multiMode)
		}
	}
	go ReportView()
}
