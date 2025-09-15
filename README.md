# RBAC Admin Server

A modern Role-Based Access Control (RBAC) admin server built with Go, supporting multiple database connections with factory pattern design.

## ✨ Features

- 🏭 **Database Factory Pattern**: Support for MySQL, PostgreSQL, SQLite, SQL Server
- 🔧 **Enterprise Configuration**: 11 configuration modules with environment variable override
- 🔄 **Auto Migration**: Automatic database table creation and updates
- 🔐 **JWT Authentication**: Secure user authentication and permission management
- 📊 **Structured Logging**: High-performance logging system based on Zap
- 🐳 **Docker Support**: Complete containerized deployment solution
- 📈 **Monitoring**: Prometheus metrics collection and health checks
- 📚 **API Documentation**: Auto-generated Swagger API docs

## 🚀 Quick Start

### 1. Environment Requirements
- Go 1.21+
- Docker (optional)
- MySQL/PostgreSQL/SQLite (choose based on configuration)

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Configuration
Copy the configuration template:
```bash
cp settings.example.yaml settings.yaml
```

Then edit `settings.yaml` with your actual database credentials and settings.

### 4. Start the Project

#### Using SQLite (Recommended for Development)
```bash
# No database installation needed, run directly
go run main.go
```

#### Using MySQL
```bash
# Update settings.yaml with your MySQL credentials
go run main.go
```

#### Using PostgreSQL
```bash
# Update settings.yaml with your PostgreSQL credentials
go run main.go
```

### 5. Verify Startup
Visit these addresses to verify the project is running:
- API Documentation: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/health
- Default Admin: admin/admin123

## 🏗️ Project Structure

```
rbac_admin_server/
├── 📂 api/                    # RESTful API endpoints
├── 📂 config/                 # Configuration management system
│   ├── config.go             # Configuration struct definitions
│   └── loader.go             # Configuration loader
├── 📂 core/                   # Core startup logic
│   ├── initializer.go        # Project initializer
│   ├── audit/                # Audit logging system
│   └── errors/               # Error handling system
├── 📂 database/               # Database factory system
│   ├── database_factory.go   # Database factory core implementation
│   ├── migrator.go           # Database migration and initialization
│   └── models/               # Data model definitions
│       ├── user.go           # User, role, permission models
├── 📂 examples/               # Usage examples
├── 📂 global/                 # Global variable management
├── 📝 main.go                # Program entry point
├── ⚙️ settings.example.yaml  # Configuration template (safe)
├── 🐳 docker-compose.yml      # Docker deployment configuration
├── 🐳 Dockerfile             # Docker image build
└── 📦 go.mod                 # Dependency management
```

## 📊 Database Support

### Configuration Examples

#### MySQL
```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  username: "your_username"      # Change to your MySQL username
  password: "your_password"      # Change to your MySQL password
  database: "rbac_admin"
  charset: "utf8mb4"
```

#### PostgreSQL
```yaml
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  username: "your_username"        # Change to your PostgreSQL username
  password: "your_password"        # Change to your PostgreSQL password
  database: "rbac_admin"
```

#### SQLite (Development Recommended)
```yaml
database:
  type: "sqlite"
  path: "./rbac_admin.db"
```

## 🐳 Docker Deployment

### Using Docker Compose
```bash
# Start all services (MySQL + Redis + RBAC server)
docker-compose up -d

# View logs
docker-compose logs -f rbac_server

# Stop services
docker-compose down
```

### Build Image Separately
```bash
# Build image
docker build -t rbac-admin-server .

# Run container
docker run -p 8080:8080 rbac-admin-server
```

## 🔐 Default Accounts

Automatically created on first startup:
- **Username**: admin
- **Password**: admin123
- **Role**: Super Administrator
- **Permissions**: All permissions

## 📚 API Documentation

After starting the project, visit:
- Swagger Documentation: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/health
- Metrics: http://localhost:8080/metrics

## 🔧 Development Guide

### Adding New API Endpoints
1. Create new route files in the `api/` directory
2. Implement corresponding handler functions
3. Register routes in main.go
4. Update Swagger documentation

### Database Model Extension
1. Add new models in the `database/models/` directory
2. Update migration logic in `migrator.go`
3. Run the project for automatic migration

### Configuration Extension
1. Add new configuration structs in `config/config.go`
2. Update `settings.example.yaml` template
3. Add validation logic in `loader.go`

## 🔒 Security Notes

⚠️ **Important Security Reminders:**

1. **Never commit sensitive configuration files**:
   - `settings.yaml` (contains database passwords)
   - `.env` files
   - SSL certificates
   - API keys

2. **Production Security Checklist**:
   - Change all default passwords
   - Use strong JWT secrets
   - Enable HTTPS
   - Configure proper CORS
   - Set up rate limiting
   - Use environment variables for sensitive data

3. **Environment Variables**:
   ```bash
   export DB_PASSWORD="your_secure_password"
   export JWT_SECRET="your_strong_jwt_secret"
   export REDIS_PASSWORD="your_redis_password"
   ```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📞 Support

For support, please open an issue in the GitHub repository or contact the development team.