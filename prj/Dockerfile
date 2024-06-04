FROM golang:1.21.5

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o prj

CMD ["./prj"]