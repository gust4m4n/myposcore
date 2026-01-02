#!/usr/bin/env python3
"""
Script to update handler files to use standardized response utilities.
This handles complex multi-line replacements that sed can't handle.
"""

import re
import os

HANDLERS_DIR = "/Users/gustaman/Desktop/GUSTAMAN7/myposcore/handlers"

# Files to update (excluding already completed ones)
FILES_TO_UPDATE = [
    "product_handler.go",
    "category_handler.go",
    "order_handler.go",
    "payment_handler.go",
    "audit_trail_handler.go",
    "superadmin_handler.go",
    "tnc_handler.go",
    "faq_handler.go",
    "health_handler.go",
    "dev_handler.go",
]

def add_utils_import(content):
    """Add utils import if not present"""
    if '"myposcore/utils"' in content:
        return content
    
    # Replace "net/http" with "myposcore/utils"
    content = content.replace('"net/http"', '"myposcore/utils"')
    return content

def update_responses(content):
    """Update all response patterns"""
    
    # Pattern 1: c.JSON(http.StatusOK, gin.H{"message": "...", "data": ...})
    content = re.sub(
        r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\s*"message":\s*"([^"]+)",\s*"data":\s*([^}]+)\}\)',
        r'utils.Success(c, "\1", \2)',
        content
    )
    
    # Pattern 2: c.JSON(http.StatusOK, gin.H{"message": "..."})
    content = re.sub(
        r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\s*"message":\s*"([^"]+)"\s*\}\)',
        r'utils.SuccessWithoutData(c, "\1")',
        content
    )
    
    # Pattern 3: c.JSON(http.StatusOK, gin.H{"data": ...})
    content = re.sub(
        r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\s*"data":\s*([^}]+)\}\)',
        r'utils.Success(c, "Success", \1)',
        content
    )
    
    # Pattern 4: c.JSON(http.StatusBadRequest, gin.H{"error": ...})
    content = re.sub(
        r'c\.JSON\(http\.StatusBadRequest,\s*gin\.H\{\s*"error":\s*([^}]+)\}\)',
        r'utils.BadRequest(c, \1)',
        content
    )
    
    # Pattern 5: c.JSON(http.StatusUnauthorized, gin.H{"error": ...})
    content = re.sub(
        r'c\.JSON\(http\.StatusUnauthorized,\s*gin\.H\{\s*"error":\s*([^}]+)\}\)',
        r'utils.Unauthorized(c, \1)',
        content
    )
    
    # Pattern 6: c.JSON(http.StatusForbidden, gin.H{"error": ...})
    content = re.sub(
        r'c\.JSON\(http\.StatusForbidden,\s*gin\.H\{\s*"error":\s*([^}]+)\}\)',
        r'utils.Forbidden(c, \1)',
        content
    )
    
    # Pattern 7: c.JSON(http.StatusNotFound, gin.H{"error": ...})
    content = re.sub(
        r'c\.JSON\(http\.StatusNotFound,\s*gin\.H\{\s*"error":\s*([^}]+)\}\)',
        r'utils.NotFound(c, \1)',
        content
    )
    
    # Pattern 8: c.JSON(http.StatusInternalServerError, gin.H{"error": ...})
    content = re.sub(
        r'c\.JSON\(http\.StatusInternalServerError,\s*gin\.H\{\s*"error":\s*([^}]+)\}\)',
        r'utils.InternalError(c, \1)',
        content
    )
    
    # Pattern 9: c.JSON(http.StatusConflict, gin.H{"error": ...})
    content = re.sub(
        r'c\.JSON\(http\.StatusConflict,\s*gin\.H\{\s*"error":\s*([^}]+)\}\)',
        r'utils.Conflict(c, \1)',
        content
    )
    
    # Pattern 10: Simple StatusOK with data directly (no gin.H)
    content = re.sub(
        r'c\.JSON\(http\.StatusOK,\s*([a-zA-Z_][a-zA-Z0-9_]*)\)',
        r'utils.Success(c, "Success", \1)',
        content
    )
    
    return content

def process_file(filepath):
    """Process a single file"""
    try:
        with open(filepath, 'r') as f:
            content = f.read()
        
        original_content = content
        content = add_utils_import(content)
        content = update_responses(content)
        
        if content != original_content:
            with open(filepath, 'w') as f:
                f.write(content)
            return True
        return False
    except Exception as e:
        print(f"Error processing {filepath}: {e}")
        return False

def main():
    """Main function"""
    print("Starting handler file updates...")
    updated_count = 0
    
    for filename in FILES_TO_UPDATE:
        filepath = os.path.join(HANDLERS_DIR, filename)
        if os.path.exists(filepath):
            print(f"Processing {filename}...")
            if process_file(filepath):
                updated_count += 1
                print(f"  ✓ Updated {filename}")
            else:
                print(f"  - No changes needed for {filename}")
        else:
            print(f"  ✗ File not found: {filename}")
    
    print(f"\nCompleted! Updated {updated_count} files.")
    print("\nNext steps:")
    print("1. Run: go fmt ./...")
    print("2. Run: go build")
    print("3. Review changes with: git diff")

if __name__ == "__main__":
    main()
