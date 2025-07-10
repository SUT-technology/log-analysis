FROM golang:1.23.6-alpine

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy everything in backend
COPY . .

# Set working directory to cmd where main.go lives
WORKDIR /app/cmd

# Build the application
RUN go build -o /app/main .

EXPOSE 8000
CMD ["/app/main"]
