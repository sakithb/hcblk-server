FROM golang:1.22

WORKDIR "/app"
COPY ./ ./

RUN apt update -y && apt install -y sqlite3 nodejs npm
RUN npm install -g tailwindcss
RUN go mod download
RUN go mod verify
RUN make build
RUN make data

EXPOSE 8080
CMD ["./builds/main"]
