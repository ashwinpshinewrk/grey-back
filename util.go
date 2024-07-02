package main

import (
	"fmt"
)

// Handles errors
func Handle_error(e error) {
	if e != nil {
		fmt.Println(e)
		return
	}
}
