FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o gts-rss .

CMD ./gts-rss
