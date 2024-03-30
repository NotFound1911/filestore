package main

import "github.com/NotFound1911/filestore/api"

func main() {
	r := api.NewRouter()
	r.Run(":8888")
}
