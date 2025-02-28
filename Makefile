# Simple Makefile for a Go project

wire:
	@cd ./cmd/api && wire

# Run the application
run:
	@go run ./cmd/api
