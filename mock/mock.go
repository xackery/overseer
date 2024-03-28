package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// w, err := os.Create("mock.log")
	// if err != nil {
	// 	fmt.Println("Failed to create mock.log")
	// 	return
	// }
	// defer w.Close()
	// os.Stdout = w
	fmt.Println("Starting")
	for {
		// 30% chance of an error
		if rand.Intn(30) == 0 {
			fmt.Println("/src/zone/main.cpp:1:1 [Error] Something went wrong")
		} else {
			fmt.Println("/src/zone/main.cpp:1:2 Everything is fine")
		}

		// 10% chance of exiting
		if rand.Intn(10) == 0 {
			// 50% chance with an error
			if rand.Intn(2) == 0 {
				fmt.Println("/src/zone/main.cpp:1:1 [Error] Something went wrong")
				//	operation.Exit(1)
			}

			fmt.Println("Exiting")
			return
		}

		time.Sleep(time.Duration(rand.Intn(5000) * int(time.Millisecond)))

	}
}
