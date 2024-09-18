FROM golang:1.17 AS build
WORKDIR /go/src/cmg
ARG TARGETARCH
ARG TARGETOS
COPY go.* ./
RUN go mod download
COPY . ./
RUN make multiarch

FROM alpine:3
RUN apk --no-cache add ca-certificates
EXPOSE 8080
ENV PORT=8080
ENV ROOT_PATH="/"
ENTRYPOINT ["/usr/bin/cmg"]
CMD ["serve"]
COPY --from=build /go/src/cmg/bin/cmg /usr/bin/cmg