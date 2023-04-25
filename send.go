package uf

import (
	"context"
	"errors"
	"io"
	"net"
	"net/netip"
)

const MTU = 2048

func Send(ctx context.Context, addr string, b io.Reader) (c int64, err error) {
	conn, err := net.DialUDP("udp", nil, net.UDPAddrFromAddrPort(netip.MustParseAddrPort(addr)))
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	for {
		offset, err := io.CopyN(conn, b, MTU)
		if err != nil || offset == 0 {
			if errors.Is(err, io.EOF) {
				conn.Write([]byte(io.EOF.Error()))
				return c, err
			}
			return offset, err
		}
		c += offset
		if err != nil {
			return c, err
		}
	}
}
