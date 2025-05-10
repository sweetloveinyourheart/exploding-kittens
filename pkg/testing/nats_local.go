package testing

import (
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	server "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
)

// Timeout for request/response.
const Timeout = 5 * time.Second

func init() {
	for range 10 {
		NATSServerPool.Put(NATSServerPool.New())
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func newServer() *server.Server {
	s, err := server.NewServer(&server.Options{
		Host:           "127.0.0.1",
		Port:           server.RANDOM_PORT,
		NoLog:          true,
		NoSigs:         true,
		MaxControlLine: 2048,
		JetStream:      true,
		StoreDir:       filepath.Join(os.TempDir(), server.JetStreamStoreDir, RandStringRunes(5)),
	})

	if err != nil {
		return nil
	}

	go func() {
		err := server.Run(s)
		if err != nil {
			log.Global().Error("starting nats server", zap.Error(err))
		}
	}()

	return s
}

var NATSServerPool = sync.Pool{
	New: func() any {
		return newServer()
	},
}

func StartLocalNATSServer(t *testing.T) (s *server.Server, shutdown func()) {
	t.Helper()

	s = NATSServerPool.Get().(*server.Server)

	if s == nil {
		t.Fatal("starting nats server: nil")
	}

	if !s.ReadyForConnections(Timeout) {
		t.Fatal("starting nats server: timeout")
	}
	return s, func() {
		s.Shutdown()

		NATSServerPool.Put(newServer())
	}
}

func NATSWaitConnected(t *testing.T, c *nats.Conn) {
	t.Helper()

	timeout := time.Now().Add(Timeout)
	for time.Now().Before(timeout) {
		if c.IsConnected() {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("client connecting timeout")
}
func StartLocalNATSServerBench(t *testing.B) (s *server.Server, shutdown func()) {
	t.Helper()

	s = NATSServerPool.Get().(*server.Server)

	if s == nil {
		t.Fatal("starting nats server: nil")
	}

	if !s.ReadyForConnections(Timeout) {
		t.Fatal("starting nats server: timeout")
	}
	return s, func() {
		s.Shutdown()

		NATSServerPool.Put(newServer())
	}
}

func NATSWaitConnectedBench(t *testing.B, c *nats.Conn) {
	t.Helper()

	timeout := time.Now().Add(Timeout)
	for time.Now().Before(timeout) {
		if c.IsConnected() {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("client connecting timeout")
}
