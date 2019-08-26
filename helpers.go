package main

import (
	"fmt"
	"os"
)

// getenv retrieves the string of the environment variable
func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		fmt.Println("environment variable is not set: " + name)
	}
	return v
}

func printSlice(s []byte) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
