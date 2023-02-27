FROM golang:alpine AS build

RUN apk update
RUN apk add ca-certificates git
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN mkdir ./out
RUN go build -o ./out/main ./app

FROM alpine:3
RUN mkdir /app
COPY --from=build /go/src/app/out/main /app
COPY --from=build /go/src/app/certs/verdandi.uno.crt /app
COPY --from=build /go/src/app/certs/verdandi.uno.middle /app
COPY --from=build /go/src/app/certs/verdandi.uno.root /app
COPY --from=build /go/src/app/certs/verdandi.uno.key /app
COPY --from=build /go/src/app/certs/verdandi.uno.pem /app

EXPOSE 443
EXPOSE 80

ENV ENV prod
ENV PORT 443

ENTRYPOINT ["/app/main"]