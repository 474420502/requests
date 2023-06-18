package requests

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/schollz/progressbar"
)

type RequestPool struct {
	runnerCount int
	timeout     time.Duration

	temps []*Temporary
	lock  sync.Mutex

	sem chan int
}

func NewRequestPool(runnerCount int) *RequestPool {
	return &RequestPool{
		runnerCount: runnerCount,
		sem:         make(chan int, runnerCount),
	}
}

func (pl *RequestPool) Add(tp *Temporary) {
	pl.lock.Lock()
	pl.temps = append(pl.temps, tp)
	pl.lock.Unlock()
}

func (pl *RequestPool) SetPerTimeout(dur time.Duration) {
	pl.timeout = dur
}

func (pl *RequestPool) Execute(errHandler func(int, error)) []*Response {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), pl.timeout)
	defer cancel()

	respChan := make(chan *Response, len(pl.temps))
	var result []*Response

	timeout := time.After(pl.timeout)

	bar := progressbar.New(len(pl.temps))

	for i, tp := range pl.temps {
		pl.sem <- 1
		ctx, cancel := context.WithTimeout(ctx, pl.timeout)
		go func(i int, tp *Temporary) {
			defer func() {
				<-pl.sem
				bar.Add(1) // 完成一个请求则进度+1
			}()
			r, err := tp.Execute()
			if err != nil {
				if errHandler != nil {
					errHandler(i, err)
				} else {
					log.Println(i, err)
				}
			} else {
				respChan <- r
			}
		}(i, tp)
	}

	// 从channel中接收响应,超时则退出循环
	for range pl.temps {
		select {
		case r := <-respChan:
			result = append(result, r)
		case <-timeout:
			return result
		}
	}

	return result
}
