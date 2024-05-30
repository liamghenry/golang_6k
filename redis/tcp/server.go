package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"redis/interface/tcp"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func ListenAndServe(addr string, handler tcp.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logrus.Info("tcp server started on ", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	errCh := make(chan error, 1)

	go func() {
		select {
		case err := <-errCh:
			logrus.Error("tcp accept error: ", err)
		case sig := <-quit:
			logrus.Info("receive signal: ", sig)
		}
		listener.Close()
		handler.Close()
	}()

	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			errCh <- err
			break
		}
		logrus.Debug("accept connection: ", conn.RemoteAddr())

		ctx := context.Background()
		wg.Add(1)

		go func() {
			defer wg.Done()
			handler.Handle(ctx, conn)
		}()
	}

	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	select {
	case <-time.After(5 * time.Second):
		logrus.Warn("tcp server stopped, but some connections are still handling")
	case <-doneCh:
		logrus.Debug("all connections are handled")
	}

	logrus.Info("tcp server stopped")
	return nil
}
