package main

import (
	"fmt"
	api "proglog/api/v1/log"
)

func main() {

	rec := api.Record{
		Value:  []byte("Hello"),
		Offset: 18,
	}

	fmt.Println(rec)
}
