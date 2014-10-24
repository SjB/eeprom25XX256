all: memdump memload

memdump: cmd/memdump/main.go
	go build -o $@ $< 

memload: cmd/memload/main.go
	go build -o $@ $<

clean:
	rm -f memload memdump

