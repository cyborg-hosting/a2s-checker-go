##
## Build
##

FROM golang:1.21-bookworm AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /a2s-checker

##
## Deploy
##

FROM gcr.io/distroless/static

WORKDIR /

COPY --from=build --chmod=755 /a2s-checker /a2s-checker

ENTRYPOINT [ "/a2s-checker" ]
