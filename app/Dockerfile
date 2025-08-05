# Yeah so this image might be a bit massive
# Ideally after running tests we two-step build or find another way to test
# But this is a project to learn GCP not Docker w/ Go so it's fine for now

FROM golang:1.24.5-alpine

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
