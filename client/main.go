package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"

	"github.com/someview/transport-benchmark/testdata"
)

var sendCount = int64(0)
var recvCount = int64(0)
var maxClientNum = 1e4

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
		StreamReorderingTimeout: time.Second * 20,
		RoundTripper: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: true,
			},
			QuicConfig: &quic.Config{
				HandshakeIdleTimeout: time.Minute,
				MaxIdleTimeout:       time.Hour,
				MaxIncomingStreams:   1 << 20,
				Allow0RTT:            true,
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancel()
	_, conn, err := d.Dial(ctx, "https://localhost:4242/", nil)
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
		maxData := make([]byte, 10)
		for {
			_, err := stream.Read(maxData)
			if err != nil {
				_ = stream.Close()
				fmt.Println("recv err:", err)
				return
			}
			atomic.AddInt64(&recvCount, 1)
		}
	}()

	go func() {
		for {
			// size is the same as application protocol on tcp
			write, err := stream.Write([]byte("hello 123456789"))
			if err != nil {
				fmt.Println("send err:", write)
				_ = stream.Close()
				return
			}
			atomic.AddInt64(&sendCount, 1)
		}
	}()
}

func ReportView() {
	for range time.NewTicker(time.Second * 20).C {
		send := atomic.LoadInt64(&sendCount)
		recv := atomic.LoadInt64(&recvCount)
		fmt.Println("时间:", time.Now(), "send:", send, "recv:", recv, "rate:", (send+recv)/20)
		atomic.StoreInt64(&sendCount, 0)
		atomic.StoreInt64(&recvCount, 0)
	}
}

func main() {
	var mode = flag.Int("mode", 0, "运行模式")
	flag.Parse()
	switch *mode {
	case slientMode:
		for i := 0; i < int(maxClientNum); i++ {
			go RunClient(slientMode)
		}
	case singleMode:
		for i := 0; i < int(maxClientNum)-1; i++ {
			go RunClient(slientMode)
		}
		go RunClient(multiMode)
	case multiMode:
		for i := 0; i < int(maxClientNum); i++ {
			go RunClient(multiMode)
		}
	}
	go ReportView()
	time.Sleep(time.Minute * 30)
}
