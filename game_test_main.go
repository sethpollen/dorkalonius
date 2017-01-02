// A small program to use for manual evaluation and tuning of the game.go
// module.

package main

import (
	"fmt"
	"github.com/sethpollen/dorkalonius"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 40; i++ {
		fmt.Println(dorkalonius.NewTargetWord())
	}
}
