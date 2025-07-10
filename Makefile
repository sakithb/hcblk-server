.PHONY: dev build

data:
	@mkdir -p data/logs data/db
	@sqlite3 assets/db/data.db "DROP TABLE IF EXISTS bikes;" 
	@sqlite3 assets/db/data.db "CREATE TABLE bikes(id TEXT, brand TEXT, model TEXT, year INT, category TEXT, engine_capacity INT, PRIMARY KEY(id));" 
	@sqlite3 assets/db/data.db ".import assets/bikes.csv bikes --csv --skip 1" 
	@sqlite3 assets/db/data.db "DROP TABLE IF EXISTS cities;" 
	@sqlite3 assets/db/data.db "CREATE TABLE cities(id TEXT, city TEXT, district TEXT, province TEXT, PRIMARY KEY(id));" 
	@sqlite3 assets/db/data.db ".import assets/cities.csv cities --csv --skip 1" 

build:
	@mkdir -p builds
	@npx tailwindcss@3 -i assets/tailwind.css -o assets/dist/styles.css
	@TEMPL_EXPERIMENT=rawgo templ generate
	@go build -o ./builds/main ./cmd/main.go

dev:
	@wgo -file=.go -file=.templ -file=.css -xfile=_templ.go npx tailwindcss@3 -i assets/tailwind.css -o assets/dist/styles.css :: TEMPL_EXPERIMENT=rawgo templ generate :: DEV_MODE=1 go run ./cmd/main.go
