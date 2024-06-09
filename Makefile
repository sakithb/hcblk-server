.PHONY: build live/server live/templ live/css live

build:
	@go build -o builds/main cmd/main/main.go

live/server:
	@wgo DEV_MODE=1 go run cmd/main/main.go

live/templ:
	@templ generate --watch --proxy="http://localhost:3000" --open-browser=false

live/css:
	@tailwindcss -i assets/tailwind.css -o ./assets/dist/styles.css --minify --watch

live:
	@make -j3 live/templ live/css live/server
