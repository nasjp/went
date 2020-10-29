CC = go

SRCS=$(wildcard *.go)

went: $(SRCS)
	rm -rf tmp*
	$(CC) build .

test: went
	./test.sh

clean:
	rm -f went tmp*

.PHONY: test clean
