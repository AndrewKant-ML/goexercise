package main

import (
	"it.uniroma2.dicii/goexercise/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
