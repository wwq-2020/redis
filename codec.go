package redis

import (
	"context"
	"io"
	"net"
)

// Codec 编解码器
type Codec interface {
	Encode(cmd string, args [][]byte) string
	EncodeTo(ctx context.Context, rw net.Conn, cmd string, args [][]byte) error
	Decode(ctx context.Context, r io.Reader) ([]byte, error)
}

type codec struct{}

func (c *codec) Encode(cmd string, args [][]byte) string {
	return ""
}

func (c *codec) EncodeTo(ctx context.Context, rw net.Conn, cmd string, args [][]byte) error {
	return nil
}

func (c *codec) Decode(ctx context.Context, r io.Reader) ([]byte, error) {
	return nil, nil
}
