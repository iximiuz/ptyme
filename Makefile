.PHONY
shim:
	cc -Wall -std=c99 shim.c -o build/shim

.PHONY
attach:
	go run attach.go
