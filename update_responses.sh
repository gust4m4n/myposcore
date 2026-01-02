#!/bin/bash

# Script to update all handler responses to use new standardized format
# Usage: ./update_responses.sh

HANDLERS_DIR="handlers"

echo "Updating handler response formats..."

# Update all .go files in handlers directory (except test files)
for file in $HANDLERS_DIR/*.go; do
    # Skip test files
    if [[ $file == *"_test.go" ]]; then
        continue
    fi
    
    echo "Processing: $file"
    
    # Backup original file
    cp "$file" "$file.bak"
    
    # Replace common patterns - Simple error responses
    sed -i '' 's/c\.JSON(http\.StatusBadRequest, gin\.H{"error": \([^}]*\)})/utils.BadRequest(c, \1)/g' "$file"
    sed -i '' 's/c\.JSON(http\.StatusUnauthorized, gin\.H{"error": \([^}]*\)})/utils.Unauthorized(c, \1)/g' "$file"
    sed -i '' 's/c\.JSON(http\.StatusForbidden, gin\.H{"error": \([^}]*\)})/utils.Forbidden(c, \1)/g' "$file"
    sed -i '' 's/c\.JSON(http\.StatusNotFound, gin\.H{"error": \([^}]*\)})/utils.NotFound(c, \1)/g' "$file"
    sed -i '' 's/c\.JSON(http\.StatusInternalServerError, gin\.H{"error": \([^}]*\)})/utils.InternalError(c, \1)/g' "$file"
    sed -i '' 's/c\.JSON(http\.StatusConflict, gin\.H{"error": \([^}]*\)})/utils.Conflict(c, \1)/g' "$file"
    
    # Simple success with message only
    sed -i '' 's/c\.JSON(http\.StatusOK, gin\.H{"message": \([^}]*\)})/utils.SuccessWithoutData(c, \1)/g' "$file"
    
done

echo "Done! Backup files created with .bak extension"
echo "Please review changes and test the application"
