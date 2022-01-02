FROM golang:1.16 as builder
WORKDIR /go/src/cmg
COPY go.* ./
RUN go mod download
COPY . ./
RUN ARCH=$TARGETARCH make multiarch

FROM alpine:3
RUN apk --no-cache add ca-certificates
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/usr/bin/cmg"]
CMD ["serve"]
COPY --from=builder /go/src/cmg/bin/$TARGETARCH/cmg /usr/bin/cmg