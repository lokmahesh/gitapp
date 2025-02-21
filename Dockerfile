FROM golang:1.20

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o github-app

CMD ["./github-app"]
EXPOSE 3000
