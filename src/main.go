package main

import (
	"goddd/api"
)

func main() {
	api.NewApi(api.NewConfig()).Serve()
}
