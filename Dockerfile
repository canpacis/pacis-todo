FROM oven/bun AS frontend-deps

WORKDIR /web

COPY ./bun.lock ./package.json ./

RUN bun install --frozen-lockfile

FROM oven/bun AS frontend-build

WORKDIR /web

COPY --from=frontend-deps web/node_modules ./node_modules

COPY . .

RUN bun run build

FROM golang:1.24-alpine AS backend-deps

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24-alpine AS backend-build

WORKDIR /app

COPY --from=backend-deps /go/pkg /go/pkg
COPY . .

COPY --from=frontend-build /web/build ./build

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -tags=prod -o server .

RUN apk add --no-cache upx && upx --best --lzma ./server

FROM scratch

COPY --from=backend-build /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
