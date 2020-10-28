went: *.go
	rm -rf tmp*
	go build .

test: went
	./test.sh

clean:
	rm -f went tmp*

.PHONY: test clean
