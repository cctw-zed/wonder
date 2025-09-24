#!/bin/bash

# Test runner script for Wonder project
# Supports different types of tests with proper environment setup

set -e

# Source environment
source .envrc

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Help function
show_help() {
    cat << EOF
Usage: $0 [OPTIONS] [TEST_TYPE]

Test runner for Wonder project

TEST_TYPE:
    unit        Run unit tests only (default)
    integration Run integration tests only
    e2e         Run end-to-end tests only
    all         Run all tests
    coverage    Run tests with coverage report

OPTIONS:
    -h, --help      Show this help message
    -v, --verbose   Run tests in verbose mode
    -s, --short     Run tests in short mode (skip long-running tests)
    --race          Run tests with race detection
    --clean         Clean test cache before running

EXAMPLES:
    $0                          # Run unit tests
    $0 all                      # Run all tests
    $0 integration -v           # Run integration tests verbosely
    $0 coverage                 # Generate coverage report
    $0 e2e --clean             # Run E2E tests with clean cache

EOF
}

# Default values
TEST_TYPE="unit"
VERBOSE=""
SHORT=""
RACE=""
CLEAN=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        -s|--short)
            SHORT="-short"
            shift
            ;;
        --race)
            RACE="-race"
            shift
            ;;
        --clean)
            CLEAN="yes"
            shift
            ;;
        unit|integration|e2e|all|coverage)
            TEST_TYPE="$1"
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Clean test cache if requested
if [[ "$CLEAN" == "yes" ]]; then
    print_step "Cleaning test cache..."
    go clean -testcache
    print_success "Test cache cleaned"
fi

# Build flags
BUILD_FLAGS="$VERBOSE $SHORT $RACE"

# Function to run tests
run_tests() {
    local pattern="$1"
    local description="$2"

    print_step "Running $description..."

    if go test $BUILD_FLAGS $pattern; then
        print_success "$description passed"
        return 0
    else
        print_error "$description failed"
        return 1
    fi
}

# Function to setup test database for integration/e2e tests
setup_test_db() {
    print_step "Setting up test database environment..."
    export DB_USERNAME="test"
    export DB_PASSWORD="test"
    export DB_DATABASE="wonder_test"
    print_success "Test database environment configured"
}

# Main test execution
case $TEST_TYPE in
    unit)
        print_step "Running unit tests..."
        run_tests "./internal/..." "Unit tests"
        ;;

    integration)
        setup_test_db
        print_step "Running integration tests..."
        run_tests "./test/integration/..." "Integration tests"
        ;;

    e2e)
        setup_test_db
        print_step "Running end-to-end tests..."
        if [[ "$SHORT" == "-short" ]]; then
            print_warning "E2E tests are skipped in short mode"
            exit 0
        fi
        run_tests "./test/e2e/..." "End-to-end tests"
        ;;

    all)
        print_step "Running all tests..."

        # Unit tests
        if ! run_tests "./internal/..." "Unit tests"; then
            exit 1
        fi

        # Integration tests
        setup_test_db
        if ! run_tests "./test/integration/..." "Integration tests"; then
            exit 1
        fi

        # E2E tests (unless in short mode)
        if [[ "$SHORT" != "-short" ]]; then
            if ! run_tests "./test/e2e/..." "End-to-end tests"; then
                exit 1
            fi
        else
            print_warning "E2E tests skipped in short mode"
        fi

        print_success "All tests passed!"
        ;;

    coverage)
        print_step "Running tests with coverage..."
        setup_test_db

        # Run tests with coverage
        go test $VERBOSE $SHORT -coverprofile=coverage.out ./...

        # Generate coverage report
        print_step "Generating coverage report..."
        go tool cover -html=coverage.out -o coverage.html

        # Show coverage summary
        print_step "Coverage summary:"
        go tool cover -func=coverage.out | tail -1

        print_success "Coverage report generated: coverage.html"
        ;;

    *)
        print_error "Invalid test type: $TEST_TYPE"
        show_help
        exit 1
        ;;
esac