#!/bin/bash
set -e
echo "📦 Initializing Go modules..."
go mod init github.com/user/erp-pos-hotel
go mod tidy
echo "✅ Dependencies ready."
