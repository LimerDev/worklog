binary_name := "worklog"
docker_image := "your-registry/worklog"
version := env_var_or_default("VERSION", "latest")
test_env := "WORKLOG_DATABASE_HOST=localhost WORKLOG_DATABASE_PORT=5432 WORKLOG_DATABASE_USER=worklog WORKLOG_DATABASE_PASSWORD=worklog WORKLOG_DATABASE_NAME=worklog"

# Build application
build:
    mkdir -p bin
    go build -o bin/{{binary_name}} main.go

# Build with optimizations (release)
build-release:
    mkdir -p bin
    go build -ldflags="-s -w" -o bin/{{binary_name}} main.go
    upx --best --lzma bin/{{binary_name}} || true

# Install to ~/.local/bin
install: build-release
    mkdir -p ~/.local/bin
    cp bin/{{binary_name}} ~/.local/bin/
    chmod +x ~/.local/bin/{{binary_name}}
    @echo "Installed {{binary_name}} to ~/.local/bin"

# Uninstall from ~/.local/bin
uninstall:
    rm -f ~/.local/bin/{{binary_name}}
    @echo "Uninstalled {{binary_name}} from ~/.local/bin"

# Run application
run *args:
    go run main.go {{args}}

# Clean build files
clean:
    go clean
    rm -rf bin/

# Build Docker image
docker-build:
    docker build -t {{docker_image}}:{{version}} .

# Build and push Docker image
docker-push: docker-build
    docker push {{docker_image}}:{{version}}

# Deploy to Kubernetes
k8s-deploy:
    kubectl apply -f k8s/namespace.yaml
    kubectl apply -f k8s/postgres.yaml
    kubectl apply -f k8s/worklog.yaml

# Remove from Kubernetes
k8s-delete:
    kubectl delete -f k8s/worklog.yaml
    kubectl delete -f k8s/postgres.yaml
    kubectl delete -f k8s/namespace.yaml

# Run tests
test:
    go test -v ./...

# Download and organize dependencies
deps:
    go mod download
    go mod tidy

# Start local PostgreSQL database
db-start:
    docker-compose up -d
    @echo "Waiting for database to be ready..."
    @sleep 3
    docker-compose ps

# Stop local PostgreSQL database
db-stop:
    docker-compose down

# Stop and delete all database data
db-reset:
    docker-compose down -v
    docker-compose up -d
    @echo "Waiting for database to be ready..."
    @sleep 3

# Show database logs
db-logs:
    docker-compose logs -f postgres

# Configure default values for quick time tracking
config-set: build
    ./bin/worklog config set -n "Alice Johnson" -c "ACME Corp" -p "E-Commerce Platform" -r 650

# View current configuration
config-show: build
    ./bin/worklog config

# Clear configuration
config-clear: build
    ./bin/worklog config clear

# Add sample data for testing
test-add: build
    @echo "Adding sample time entries..."
    sh -c "{{test_env}} ./bin/worklog add -t 8 -d 'Backend API development' -p 'E-Commerce Platform' -c 'ACME Corp' -n 'Alice Johnson' -r 650"
    sh -c "{{test_env}} ./bin/worklog add -t 6 -d 'Frontend design improvements' -p 'E-Commerce Platform' -c 'ACME Corp' -n 'Bob Smith' -r 600"
    sh -c "{{test_env}} ./bin/worklog add -t 4.5 -d 'Bug fixes and testing' -p 'Mobile App' -c 'TechStart AB' -n 'Alice Johnson' -r 650"
    sh -c "{{test_env}} ./bin/worklog add -t 7.5 -d 'Database optimization' -p 'Data Pipeline' -c 'TechStart AB' -n 'Charlie Davis' -r 750"
    sh -c "{{test_env}} ./bin/worklog add -t 5 -d 'UI/UX improvements' -p 'Dashboard' -c 'WebDev Inc' -n 'Bob Smith' -r 600"
    @echo "✓ Sample data added successfully"

# Add quick test entry using defaults
test-quick: build
    sh -c "{{test_env}} ./bin/worklog add -t 3 -d 'Quick task'"

# Get all entries
test-get-all: build
    @echo "Getting all time entries..."
    sh -c "{{test_env}} ./bin/worklog get"

# Get entries for specific consultant
test-get-consultant: build
    @echo "Getting entries for Alice Johnson..."
    sh -c "{{test_env}} ./bin/worklog get -n 'Alice Johnson'"

# Get entries for specific customer
test-get-customer: build
    @echo "Getting entries for ACME Corp..."
    sh -c "{{test_env}} ./bin/worklog get -c 'ACME Corp'"

# Get entries for specific date
test-get-date: build
    @echo "Getting entries for 2025-12-05..."
    sh -c "{{test_env}} ./bin/worklog get -D 2025-12-05"

# Get entries with combined filters
test-get-combined: build
    @echo "Getting entries for Alice Johnson in E-Commerce Platform project..."
    sh -c "{{test_env}} ./bin/worklog get -n 'Alice Johnson' -p 'E-Commerce Platform'"

# Get today's entries
test-get-today: build
    @echo "Getting today's entries..."
    sh -c "{{test_env}} ./bin/worklog get --today"

# Get today's entries for specific consultant
test-get-today-consultant: build
    @echo "Getting today's entries for Alice Johnson..."
    sh -c "{{test_env}} ./bin/worklog get --today -n 'Alice Johnson'"

# Run full test: build, add sample data and test get commands
test-full: build test-add test-get-all test-get-consultant test-get-customer test-get-date test-get-combined test-get-today test-get-today-consultant

# Show all available commands
list:
    @just --list
