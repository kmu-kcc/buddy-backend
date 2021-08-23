###################
### build stage ###
###################

FROM golang:1.17 AS builder
WORKDIR /go/src/github.com/kmu-kcc/buddy-backend/
RUN go get github.com/kmu-kcc/buddy-backend && \
  CGO_ENABLED=0 go build -gcflags -m -o buddy -v .

###################
#### run stage ####
###################

FROM alpine:latest
WORKDIR /bin/
COPY --from=builder /go/src/github.com/kmu-kcc/buddy-backend/buddy ./
RUN chmod +x ./buddy
EXPOSE 3000
CMD ["./buddy", "--port", "3000"]
