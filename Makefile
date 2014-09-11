all: memdump memload

memdump: cmd/memdump/memdump.go
	go build -o $@ $< 

memload: cmd/memload/memload.go
	go build -o $@ $<

clean:
	rm -f memload memdump

