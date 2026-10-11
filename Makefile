.PHONY: all generate booklet clean test lint fmt

TEXINPUTS := .:../sty/:

all: out/book.pdf

out/book.pdf: generate sty/gloss.tex sty/pre.tex
	cd out && TEXINPUTS=$(TEXINPUTS) lualatex book.tex
	cd out && TEXINPUTS=$(TEXINPUTS) lualatex book.tex

generate: data/ templates/ cmd/songbook/main.go internal/generator/*.go go.mod go.sum
	go run ./cmd/songbook

booklet: out/booklet.pdf

out/booklet.pdf: out/book.pdf post-processing/booklet.tex
	cd post-processing && lualatex -output-directory=../out booklet.tex

clean:
	rm -rf out/

test:
	go test ./...

lint:
	test -z "$$(gofmt -l .)"
	go vet ./...

fmt:
	gofmt -w .
