FROM quay.io/cybozu/golang:1.19-focal as builder

WORKDIR /workspace
COPY ./client .
RUN go build -trimpath .

FROM quay.io/cybozu/ubuntu:20.04

COPY --from=builder /workspace/client /
USER 10000:10000
ENTRYPOINT ["/client"]
