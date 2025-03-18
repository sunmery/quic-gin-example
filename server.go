package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	// 1. 创建 Gin 路由引擎
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello over HTTP/3!"})
	})

	// 2. 加载 TLS 证书
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	// 3. 创建 UDP 监听
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 443})
	if err != nil {
		log.Fatal(err)
	}

	// 4. 创建 QUIC Transport
	tr := quic.Transport{
		Conn:               udpConn,
		ConnectionIDLength: 8,
		StatelessResetKey:  nil,
	}

	// 5. 创建早期监听器 (支持 0-RTT)
	quicConf := &quic.Config{
		EnableDatagrams: true,
		Allow0RTT:       true, // 启用 0-RTT 支持
	}

	earlyListener, err := tr.ListenEarly(
		&tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"h3"},
		},
		quicConf,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer earlyListener.Close()

	// 6. 创建 HTTP/3 服务器
	server := http3.Server{
		Handler: router,
	}

	fmt.Println("HTTP/3 server listening on :443")
	if err := server.ServeListener(earlyListener); err != nil {
		log.Fatal(err)
	}
}
