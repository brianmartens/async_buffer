package async_buffer

import (
	"fmt"
	"sync"
	"testing"
)

func TestAsyncBuffer1wr1rd(t *testing.T) {
	defer func() {
		fmt.Println(recover())
	}()
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
}

func TestAsyncBuffer2wr1rd(t *testing.T) {
	defer func() {
		fmt.Println(recover())
	}()
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
}

func TestAsyncBuffer1wr2rd(t *testing.T) {
	defer func() {
		fmt.Println(recover())
	}()
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
}

func TestAsyncBuffer2wr2rd(t *testing.T) {
	defer func() {
		fmt.Println(recover())
	}()
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
}

func putTestData(buffer *AsyncBuffer) {
	for i := 0; i < 100; i++ {
		buffer.Put(i)
		fmt.Printf("put: %v\n", i)
	}
}

func readTestData(buffer *AsyncBuffer) {
	for {
		v := buffer.Get()
		if v == nil && buffer.IsClosed() {
			break
		}
		fmt.Printf("get: %v\n", v)
	}
}
