package main

import (
	"github.com/Tansuozhe1num/codedream/internal/router"
)

func main() {
	r := router.NewRouter()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
