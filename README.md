# PolyglotWeb-Go

Declarative Fullstack (Backend/Frontend) Framework using Htmx Like frontend and React Like backend implemented using Go

## Local Development

1. Install dependencies:
```bash
go mod tidy
```

2. Copy the example environment file and modify as needed:
```bash
cp .env.example .env
```

3. Run the application:
```bash
go run .
```

By default, the server will start on port 3000.

## Docker

Build the image:
```bash
docker build -t polyglot-web-go .
```

Run the container:
```bash
docker run -p 3000:3000 -v ./data:/app/data --env DB_PATH=/app/data/data.db polyglot-web-go
```
