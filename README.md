# Worklog

A simple CLI application for time tracking and consultant billing, written in Go.

## Features

- Register work logs with hours, description, project, client, and consultant
- Flexible filtering and retrieval of work logs by consultant, project, customer, date, or date range
- Export work logs to CSV format with customizable filters
- Hourly rates stored per time entry for cost calculation with historical accuracy
- Calculate costs based on hourly rates and worked hours
- Normalized database structure: Client → Project → Time Entry
- PostgreSQL database for storage
- Kubernetes deployment support for cluster hosting

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

# Run with Docker (requires PostgreSQL and ~/.worklog/config.json)
docker run \
  -v ~/.worklog:/root/.worklog \
  worklog add --help

# Or use environment variables to override config
docker run \
  -e WORKLOG_DATABASE_HOST=db.example.com \
  -e WORKLOG_DATABASE_PASSWORD=yourpassword \
  -v ~/.worklog:/root/.worklog \
  worklog add --help
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
   kubectl apply -f k8s/postgres-storage.yaml
   kubectl apply -f k8s/postgres-config.yaml
   kubectl apply -f k8s/postgres-deployment.yaml
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

### Retrieve and filter work logs

Get all work logs:
```bash
worklog get
```

Get entries for a specific consultant:
```bash
worklog get -n "Alice Johnson"
```

Get entries for a specific project:
```bash
worklog get -p "E-Commerce Platform"
```

Get entries for a specific customer:
```bash
worklog get -c "ACME Corp"
```

Get entries for a specific date:
```bash
worklog get -D 2025-11-29
```

Get entries for a specific month:
```bash
worklog get -m 2025-11
```

Get entries for a date range:
```bash
worklog get --from 2025-11-01 --to 2025-11-30
```

Get today's entries:
```bash
worklog get --today
```

Combine multiple filters:
```bash
worklog get -n "Alice Johnson" -p "E-Commerce Platform"
worklog get --today -c "ACME Corp"
```

The output shows:
- Table with all matching work logs (date, consultant, hours, rate, cost, project, customer, description)
- Total hours and costs

### Export work logs to CSV

Export all entries:
```bash
worklog export -o report.csv
```

Export to stdout:
```bash
worklog export | head -20
```

Export with filters (same as get command):
```bash
worklog export -n "Alice Johnson" -o alice_report.csv
worklog export -m 2025-11 -o november_report.csv
worklog export --from 2025-11-01 --to 2025-11-30 -o period_report.csv
worklog export -c "ACME Corp" -o acme_report.csv
```

The CSV file includes:
- Headers: DATE, CONSULTANT, PROJECT, CUSTOMER, DESCRIPTION, HOURS, RATE, COST
- All matching work log entries
- Total row with summed hours and costs

### Using with Kubernetes

Run commands in the K8s pod:

```bash
# Add a time entry
kubectl exec -it -n worklog deployment/worklog -- ./worklog add \
  -t 8 -d "Development" -p "Project A" -c "Client AB" -n "Consultant" -r 650

# Retrieve work logs
kubectl exec -it -n worklog deployment/worklog -- ./worklog get -n "Consultant"
```

Alternatively, create a shell alias:
```bash
alias worklog='kubectl exec -it -n worklog deployment/worklog -- ./worklog'

# Now you can run:
worklog add -t 8 -d "Development" -p "Project A" -c "Client AB" -n "Consultant" -r 650
worklog get -n "Consultant"
```

## Configuration

### Configuration File

The application reads configuration from `~/.worklog/config.json`:

```json
{
  "default_consultant": "Your Name",
  "default_client": "Client Name",
  "default_project": "Project Name",
  "default_rate": 650,
  "database": {
    "host": "192.168.0.20",
    "port": "30432",
    "user": "wl_admin",
    "password": "your-secure-password",
    "name": "worklog"
  }
}
```

### Environment Variables

Database configuration can be overridden with environment variables (prefix: `WORKLOG_`):

- `WORKLOG_DATABASE_HOST` - Database host
- `WORKLOG_DATABASE_PORT` - Database port
- `WORKLOG_DATABASE_USER` - Database user
- `WORKLOG_DATABASE_PASSWORD` - Database password
- `WORKLOG_DATABASE_NAME` - Database name
- `WORKLOG_DEFAULT_CONSULTANT` - Default consultant name
- `WORKLOG_DEFAULT_CLIENT` - Default client name
- `WORKLOG_DEFAULT_PROJECT` - Default project name
- `WORKLOG_DEFAULT_RATE` - Default hourly rate

**Example with test database:**
```bash
WORKLOG_DATABASE_HOST=localhost \
WORKLOG_DATABASE_PORT=5432 \
WORKLOG_DATABASE_USER=testuser \
WORKLOG_DATABASE_PASSWORD=testpass \
WORKLOG_DATABASE_NAME=testdb \
worklog report
```

All test commands automatically use the test database configuration from the environment variables.

## Security

**IMPORTANT**: Change the default password in `k8s/postgres.yaml` before deploying to production!

```yaml
stringData:
  POSTGRES_PASSWORD: your-secure-password-here
```

## Tech Stack

- **Go** - Programming language
- **Cobra** - CLI framework
- **Viper** - Configuration management
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

# Test commands (automatically use test database)
just test-add              # Add sample test data
just test-quick            # Add a quick test entry
just test-get-all          # Get all work logs
just test-export-all       # Export all entries to CSV
just test-export-consultant # Export entries for Alice Johnson to CSV
just test-full             # Build + add sample data + run all get tests
```

## Development

### Local Setup

1. **Create config file:**
   ```bash
   mkdir -p ~/.worklog

   # Create config.json with your database details
   cat > ~/.worklog/config.json << 'EOF'
   {
     "default_consultant": "Your Name",
     "default_client": "Your Client",
     "default_project": "Your Project",
     "default_rate": 650,
     "database": {
       "host": "your-db-host",
       "port": "5432",
       "user": "your-user",
       "password": "your-password",
       "name": "worklog"
     }
   }
   EOF
   ```

2. **Run with Docker Compose (optional):**
   ```bash
   # Start PostgreSQL
   docker-compose up -d

   # Run the application with test database (all test commands use test DB automatically)
   just test-add      # Add sample data
   just test-get-all  # Get all work logs
   just test-full     # Build + add data + run all get tests

   # Stop the database
   docker-compose down
   ```

### Build and test

```bash
# Download dependencies
go mod tidy

# Build
just build

# Run tests (when implemented)
go test ./...

# Test with sample data (uses test database configuration)
just test-full              # Builds, adds sample data, and runs all get tests
just test-add               # Add more sample data
just test-quick             # Add a quick single test entry
just test-get-all           # Get and display all work logs
just test-export-all        # Export all work logs to CSV
just test-export-consultant # Export work logs for specific consultant to CSV
```

## License

MIT
