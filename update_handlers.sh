#!/bin/bash

# Script to update remaining handler files to use standardized response format
# Run this script from the myposcore directory

echo "Starting handler updates..."

# Function to add utils import and replace responses in a file
update_handler() {
    local file=$1
    echo "Updating $file..."
    
    # Add utils import if not present
    if ! grep -q '"myposcore/utils"' "$file"; then
        # Replace net/http import with utils
        sed -i '' 's|"net/http"|"myposcore/utils"|g' "$file"
    fi
    
    # Replace common response patterns
    sed -i '' 's|c\.JSON(http\.StatusOK, gin\.H{"message": "\([^"]*\)", "data": \([^}]*\)})|utils.Success(c, "\1", \2)|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusOK, gin\.H{"message": "\([^"]*\)"})|utils.SuccessWithoutData(c, "\1")|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusOK, gin\.H{"data": \([^}]*\)})|utils.Success(c, "Success", \1)|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusOK, \([^)]*\))|utils.Success(c, "Success", \1)|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusBadRequest, gin\.H{"error": \([^}]*\)})|utils.BadRequest(c, \1)|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusUnauthorized, gin\.H{"error": \([^}]*\)})|utils.Unauthorized(c, \1)|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusForbidden, gin\.H{"error": \([^}]*\)})|utils.Forbidden(c, \1)|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusNotFound, gin\.H{"error": \([^}]*\)})|utils.NotFound(c, \1)|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusInternalServerError, gin\.H{"error": \([^}]*\)})|utils.InternalError(c, \1)|g' "$file"
    sed -i '' 's|c\.JSON(http\.StatusConflict, gin\.H{"error": \([^}]*\)})|utils.Conflict(c, \1)|g' "$file"
}

# Update remaining handlers
cd /Users/gustaman/Desktop/GUSTAMAN7/myposcore/handlers

# Files that need complete updates
for file in category_handler.go order_handler.go payment_handler.go audit_trail_handler.go superadmin_handler.go tnc_handler.go faq_handler.go health_handler.go dev_handler.go; do
    if [ -f "$file" ]; then
        update_handler "$file"
    fi
done

# Note: product_handler.go needs manual attention due to complex multipart sections

echo "Handler updates complete!"
echo ""
echo "Note: You may need to manually review and adjust:"
echo "  - Multi-line response objects"
echo "  - Complex conditional responses"
echo "  - product_handler.go (partially updated)"
echo ""
echo "After running this script, run: go fmt ./... to format the code"
