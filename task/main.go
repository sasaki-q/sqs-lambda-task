package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("Start task \n")
	fmt.Printf("Message ID === %s \n", os.Getenv("MessageId"))
	fmt.Print("End task \n")
}
