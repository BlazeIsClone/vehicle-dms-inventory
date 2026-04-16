FROM golang:1.22.1

WORKDIR /app

ADD . .

RUN make build

EXPOSE 3000

CMD ["./build"]