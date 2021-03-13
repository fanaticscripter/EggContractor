.PHONY: all go alert update

all: alert

alert: go
	./monitor --error-on-diff

update: go
	./monitor

go:
	cd cmd/monitor && go build -o ../../monitor
