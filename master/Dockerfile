FROM golang:1.16 AS builder 
WORKDIR /go/src/
RUN git clone https://github.com/smartfog/fogflow.git  
WORKDIR /go/src/fogflow/master
RUN go get; go build;

FROM scratch
COPY --from=builder /go/src/fogflow/master/master  /
CMD ["/master"]
