FROM golang:1.23-alpine

WORKDIR /app

# Install git (required for Go modules)
RUN apk add --no-cache git

# Copy module files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Run tests by default when the container starts
CMD ["go", "test", "./..."]
