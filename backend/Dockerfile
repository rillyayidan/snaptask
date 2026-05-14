FROM golang:1.25-bookworm AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /snaptask-api ./cmd

FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=build /snaptask-api /snaptask-api
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/snaptask-api"]
