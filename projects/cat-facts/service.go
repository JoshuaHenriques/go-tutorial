package main

import "context"

type Service interface {
	GetCatFace(context.Context) (*CatFact, error)
}
