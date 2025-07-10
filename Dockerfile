FROM golang:1.22

RUN apt update -y && apt install -y sqlite3 nodejs npm
RUN go install github.com/a-h/templ/cmd/templ@v0.2.747

WORKDIR "/app"
COPY ./ ./

RUN go mod download
RUN go mod verify
RUN make build
RUN make data

EXPOSE 8080
CMD ["./builds/main"]
