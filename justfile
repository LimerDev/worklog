binary_name := "timetrack"
docker_image := "your-registry/timetrack"
version := env_var_or_default("VERSION", "latest")

# Bygg applikationen
build:
    go build -o {{binary_name}} main.go

# Kör applikationen
run *args:
    go run main.go {{args}}

# Rensa byggfiler
clean:
    go clean
    rm -f {{binary_name}}

# Bygg Docker-imagen
docker-build:
    docker build -t {{docker_image}}:{{version}} .

# Bygg och pusha Docker-imagen
docker-push: docker-build
    docker push {{docker_image}}:{{version}}

# Deploya till Kubernetes
k8s-deploy:
    kubectl apply -f k8s/namespace.yaml
    kubectl apply -f k8s/postgres.yaml
    kubectl apply -f k8s/timetrack.yaml

# Ta bort från Kubernetes
k8s-delete:
    kubectl delete -f k8s/timetrack.yaml
    kubectl delete -f k8s/postgres.yaml
    kubectl delete -f k8s/namespace.yaml

# Kör tester
test:
    go test -v ./...

# Hämta och organisera dependencies
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
config-set:
    ./timetrack config set -n "Alice Johnson" -c "ACME Corp" -p "E-Commerce Platform" -r 650

# View current configuration
config-show:
    ./timetrack config

# Clear configuration
config-clear:
    ./timetrack config clear

# Add sample data for testing
test-add:
    @echo "Adding sample time entries..."
    ./timetrack add -t 8 -d "Backend API development" -p "E-Commerce Platform" -c "ACME Corp" -n "Alice Johnson" -r 650
    ./timetrack add -t 6 -d "Frontend design improvements" -p "E-Commerce Platform" -c "ACME Corp" -n "Bob Smith" -r 600
    ./timetrack add -t 4.5 -d "Bug fixes and testing" -p "Mobile App" -c "TechStart AB" -n "Alice Johnson" -r 650
    ./timetrack add -t 7.5 -d "Database optimization" -p "Data Pipeline" -c "TechStart AB" -n "Charlie Davis" -r 750
    ./timetrack add -t 5 -d "UI/UX improvements" -p "Dashboard" -c "WebDev Inc" -n "Bob Smith" -r 600
    @echo "✓ Sample data added successfully"

# Add quick test entry using defaults
test-quick:
    ./timetrack add -t 3 -d "Quick task"

# Generate monthly report
test-report:
    @echo "Generating monthly report..."
    ./timetrack report

# Run full test: build, add sample data, generate report
test-full: build test-add test-report

# Visa alla tillgängliga kommandon
list:
    @just --list
