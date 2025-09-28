#!/bin/bash

# Test Coverage Script for Wonder Project
# This script runs all tests and generates coverage reports

set -e

echo "🧪 Running tests with coverage..."

# Source project environment
source .envrc

# Create coverage directory if it doesn't exist
mkdir -p coverage

# Run tests for main modules with coverage
echo "Running domain layer tests..."
go test -coverprofile=coverage/domain.out ./internal/domain/user

echo "Running application layer tests..."
go test -coverprofile=coverage/application.out ./internal/application/service

echo "Running interface layer tests..."
go test -coverprofile=coverage/interface.out ./internal/interfaces/http

echo "Running test utilities tests..."
go test -coverprofile=coverage/testutil.out ./internal/testutil/builder

# Combine coverage profiles
echo "Combining coverage profiles..."
echo "mode: set" > coverage/coverage.out
grep -h -v "^mode:" coverage/domain.out coverage/application.out coverage/interface.out coverage/testutil.out >> coverage/coverage.out

# Generate coverage report
echo "Generating coverage report..."
go tool cover -func=coverage/coverage.out

# Generate HTML report
echo "Generating HTML coverage report..."
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Check coverage threshold
COVERAGE=$(go tool cover -func=coverage/coverage.out | grep total: | awk '{print $3}' | sed 's/%//')
THRESHOLD=80

echo "📊 Coverage Summary:"
echo "   Total Coverage: ${COVERAGE}%"
echo "   Target Threshold: ${THRESHOLD}%"

if (( $(echo "$COVERAGE >= $THRESHOLD" | bc -l) )); then
    echo "✅ Coverage target met! (${COVERAGE}% >= ${THRESHOLD}%)"
    exit 0
else
    echo "❌ Coverage target not met! (${COVERAGE}% < ${THRESHOLD}%)"
    exit 1
fi