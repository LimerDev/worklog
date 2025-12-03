# Worklog

A simple CLI application for time tracking and consultant billing, written in Go.

## Features

- ✅ Register time entries with hours, description, project, client, and consultant
- ✅ Manage consultants with individual hourly rates
- ✅ Calculate costs based on hourly rates and worked hours
- ✅ Generate monthly reports with financial summaries
- ✅ Normalized database structure: Client → Project → Time Entry
- ✅ PostgreSQL database for storage
- ✅ Kubernetes deployment support for cluster hosting

## Installation

### Local (without Docker)

```bash
# Download dependencies
go mod download

# Build the application
go build -o worklog

# Run the application
./worklog --help
```

### Docker

```bash
# Build Docker image
docker build -t worklog:latest .

# Run with Docker (requires PostgreSQL)
docker run -e DB_HOST=localhost -e DB_PASSWORD=yourpassword worklog add --help
```

### Kubernetes

1. Update the image name in `k8s/worklog.yaml` to your registry
2. Build and push the image:
   ```bash
   docker build -t your-registry/worklog:latest .
   docker push your-registry/worklog:latest
   ```
3. Apply Kubernetes manifests:
   ```bash
   kubectl apply -f k8s/namespace.yaml
   kubectl apply -f k8s/postgres.yaml
   kubectl apply -f k8s/worklog.yaml
   ```
4. Wait for pods to be ready:
   ```bash
   kubectl wait --for=condition=ready pod -l app=postgres -n worklog --timeout=300s
   kubectl wait --for=condition=ready pod -l app=worklog -n worklog --timeout=300s
   ```

## Usage

### Configure default values (first time setup)

Set your default consultant, client, project, and hourly rate to avoid entering them repeatedly:

```bash
worklog config set -n "Alice Johnson" -c "ACME Corp" -p "E-Commerce Platform" -r 650
```

View your configuration:
```bash
worklog config
```

Clear your configuration:
```bash
worklog config clear
```

### Add a time entry

**Simple (with configured defaults):**
```bash
worklog add -t 5 -d "Code review and meeting"
```

**With all parameters:**
```bash
worklog add \
  --hours 8 \
  --description "Feature development" \
  --project "Project A" \
  --client "Client AB" \
  --consultant "Alice Johnson" \
  --rate 650
```

**With specific date:**
```bash
worklog add \
  --hours 4.5 \
  --description "Bug fixes and testing" \
  --date 2025-11-29
```

**Short syntax:**
```bash
worklog add -t 8 -d "Development" -p "Project A" -c "Client AB" -n "Alice Johnson" -r 650
```

Configuration file location: `~/.worklog/config.json`

### Generate monthly report

Current month:
```bash
worklog report
```

Specific month:
```bash
worklog report --month 2025-11
```

The report shows:
- All time entries with consultant, hours, and calculated costs
- Total hours and total costs
- Summary per consultant (hours and costs)
- Summary per project (hours and costs)
- Summary per client (hours and costs)

### Using with Kubernetes

Run commands in the K8s pod:

```bash
# Add a time entry
kubectl exec -it -n worklog deployment/worklog -- ./worklog add \
  -t 8 -d "Development" -p "Project A" -c "Client AB" -n "Consultant" -r 650

# Generate report
kubectl exec -it -n worklog deployment/worklog -- ./worklog report
```

Alternatively, create a shell alias:
```bash
alias worklog='kubectl exec -it -n worklog deployment/worklog -- ./worklog'

# Now you can run:
worklog add -t 8 -d "Development" -p "Project A" -c "Client AB" -n "Consultant" -r 650
worklog report
```

## Environment Variables

The application uses the following environment variables for database configuration:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: worklog)
- `DB_PASSWORD` - Database password (default: worklog)
- `DB_NAME` - Database name (default: worklog)

## Security

**IMPORTANT**: Change the default password in `k8s/postgres.yaml` before deploying to production!

```yaml
stringData:
  POSTGRES_PASSWORD: your-secure-password-here
```

## Tech Stack

- **Go** - Programming language
- **Cobra** - CLI framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Database
- **Docker** - Containerization
- **Kubernetes** - Orchestration
- **just** - Command runner (alternative to Make)

## Quick Commands with just

The project uses [just](https://github.com/casey/just) to simplify common tasks:

```bash
# Show all available commands
just list

# Build the application
just build

# Run the application with arguments
just run add -t 8 -d "Test" -p "Project" -c "Client" -n "Consultant" -r 650

# Build and push Docker image
just docker-push

# Deploy to Kubernetes
just k8s-deploy

# Download dependencies
just deps

# Database commands
just db-start    # Start PostgreSQL database
just db-stop     # Stop the database
just db-reset    # Reset database (delete all data)
just db-logs     # Show database logs

# Test commands
just test-add    # Add sample test data
just test-report # Generate test report
just test-full   # Build + add sample data + generate report
```

## Development

### Run locally with PostgreSQL in Docker

#### With Docker Compose (recommended)

```bash
# Start PostgreSQL database
docker-compose up -d

# Wait for database to be ready
docker-compose ps

# Copy .env.example to .env if not already done
cp .env.example .env

# Run the application
go run main.go add -t 8 -d "Test" -p "Project" -c "Client" -n "Consultant" -r 650
go run main.go report

# Stop the database when done
docker-compose down

# Remove data as well (permanent deletion)
docker-compose down -v
```

#### With docker run

```bash
# Start PostgreSQL
docker run -d \
  --name worklog-db \
  -e POSTGRES_USER=worklog \
  -e POSTGRES_PASSWORD=worklog \
  -e POSTGRES_DB=worklog \
  -p 5432:5432 \
  postgres:16-alpine

# Run the application
go run main.go add -t 8 -d "Test" -p "Project" -c "Client" -n "Consultant" -r 650
go run main.go report
```

### Build and test

```bash
# Download dependencies
go mod tidy

# Run tests (when implemented)
go test ./...

# Build
go build -o worklog
```

## License

MIT
