# Chirpy

**Chirpy** is a lightweight HTTP-based social microblogging API server written in Go. Built for learning and simplicity, Chirpy provides core functionality such as user registration, login, posting short messages ("chirps"), and secure token-based authentication. This project is designed as part of the [Boot.dev](https://boot.dev) HTTP Servers course.

---

## Features

- User registration and login  
- Secure password hashing with bcrypt  
- Token-based authentication with refresh support  
- Create and retrieve chirps  
- Rate-limiting endpoint tracking  
- Basic XSS protection for chirp content  
- Health check endpoint

---

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/chirpy.git
cd chirpy
```

Install dependencies:

```bash
go mod tidy
```

Run the server:

```bash
go run main.go
```

By default, the server will start on `localhost:8080`.

---

## API Endpoints

### Authentication

- `POST /api/users` - Register a new user  
- `POST /api/login` - Login and receive tokens  
- `POST /api/refresh` - Refresh access token  

### Chirps

- `POST /api/chirps` - Post a new chirp  
- `GET /api/chirps` - Retrieve all chirps  

### Health Check

- `GET /api/healthz` - Returns 200 OK if the server is healthy

---

## Folder Structure

```
chirpy/
├── internal/         # Project internals (e.g., database logic)
├── sql/              # SQL schema and migrations
├── handler_*.go      # HTTP route handlers
├── users.go          # User logic
├── chirpy/           # Project module name
├── main.go           # Entry point
└── ...
```

---

## Development

To run tests:

```bash
go test ./...
```

To lint:

```bash
golint ./...
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contributing

PRs are welcome! This project is educational in nature, but contributions that improve structure, security, or testing are appreciated.

---

## Credits

Built as part of the Boot.dev Go HTTP Server curriculum. Shoutout to the Boot.dev team for structured backend education.
