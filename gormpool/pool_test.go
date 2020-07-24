package gormpool

import (
	"context"
	"fmt"
	godbpool "github.com/ALiuGuanyan/go-db-pool"
	"sync"
	"testing"
	"time"
)

func TestMySQLNewPool(t *testing.T) {
	ctx, canc := context.WithCancel(context.Background())
	opts := Options{
		Type:            godbpool.MySQL,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		KeepConn:        2,
		Capacity:        5,
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	p, err := NewPool(ctx, opts)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 5*time.Second)
	}
	wg.Wait()
	canc()
}

func TestMySQLClose(t *testing.T) {
	ctx, canc := context.WithCancel(context.Background())
	opts := Options{
		Type:            godbpool.MySQL,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		KeepConn:        2,
		Capacity:        5,
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	p, err := NewPool(ctx, opts)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 2*time.Second)
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 8*time.Second)
	}
	time.Sleep(4 * time.Second)
	canc()
	wg.Wait()
}

func mockJob(wg *sync.WaitGroup, p *Pool, duration time.Duration) {
	conn, err := p.Get()
	if err != nil {
		fmt.Println(err)
	} else {
		time.Sleep(duration)
		p.Put(conn)
	}
	wg.Done()
}
