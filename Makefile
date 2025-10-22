# Valforge Makefile

.PHONY: build test install clean

build:
	go build -o valforge .

install:
	go install valforge

fmt:
	go fmt ./...

vet:
	go vet ./...

.DEFAULT_GOAL := build
