##
## Build
##

FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /docker-a2s-checker

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /docker-a2s-checker /docker-a2s-checker

ENTRYPOINT [ "/docker-a2s-checker" ]
