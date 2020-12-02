all: clean png

dot:
	go run main.go
png: dot
	cd go-diagrams && dot network.dot -Tpng > network.png
clean:
	rm -rf go-diagrams