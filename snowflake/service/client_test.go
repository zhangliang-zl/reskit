package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
	"time"
)

func startServer(opts ServerOptions) error {
	logger := log.NewHelper(log.With(log.DefaultLogger, "tag", "uuid"))
	ser, err := NewServer(opts, logger)
	if err != nil {
		return err
	}
	if err := ser.Run(); err != nil {
		return err
	}

	return nil
}

func TestClient_UUID(t *testing.T) {
	opts1 := ServerOptions{
		WorkerID:   1,
		WorkerBits: 4,
		NumberBits: 18,
		Epoch:      time.Now().Unix() - 86400,
		Port:       5001,
	}

	opts2 := ServerOptions{
		WorkerID:   2,
		WorkerBits: 4,
		NumberBits: 18,
		Epoch:      time.Now().Unix() - 86400,
		Port:       5002,
	}

	go startServer(opts2)
	go startServer(opts1)

	time.Sleep(time.Millisecond * 100)
	c := NewClient("127.0.0.1:5001", "127.0.0.1:5002")

	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond)
		if _, err := c.UUID(); err != nil {
			t.Fatal(err)
		}
	}
}
