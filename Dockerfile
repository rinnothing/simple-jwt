FROM golang:1.24

WORKDIR /usr/src/server
COPY . .

RUN make test
RUN make build
CMD ["./server"]
