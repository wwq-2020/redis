package redis

import (
	"bytes"
	"context"
	"errors"
	"net"
	"time"
)

// Client 客户端
type Client interface {
	Set(key string, val interface{}, expiration time.Duration) (string, error)
	Get(key string, val interface{}) error
}

type client struct {
	options *ClientOptions
}

// ClientOptions 客户端参数
type ClientOptions struct {
	Pool       Pool
	Codec      Codec
	Password   string
	Serializer Serializer
}

func (o *ClientOptions) fill() *ClientOptions {
	if o == nil {
		return &ClientOptions{
			Pool:       NewPool(nil),
			Codec:      NewDefaultCodec(),
			Serializer: NewDefaultSerializer(),
		}
	}
	return o
}

// NewClient 初始化客户端
func NewClient(options *ClientOptions) Client {
	return &client{
		options: options.fill(),
	}
}

func (c *client) Set(key string, val interface{}, expiration time.Duration) (string, error) {
	result, err := c.SetContext(context.Background(), key, val, expiration)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *client) SetContext(ctx context.Context, key string, val interface{}, expiration time.Duration) (string, error) {
	result, err := c.roudtrip(ctx, "set", c.buildSetArgs(key, val, expiration))
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (c *client) buildSetArgs(key string, val interface{}, expiration time.Duration) []interface{} {
	args := make([]interface{}, 2, 4)
	args[0] = key
	args[1] = val
	if expiration <= 0 {
		return args
	}
	if usePrecise(expiration) {
		args = append(args, "px", int64(expiration/time.Millisecond))
	} else {
		args = append(args, "ex", int64(expiration/time.Second))
	}
	return args
}

func (c *client) Get(key string, val interface{}) error {
	if err := c.GetContext(context.Background(), key, val); err != nil {
		return err
	}
	return nil
}

func (c *client) GetContext(ctx context.Context, key string, val interface{}) error {
	args := make([]interface{}, 2)
	args[0] = "get"
	args[1] = val
	result, err := c.roudtrip(ctx, "set", args...)
	if err != nil {
		return err
	}
	if err := c.options.Serializer.Unmarshal(result, val); err != nil {
		return err
	}
	return nil
}

func (c *client) roudtrip(ctx context.Context, cmd string, args ...interface{}) ([]byte, error) {
	conn, err := c.options.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer c.options.Pool.Put(conn)
	if err := c.authIfNeeded(ctx, conn); err != nil {
		return nil, err
	}
	result, err := c.roudtripWithConn(ctx, conn, cmd, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *client) roudtripWithConn(ctx context.Context, conn net.Conn, cmd string, args ...interface{}) ([]byte, error) {
	bytesArgs := make([][]byte, 0, len(args))
	for _, arg := range args {
		bytesArg, err := c.options.Serializer.Marshal(arg)
		if err != nil {
			return nil, err
		}
		bytesArgs = append(bytesArgs, bytesArg)
	}
	if err := c.options.Codec.EncodeTo(ctx, conn, cmd, bytesArgs); err != nil {
		return nil, err
	}
	result, err := c.options.Codec.Decode(ctx, conn)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *client) authIfNeeded(ctx context.Context, conn net.Conn) error {
	if c.options.Password == "" {
		return nil
	}
	if err := c.auth(ctx, conn); err != nil {
		return err
	}
	return nil
}

func (c *client) auth(ctx context.Context, conn net.Conn) error {
	cmd := "auth"
	args := []interface{}{c.options.Password}
	result, err := c.roudtripWithConn(ctx, conn, cmd, args...)
	if err != nil {
		return err
	}
	if !bytes.Equal(result, []byte("OK")) {
		return errors.New("auth failed")
	}
	return nil
}

func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}
