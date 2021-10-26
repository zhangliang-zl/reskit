package cron

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	c := cron.New(cron.WithSeconds())
	i := 0
	runTimes := 5

	locker := sync.Mutex{}
	id, err := c.AddFunc("* * * * * *", func() {
		locker.Lock()
		i++
		fmt.Println(i)
		locker.Unlock()
	})
	if err != nil {
		t.Fatal(id, err)
	}

	srv := NewServer(c)
	ticker := time.NewTimer(time.Second * time.Duration(runTimes))

	go func() {
		err = srv.Start(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	}()

loop:
	for {
		select {
		case <-ticker.C:
			err := srv.Stop(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			break loop
		}
	}

	if i != runTimes {
		t.Fatalf("cron task err, i is not %d", runTimes)
	}
}
