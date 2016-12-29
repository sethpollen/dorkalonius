// Utilities for memoizing expensive operations (like loading corpus data).

package dorkalonius

// Memoizes a single result. This struct consists solely of a channel used
// to communicate requests to a worker goroutine.
type Memo struct {
	requestChan chan<- chan<- result
}

func NewMemo(f func() (interface{}, error)) Memo {
	requestChan := make(chan chan<- result)
	go worker(f, requestChan)
	return Memo{requestChan}
}

type result struct {
	Object interface{}
	Err    error
}

// Fetches the memoized object.
func (self Memo) Get() (interface{}, error) {
	response := make(chan result)
	self.requestChan <- response
	result := <-response
	return result.Object, result.Err
}

// The worker goroutine.
func worker(f func() (interface{}, error),
	requestChan chan chan<- result) {
	var called bool = false
	var result result

	for response := range requestChan {
		if !called {
			// We construct the object lazily.
			result.Object, result.Err = f()
			called = true
		}
		response <- result
	}
}
