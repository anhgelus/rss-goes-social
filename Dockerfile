FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o rss-goes-social .

CMD ./rss-goes-social run
