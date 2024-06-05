package main

import "context"

func main() {
	var store = NewStore()
	var orderService = NewService(store)

	orderService.Create(context.Background())
}
