#!/bin/bash
set -e
echo "📦 Tidying Go module dependencies..."
go mod tidy
echo "✅ Dependencies ready."
