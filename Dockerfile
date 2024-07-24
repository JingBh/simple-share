FROM node:lts AS web-builder

WORKDIR /app

COPY ./web/package.json ./web/yarn.lock ./

RUN yarn install --frozen-lockfile

COPY ./web .

RUN yarn build

FROM golang:1.22 AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /go/src/github.com/jingbh/simple-share

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY --from=web-builder /app/dist ./web/dist

RUN go build -a -o /go/bin/simple-share .

FROM gcr.io/distroless/base-debian12:nonroot AS runner

COPY --from=builder --chown=nonroot:nonroot /go/bin/simple-share /usr/bin/simple-share

EXPOSE 8080

ENTRYPOINT ["/usr/bin/simple-share"]
