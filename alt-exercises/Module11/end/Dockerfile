ARG TARGETOS
ARG TARGETARCH
FROM golang:1.14.3-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go mod download
WORKDIR /src/cmd
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/gameserver .

RUN mkdir /app/keys
RUN mkdir /app/configs
RUN apk upgrade --update-cache --available && \
    apk add openssl && \
    rm -rf /var/cache/apk/*
RUN openssl genrsa -out /app/keys/app.rsa
RUN openssl rsa -in /app/keys/app.rsa -pubout > /app/keys/app.rsa.pub
RUN cp ../configs/app.env /app/configs

FROM scratch AS bin-unix
COPY --from=build /app /gameserverapp

FROM bin-unix AS bin-linux
FROM bin-unix AS bin-darwin

FROM scratch AS bin-windows
COPY --from=build /app /gameserverapp

FROM bin-${TARGETOS} as bin
ENV GAME_SERVER_HOMEDIR="/gameserverapp"
ENTRYPOINT ["/gameserverapp/gameserver"]
