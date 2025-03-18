package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	// 1. 创建 HTTP/3 客户端
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证（仅开发用）
			NextProtos:         []string{"h3"},
		},
	}
	defer roundTripper.Close()

	client := &http.Client{
		Transport: roundTripper,
	}

	// 2. 发送请求
	resp, err := client.Get("https://localhost:443/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 3. 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n", body)
}
