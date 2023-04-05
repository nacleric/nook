.PHONY: run
run: mymain
	./$<

mymain: *.go go.mod
	go build -o $@ .

.PHONY: all
all: mymain