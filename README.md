# PolyglotWeb-Go

Declarative Fullstack (Backend/Frontend) Framework using Htmx Like frontend and React Like backend implemented using Go

## Local Development

1. Copy the example environment file and modify as needed:
```bash
cp .env.example .env
```

2. Run the application:
```bash
go run .
```

3. Regenerate swagger when there are changes (see `generate.go` for required packages):
```bash
go generate .
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
