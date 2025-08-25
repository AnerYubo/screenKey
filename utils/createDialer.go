// utils/utils.go
package utils

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

var Dialer proxy.Dialer = proxy.Direct // 默认直连

func InitProxy(proxyAddr string) error {
	if proxyAddr == "" {
		Dialer = proxy.Direct
		return nil
	}

	// 使用 SOCKS5 代理，并通过 proxy.Direct 进行 TCP 连接（包括 DNS 解析走代理）
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return fmt.Errorf("无法连接到代理: %v", err)
	}

	Dialer = dialer
	return nil
}

func DialWithTimeout(network, addr string, timeout time.Duration) (net.Conn, error) {
	type result struct {
		conn net.Conn
		err  error
	}

	resultChan := make(chan result, 1)

	go func() {
		conn, err := Dialer.Dial(network, addr)
		resultChan <- result{conn, err}
	}()

	select {
	case res := <-resultChan:
		return res.conn, res.err
	case <-time.After(timeout):
		return nil, fmt.Errorf("连接超时：%s", addr)
	}
}
