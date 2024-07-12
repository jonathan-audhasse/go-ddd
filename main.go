package main

import (
	"goddd/src/api"
)

func main() {
	api.NewApi(api.NewConfig()).Serve()
}
