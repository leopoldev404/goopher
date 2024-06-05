package main

import "context"

type store struct {
}

func NewStore() *store {
	return &store{}
}

func (store *store) Save(context.Context) error {
	return nil
}
