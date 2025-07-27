# herdcontrol

`herdcontrol` is a Go package that prevents request stampedes by coalescing concurrent requests with the same key into a single function call. This ensures that expensive or rate-limited operations run only once even if requested concurrently.

---

## Features

- Request coalescing / single-flight control
- Prevent duplicate execution of concurrent requests with the same key
- Simple API with `Group.Do`
- Lightweight and efficient

---

## Installation

```bash
go get github.com/prabhavdogra/herdcontrol
```

## Example
```go
package main

import (
	"fmt"
	"time"

	herdcontrol "github.com/prabhavdogra/herdcontrol"
)

func main() {
	var g herdcontrol.Group
	key := "data"

	goroutine := 15
	results := make(chan string, goroutine)

	for i := 0; i < goroutine; i++ {
		go func(i int) {
			v, err := g.Do(key, func() (interface{}, error) {
				time.Sleep(500 * time.Millisecond) // simulate work
				fmt.Println("function executed")
				return fmt.Sprintf("value-%d", i), nil
			})

			if err == nil {
				results <- v.(string)
			}
		}(i)
	}

	for i := 0; i < goroutine; i++ {
		fmt.Println(<-results)
	}
}
```
