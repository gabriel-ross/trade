FROM golang

WORKDIR /app/server

ENV PORT=80

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p bin
RUN go build -o ./bin/server ./cmd/main.go

EXPOSE ${PORT}

CMD ./bin/server