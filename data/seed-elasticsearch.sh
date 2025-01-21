#!/bin/bash

# Wait for Elasticsearch to be ready
echo "Waiting for Elasticsearch to be ready..."
until curl -s "http://localhost:9200/_cluster/health" > /dev/null; do
    sleep 3
done

# Delete existing indices if they exist
echo "Deleting existing indices..."
curl -X DELETE "http://localhost:9200/authors" 2>/dev/null
curl -X DELETE "http://localhost:9200/titles" 2>/dev/null
curl -X DELETE "http://localhost:9200/contents" 2>/dev/null

# Create indices with mappings
echo "Creating indices..."

# Authors index
curl -X PUT "http://localhost:9200/authors" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": { "type": "text" },
      "bio": { "type": "text" },
      "image_url": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}'

# Titles index
curl -X PUT "http://localhost:9200/titles" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "title": { "type": "text" },
      "author_id": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}'

# Contents index
curl -X PUT "http://localhost:9200/contents" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "content": { "type": "text" },
      "author_id": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" },
      "tags": { "type": "keyword" },
      "image_url": { "type": "keyword" }
    }
  }
}'

echo "Loading data..."

# Load authors data
jq -c '.[]' authors.json | while read -r line; do
    id=$(echo $line | jq -r '.id')
    curl -X POST "http://localhost:9200/authors/_doc/$id" -H "Content-Type: application/json" -d "$line"
done

# Load titles data
jq -c '.[]' titles.json | while read -r line; do
    id=$(echo $line | jq -r '.id')
    curl -X POST "http://localhost:9200/titles/_doc/$id" -H "Content-Type: application/json" -d "$line"
done

# Load contents data
jq -c '.[]' contents.json | while read -r line; do
    id=$(echo $line | jq -r '.id')
    curl -X POST "http://localhost:9200/contents/_doc/$id" -H "Content-Type: application/json" -d "$line"
done

echo "Verifying data..."

# Verify the data was loaded
echo "Document counts:"
echo "Authors: $(curl -s -X GET "http://localhost:9200/authors/_count" | jq '.count')"
echo "Titles: $(curl -s -X GET "http://localhost:9200/titles/_count" | jq '.count')"
echo "Contents: $(curl -s -X GET "http://localhost:9200/contents/_count" | jq '.count')"

echo "Seeding completed!" 