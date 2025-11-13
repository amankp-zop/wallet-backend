# Wallet Backend

A robust and scalable wallet management system built with Go, providing user authentication, wallet management, and transaction processing capabilities.

## Features

- **User Management**
  - User signup and authentication
  - JWT-based authorization
  - Secure password hashing
  - User profile management

- **Wallet System**
  - Multi-currency wallet support
  - Balance management
  - Wallet creation and retrieval
  - User-wallet association

- **Transaction Processing**
  - Peer-to-peer transfers
  - Transaction status tracking (PENDING, COMPLETED, FAILED)
  - Transaction history

- **Background Worker**
  - Asynchronous task processing
  - Background job handling

## Technology Stack

- **Language**: Go 1.24.3
- **Web Framework**: Chi Router
- **Database**: MySQL 8.0
- **Cache**: Redis 7
- **Authentication**: JWT (JSON Web Tokens)
- **Configuration**: Viper
- **Validation**: go-playground/validator
- **Decimal Handling**: shopspring/decimal
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Go 1.24.3 or higher
- Docker and Docker Compose (for containerized deployment)
- MySQL 8.0 (for local development)
- Redis 7 (for local development)

## Installation

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/amankp-zop/wallet-backend.git
cd wallet-backend
```

2. Start all services using Docker Compose:
```bash
docker-compose up -d
```

This will start:
- MySQL database on port 3306
- Redis on port 6379
- API server on port 8080
- Background worker

3. Check service health:
```bash
docker-compose ps
```

### Local Development Setup

1. Clone the repository:
```bash
git clone https://github.com/amankp-zop/wallet-backend.git
cd wallet-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up MySQL database and run migrations:
```bash
# Create database
mysql -u root -p -e "CREATE DATABASE wallet;"

# Run migrations
mysql -u root -p wallet < migrations/000001_create_initial_tables.up.sql
```

4. Configure the application:
```bash
# Copy and edit the configuration file
cp configs/config.yaml configs/config.local.yaml
# Edit config.local.yaml with your local database credentials
```

5. Run the API server:
```bash
go run cmd/api/main.go
```

6. Run the worker (in a separate terminal):
```bash
go run cmd/worker/main.go
```

## Configuration

The application uses a YAML configuration file located at `configs/config.yaml`. You can create a local override at `configs/config.local.yaml` (which is git-ignored).

### Configuration Options

```yaml
server:
  port: 8080                                           # API server port

database:
  dsn: 'user:password@tcp(localhost:3306)/wallet?parseTime=true'  # MySQL connection string

redis:
  addr: 'localhost:6379'                               # Redis connection address

auth_config:
  jwt_secret: 'your-secret-key'                        # JWT signing secret
```

## API Endpoints

### Public Endpoints

#### Health Check
```
GET /health
```
Returns the health status of the API.

**Response:**
```json
{
  "status": "ok"
}
```

#### User Signup
```
POST /users/signup
```
Register a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### User Login
```
POST /users/login
```
Authenticate and receive a JWT token.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Protected Endpoints

These endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

#### Get User Profile
```
GET /users/profile
```
Retrieve the authenticated user's profile information.

#### Get User Wallet
```
GET /users/wallets
```
Retrieve the authenticated user's wallet information.

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "balance": "100.0000",
  "currency": "USD",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Database Schema

### Users Table
```sql
- id: BIGINT UNSIGNED (Primary Key)
- name: VARCHAR(255)
- email: VARCHAR(255) (Unique)
- password: VARCHAR(255) (Hashed)
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

### Wallets Table
```sql
- id: BIGINT UNSIGNED (Primary Key)
- user_id: BIGINT UNSIGNED (Foreign Key -> users.id)
- balance: DECIMAL(19,4)
- currency: VARCHAR(3) (Default: 'USD')
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

### Transactions Table
```sql
- id: BIGINT UNSIGNED (Primary Key)
- sender_wallet_id: BIGINT UNSIGNED (Foreign Key -> wallets.id)
- receiver_wallet_id: BIGINT UNSIGNED (Foreign Key -> wallets.id)
- amount: DECIMAL(19,4)
- status: ENUM('PENDING', 'COMPLETED', 'FAILED')
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

## Project Structure

```
wallet-backend/
├── cmd/                      # Application entry points
│   ├── api/                  # API server
│   │   └── main.go
│   └── worker/               # Background worker
│       └── main.go
├── configs/                  # Configuration files
│   └── config.yaml
├── internal/                 # Internal application code
│   ├── api/                  # API layer
│   │   ├── handler/          # HTTP handlers
│   │   └── middleware/       # HTTP middleware
│   ├── config/               # Configuration management
│   ├── database/             # Database connection
│   ├── domain/               # Domain models and interfaces
│   ├── repository/           # Data access layer
│   └── service/              # Business logic layer
├── migrations/               # Database migrations
│   ├── 000001_create_initial_tables.up.sql
│   └── 000001_create_initial_tables.down.sql
├── docker-compose.yml        # Docker Compose configuration
├── docker-compose.dev.yml    # Development Docker Compose
├── Dockerfile                # Docker build instructions
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
└── README.md                 # This file
```

## Development

### Building the Application

```bash
# Build API server
go build -o bin/api ./cmd/api

# Build worker
go build -o bin/worker ./cmd/worker
```

### Running Tests

```bash
go test ./...
```

### Docker Commands

```bash
# Build and start all services
docker-compose up --build

# Stop all services
docker-compose down

# View logs
docker-compose logs -f api
docker-compose logs -f worker

# Rebuild specific service
docker-compose up --build api
```

## Architecture

The application follows a clean architecture pattern with clear separation of concerns:

- **Handler Layer**: Handles HTTP requests and responses
- **Service Layer**: Contains business logic
- **Repository Layer**: Manages data access and persistence
- **Domain Layer**: Defines core business entities and interfaces

### Middleware

- **Logger**: Logs all incoming HTTP requests
- **Recoverer**: Recovers from panics and returns 500 errors
- **AuthMiddleware**: Validates JWT tokens for protected routes

## Security

- Passwords are hashed using bcrypt before storage
- JWT tokens are used for stateless authentication
- SQL injection protection through prepared statements
- CORS and security headers can be configured as needed

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues, questions, or contributions, please open an issue on GitHub.
