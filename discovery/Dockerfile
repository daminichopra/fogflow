FROM golang:1.16 AS builder 
WORKDIR /go/src/
RUN git clone https://github.com/smartfog/fogflow.git  
WORKDIR /go/src/fogflow/discovery
RUN go get; go build;

FROM scratch
COPY --from=builder /go/src/fogflow/discovery/discovery  /
CMD ["/discovery"]

