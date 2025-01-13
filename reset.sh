#!/bin/bash

echo "🗑️  Cleaning up Nogram data..."
docker-compose down -v

echo "🚀 Starting fresh Nogram instance..."
docker-compose up --build
