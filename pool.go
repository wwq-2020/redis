package redis

import (
	"context"
	"github.com/wwq1988/datastructure/lockfree/queue"
	"net"
)

// Dialer 链接简历方法
type Dialer func() (net.Conn, error)

// Pool 连接池
type Pool interface {
	Get(ctx context.Context) (net.Conn, error)
	Put(net.Conn)
}

// PoolOptions 连接池参数
type PoolOptions struct {
	Addr string
}

type pool struct {
	options *PoolOptions
	queue   *queue.Queue
}

// NewPool 初始化连接池
func NewPool(options *PoolOptions) Pool {
	return &pool{
		queue: queue.New(),
	}
}

func (p *pool) Get(ctx context.Context) (net.Conn, error) {
	item, err := p.queue.Pop()
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item.(net.Conn), nil
	}
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp4", p.options.Addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (p *pool) Put(conn net.Conn) {
	p.queue.Push(conn)
}
