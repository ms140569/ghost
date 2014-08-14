package main

type Name int

const (
	COMMAND Name = iota
	HEADER
	DATA
)

type Token struct {
	name  Name
	value interface{}
}
