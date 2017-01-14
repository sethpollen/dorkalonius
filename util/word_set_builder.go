package util

// Runs all of the 'tasks' in parallel and returns a WordSet containing their
// combined outputs.
func BuildWordSet(tasks []func() WordSet) WordSet {
	responseChans := make([]chan WordSet, len(tasks))
	for i := range tasks {
		responseChans[i] = make(chan WordSet)
		go func(task func() WordSet, responseChan chan<- WordSet) {
			responseChan <- task()
		}(tasks[i], responseChans[i])
	}

	// Collect outputs from workers.
	wordSet := NewWordSet()
	for _, responseChan := range responseChans {
		wordSet.AddAll(<-responseChan)
	}
	return wordSet
}
