#!/bin/bash

# Wait for Elasticsearch to be ready
echo "Waiting for Elasticsearch to be ready..."
until curl -s "http://localhost:9200/_cluster/health" > /dev/null; do
    sleep 3
done

# Delete existing indices if they exist
echo "Deleting existing indices..."
curl -X DELETE "http://localhost:9200/authors" 2>/dev/null
curl -X DELETE "http://localhost:9200/news" 2>/dev/null

# Create indices with mappings
echo "Creating indices..."

# Authors index
curl -X PUT "http://localhost:9200/authors" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": { "type": "text" },
      "bio": { "type": "text" },
      "imageUrl": { "type": "keyword" },
      "createdAt": { "type": "date" },
      "updatedAt": { "type": "date" }
    }
  }
}'

# News index
curl -X PUT "http://localhost:9200/news" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "title": { "type": "text" },
      "content": { "type": "text" },
      "authorId": { "type": "keyword" },
      "tags": { "type": "keyword" },
      "imageUrl": { "type": "keyword" },
      "createdAt": { "type": "date" },
      "updatedAt": { "type": "date" }
    }
  }
}'

echo "Loading data..."

# Load authors data
jq -c '.[]' authors.json | while read -r line; do
    id=$(echo $line | jq -r '.id')
    # Add timestamps if not present
    line=$(echo $line | jq '. + {createdAt: (now | todate), updatedAt: (now | todate)}')
    curl -X POST "http://localhost:9200/authors/_doc/$id" -H "Content-Type: application/json" -d "$line"
done

# Load news data
jq -c '.[]' news.json | while read -r line; do
    id=$(echo $line | jq -r '.id')
    # Add timestamps if not present
    line=$(echo $line | jq '. + {createdAt: (now | todate), updatedAt: (now | todate)}')
    curl -X POST "http://localhost:9200/news/_doc/$id" -H "Content-Type: application/json" -d "$line"
done

echo "Verifying data..."

# Verify the data was loaded
echo "Document counts:"
echo "Authors: $(curl -s -X GET "http://localhost:9200/authors/_count" | jq '.count')"
echo "News: $(curl -s -X GET "http://localhost:9200/news/_count" | jq '.count')"

echo "Seeding completed!" 