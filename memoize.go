// Utilities for memoizing expensive operations (like loading corpus data).

package dorkalonius

// Memoizes a single result. This struct consists solely of a channel used
// to communicate requests to a worker goroutine.
type Memo struct {
	requestChan chan<- chan<- interface{}
}

func NewMemo(f func() interface{}) Memo {
	requestChan := make(chan chan<- interface{})
	go worker(f, requestChan)
	return Memo{requestChan}
}

// Fetches the memoized object.
func (self Memo) Get() interface{} {
	response := make(chan interface{})
	self.requestChan <- response
	return <-response
}

// The worker goroutine.
func worker(f func() interface{}, requestChan chan chan<- interface{}) {
	var called bool = false
	var result interface{} = nil

	for response := range requestChan {
		if !called {
			// We construct the object lazily.
			result = f()
			called = true
		}
		response <- result
	}
}
