binary_name := "worklog"
docker_image := "your-registry/worklog"
version := env_var_or_default("VERSION", "latest")

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
    ./bin/worklog add -t 8 -d "Backend API development" -p "E-Commerce Platform" -c "ACME Corp" -n "Alice Johnson" -r 650
    ./bin/worklog add -t 6 -d "Frontend design improvements" -p "E-Commerce Platform" -c "ACME Corp" -n "Bob Smith" -r 600
    ./bin/worklog add -t 4.5 -d "Bug fixes and testing" -p "Mobile App" -c "TechStart AB" -n "Alice Johnson" -r 650
    ./bin/worklog add -t 7.5 -d "Database optimization" -p "Data Pipeline" -c "TechStart AB" -n "Charlie Davis" -r 750
    ./bin/worklog add -t 5 -d "UI/UX improvements" -p "Dashboard" -c "WebDev Inc" -n "Bob Smith" -r 600
    @echo "✓ Sample data added successfully"

# Add quick test entry using defaults
test-quick: build
    ./bin/worklog add -t 3 -d "Quick task"

# Generate monthly report
test-report: build
    @echo "Generating monthly report..."
    ./bin/worklog report

# Run full test: build, add sample data, generate report
test-full: build test-add test-report

# Show all available commands
list:
    @just --list
