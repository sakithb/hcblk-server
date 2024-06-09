live/server:
	@wgo go run cmd/main/main.go

live/templ:
	@templ generate --watch --proxy="http://localhost:3000" --open-browser=false

live:
	@make -j2 live/server live/templ
