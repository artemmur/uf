package uf

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/netip"

	"golang.org/x/exp/slices"
)

func Recv(ctx context.Context, addr string, handle func(ctx context.Context, msg []byte)) error {
	conn, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(netip.MustParseAddrPort(addr)))
	if err != nil {
		return err
	}
	defer conn.Close()

	var data []byte
LOOP:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			buf := make([]byte, MTU)
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println(err)
			}

			if n > 0 {
				if slices.Compare(buf[:len([]byte(io.EOF.Error()))], []byte(io.EOF.Error())) == 0 {
					go handle(ctx, data)
					data = data[:0]
					continue LOOP
				}

				data = append(data, buf...)
			}
		}
	}
}
