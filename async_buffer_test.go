package async_buffer

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestAsyncBuffer1wr1rd(t *testing.T) {
	buf, exit := NewBuffer()
	var wg sync.WaitGroup
	wg.Add(2)
	// writer 1
	go func() {
		defer wg.Done()
		putTestData(buf)
		exit()
	}()
	// reader 1
	go func() {
		defer wg.Done()
		readTestData(buf)
	}()
	wg.Wait()
	t.Log("successfully tested 1 reader and 1 writer")
}

func TestAsyncBuffer2wr1rd(t *testing.T) {
	buf, exit := NewBuffer()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		var wrwg sync.WaitGroup
		wrwg.Add(2)
		// writer 1
		go func() {
			defer wrwg.Done()
			putTestData(buf)
		}()
		// writer 2
		go func() {
			defer wrwg.Done()
			putTestData(buf)
		}()
		wrwg.Wait()
		exit()
	}()
	// reader 1
	go func() {
		defer wg.Done()
		readTestData(buf)
	}()
	wg.Wait()
	t.Log("successfully tested 1 reader and 2 writers")
}

func TestAsyncBuffer1wr2rd(t *testing.T) {
	buf, exit := NewBuffer()
	var wg sync.WaitGroup
	wg.Add(3)
	// writer 1
	go func() {
		defer wg.Done()
		putTestData(buf)
		exit()
	}()
	// reader 1
	go func() {
		defer wg.Done()
		readTestData(buf)
	}()
	// reader 2
	go func() {
		defer wg.Done()
		readTestData(buf)
	}()
	wg.Wait()
	t.Log("successfully tested 2 readers and 1 writer")
}

func TestAsyncBuffer2wr2rd(t *testing.T) {
	buf, exit := NewBuffer()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		var wrwg sync.WaitGroup
		wrwg.Add(2)
		// writer 1
		go func() {
			defer wrwg.Done()
			putTestData(buf)
		}()
		// writer 2
		go func() {
			defer wrwg.Done()
			putTestData(buf)
		}()
		wrwg.Wait()
		exit()
	}()
	// reader 1
	go func() {
		defer wg.Done()
		readTestData(buf)
	}()
	// reader 2
	go func() {
		defer wg.Done()
		readTestData(buf)
	}()
	wg.Wait()
	t.Log("successfully tested 2 readers and 2 writers")
}

func TestRandomizedLoad(t *testing.T) {
	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().UnixNano())
		buf, exit := NewBuffer()
		var wg sync.WaitGroup
		nRand, gRand := rand.Intn(50)+1, rand.Intn(50)+1
		// g readers
		for g := 0; g < gRand; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				readTestData(buf)
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			var wrwg sync.WaitGroup
			for n := 0; n < nRand; n++ {
				wrwg.Add(1)
				go func() {
					defer wrwg.Done()
					putTestData(buf)
				}()
			}
			wrwg.Wait()
			exit()
		}()
		wg.Wait()
		t.Logf("successfully tested %d writers and %d readers\n", nRand, gRand)
	}
}

func putTestData(buffer *AsyncBuffer) {
	for i := 0; i < 100; i++ {
		buffer.Put(i)
	}
}

func readTestData(buffer *AsyncBuffer) {
	for {
		v := buffer.Get()
		if v == nil {
			break
		}
	}
}
