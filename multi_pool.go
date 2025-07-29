package requests

import (
	"sync"

	"github.com/schollz/progressbar"
)

type RequestPool struct {
	runnerCount int
	isBar       bool

	requests []*Request
	lock     sync.Mutex
	sem      chan int
}

type MultiResponse struct {
	*Response
	Error error
}

func NewRequestPool(runnerCount int) *RequestPool {
	return &RequestPool{
		isBar:       false,
		runnerCount: runnerCount,
		sem:         make(chan int, runnerCount),
	}
}

func (pl *RequestPool) Add(req *Request) {
	pl.lock.Lock()
	pl.requests = append(pl.requests, req)
	pl.lock.Unlock()
}

func (pl *RequestPool) SetBar(is bool) {
	pl.lock.Lock()
	pl.isBar = is
	pl.lock.Unlock()
}

func (pl *RequestPool) Execute() []*MultiResponse {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	respChan := make(chan *MultiResponse, len(pl.requests))
	var result []*MultiResponse

	var bar *progressbar.ProgressBar
	if pl.isBar {
		bar = progressbar.New(len(pl.requests))
	}

	for i, req := range pl.requests {
		pl.sem <- 1

		go func(i int, req *Request) {
			defer func() {
				<-pl.sem
			}()

			resp, err := req.Execute()
			respChan <- &MultiResponse{Response: resp, Error: err}

			if bar != nil {
				bar.Add(1)
			}
		}(i, req)
	}

	for range pl.requests {
		result = append(result, <-respChan)
	}

	return result
}
