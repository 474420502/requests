package requests

import (
	"log"
	"sync"

	"github.com/schollz/progressbar"
)

type RequestPool struct {
	runnerCount int
	isBar       bool

	temps []*Temporary
	lock  sync.Mutex
	sem   chan int
}

func NewRequestPool(runnerCount int) *RequestPool {
	return &RequestPool{
		isBar:       false,
		runnerCount: runnerCount,
		sem:         make(chan int, runnerCount),
	}
}

func (pl *RequestPool) Add(tp *Temporary) {
	pl.lock.Lock()
	pl.temps = append(pl.temps, tp)
	pl.lock.Unlock()
}

func (pl *RequestPool) SetBar(is bool) {
	pl.lock.Lock()
	pl.isBar = is
	pl.lock.Unlock()
}

func (pl *RequestPool) Execute(errHandler func(int, error)) []*Response {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	respChan := make(chan *Response, len(pl.temps))
	var result []*Response

	var bar *progressbar.ProgressBar
	if pl.isBar {
		bar = progressbar.New(len(pl.temps))
	}

	for i, tp := range pl.temps {
		pl.sem <- 1

		go func(i int, tp *Temporary) {
			defer func() {
				<-pl.sem
				if pl.isBar {
					bar.Add(1) // 完成一个请求则进度+1
				}
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
		result = append(result, <-respChan)
	}

	return result
}
