package main

import "context"

type OrdersService interface {
	Create(context.Context) error
}

type OrdersStore interface {
	Save(context.Context) error
}
