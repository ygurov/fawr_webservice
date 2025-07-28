package main

import (
	"github.com/fawrwebservice/api"
	"github.com/fawrwebservice/storage"
)

func main() {
	api.Register(":80", storage.NewDB())
}
