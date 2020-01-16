package redis

import (
	"context"
	"io"
	"net"
)

// Dialer 链接简历方法
type Dialer func() (net.Conn, error)

// Pool 连接池
type Pool interface {
	Get(ctx context.Context) (io.ReadWriter, error)
	Put(io.ReadWriter)
}

// PoolOptions 连接池参数
type PoolOptions struct {
	Addr string
}

type pool struct {
	options *PoolOptions
}

// NewPool 初始化连接池
func NewPool(options *PoolOptions) Pool {
	return &pool{}
}

func (p *pool) Get(ctx context.Context) (io.ReadWriter, error) {
	return nil, nil
}

func (p *pool) Put(io.ReadWriter) {
}
