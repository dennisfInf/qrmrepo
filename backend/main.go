package main

import "github.com/enclaive/backend/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
