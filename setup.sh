#!/bin/bash
set -e
echo "📦 Initializing Go modules..."
go mod init github.com/allyelvis/erp-pos-hotel-system-go
go mod tidy
echo "✅ Dependencies ready."
