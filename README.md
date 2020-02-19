# async buffer
AsyncBuffer is a lightweight solution for asynchronously sending and retrieving data. 
The basic idea is that writers to the buffer should not block when attempting to write new values, 
but readers should block until either values are available or the buffer is closed via the cancel function.

At its core, the buffer offers two methods:
```go
// blocking call to retrieve data
func (b *AsyncBuffer) Get() interface{} 

// non-blocking call to put new data
func (b *AsyncBuffer) Put(v interface{}) 
```

When a new buffer is allocated, a context.CancelFunc is returned with it. The intended purpose of the cancel function
is to tell the buffer when no more writes are to be expected and to wait until all written values have been successfully retrieved.
Once a buffer has been closed, all attempts to retrieve new values using Get() will return nil values.
```go
func NewBuffer() (*AsyncBuffer, context.CancelFunc) 
```

Below is an example showcasing the use of an async buffer with 1 reader and 1 writer. For an example of multiple readers
and writers, please review the TestAsyncBuffer2wr2rd function located in `async_buffer_test.go`.
```go
func main() {
    // initialize a new buffer with its corresponding cancel function
	buf, cancel := NewBuffer()
    // since we know we will have 1 reader and 1 writer, initialize a waitgroup with a count of 2
	var wg sync.WaitGroup
	wg.Add(2)
	// writer 1
	go func() {
		defer wg.Done()
		putTestData(buf)
        // once all values have been sent to the buffer to be queued, wait until they are received by using the cancel function.
		cancel()
	}()
	// reader 1
	go func() {
		defer wg.Done()
        // asynchronously read values from the buffer
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
		if v == nil && buffer.IsClosed(){
			break
		} else {
            // handle nil value
        }
		fmt.Printf("get: %v\n", v)
	}
}
```
Notice the use of the third and final method the buffer offers, `IsClosed()`. This method allows the user to differentiate
between nil values written to the buffer and nil values returned by the buffer after it has been closed.