FROM golang:1.23-alpine
RUN go get -u github.com/cosmtrek/air
WORKDIR /app

# Copy go.mod and go.sum first, to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source
COPY . .

# Install air
RUN go install github.com/cosmtrek/air@latest

CMD ["air"]
