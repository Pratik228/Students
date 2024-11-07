package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
    var wg sync.WaitGroup
    fmt.Println("Testing concurrent student creation...")
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            student := map[string]interface{}{
                "id":    id,
                "name":  fmt.Sprintf("Student %d", id),
                "age":   20,
                "email": fmt.Sprintf("student%d@example.com", id),
            }
            
            jsonData, _ := json.Marshal(student)
            resp, err := http.Post(
                "http://localhost:8080/students",
                "application/json",
                bytes.NewBuffer(jsonData),
            )
            if err != nil {
                fmt.Printf("Error creating student %d: %v\n", id, err)
            } else {
                fmt.Printf("Created student %d, status: %d\n", id, resp.StatusCode)
            }
        }(i)
        time.Sleep(100 * time.Millisecond) 
    }

    fmt.Println("\nTesting concurrent student retrieval...")
    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func(iteration int) {
            defer wg.Done()
            resp, err := http.Get("http://localhost:8080/students")
            if err != nil {
                fmt.Printf("Error reading students (iteration %d): %v\n", iteration, err)
            } else {
                fmt.Printf("Read students (iteration %d), status: %d\n", iteration, resp.StatusCode)
            }
        }(i)
    }

    wg.Wait()
    fmt.Println("\nConcurrency test completed!")
}