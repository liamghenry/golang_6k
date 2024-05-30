package tcp

import (
	"context"
	"net"
)

// Handler is a interface that handle tcp conn
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
