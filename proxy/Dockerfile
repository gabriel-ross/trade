FROM golang

WORKDIR /app

ENV PORT=81
ENV SERVER_ADDRESS=localhost:80

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p bin
RUN go build -o ./bin/proxy cmd/proxy

EXPOSE ${PORT}

CMD ./bin/proxy