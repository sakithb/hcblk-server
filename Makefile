.PHONY: dev build

build:
	@go build -o builds/main cmd/main/main.go

dev:
	@wgo -file=.go -file=.templ -file=.css -xfile=_templ.go tailwindcss -i assets/tailwind.css -o assets/dist/styles.css :: templ generate :: go run cmd/main.go
