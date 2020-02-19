package async_buffer

import (
	"context"
	"sync"
)

type AsyncBuffer struct {
	mu       sync.Mutex
	isClosed uint8
	wg       sync.WaitGroup
	ctx      context.Context
	cache    []interface{}
	chOut    chan interface{}
	chIn     chan interface{}
	chToggle chan struct{}
}

func NewBuffer() (*AsyncBuffer, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	b := &AsyncBuffer{
		isClosed: 0,
		ctx:      ctx,
		chIn:     make(chan interface{}),
		chOut:    make(chan interface{}),
		chToggle: make(chan struct{}),
	}
	go b.bufferIn()
	go b.bufferOut()
	go b.listenClose()
	return b, cancel
}

func (b *AsyncBuffer) Get() interface{} {
	return <-b.chOut
}

func (b *AsyncBuffer) Put(v interface{}) {
	b.wg.Add(1)
	go func() {
		b.chIn <- v
	}()
}

func (b *AsyncBuffer) IsClosed() bool {
	return b.isClosed == 1
}

func (b *AsyncBuffer) listenClose() {
	defer close(b.chIn)
	defer func() {
		b.isClosed = 1
	}()
	<-b.ctx.Done()
	b.wg.Wait()
}

func (b *AsyncBuffer) bufferIn() {
	defer close(b.chToggle)
	for {
		v, ok := <-b.chIn
		if ok {
			b.push(v)
			b.chToggle <- struct{}{}
		} else {
			return
		}
	}
}

func (b *AsyncBuffer) bufferOut() {
	defer close(b.chOut)
	for {
		select {
		case _, ok := <-b.chToggle:
			if !ok {
				for {
					if !b.pullAndSend() {
						return
					}
				}
			}
		default:
			if !b.pullAndSend() {
				<-b.chToggle
			}
		}
	}
}

func (b *AsyncBuffer) pullAndSend() bool {
	v := b.pull()
	if v != nil {
		defer b.wg.Done()
		b.chOut <- v
		return true
	}
	return false
}

func (b *AsyncBuffer) push(v interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cache = append(b.cache, v)
}

func (b *AsyncBuffer) pull() interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.cache) > 0 {
		v := b.cache[0]
		b.cache = b.cache[1:]
		return v
	}
	return nil
}
