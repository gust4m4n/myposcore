#!/usr/bin/env python3
"""
Update complex multi-line gin.H{} responses to use utils response functions.
"""

import re
import os

def update_file(filepath):
    """Update a single file to convert gin.H responses to utils functions."""
    with open(filepath, 'r') as f:
        content = f.read()
    
    original_content = content
    
    # Pattern 1: Success responses with data (multi-line)
    # c.JSON(http.StatusOK, gin.H{
    #     "message": "...",
    #     "data": ...,
    # })
    pattern1 = r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\s*\n\s*"message":\s*"([^"]+)",\s*\n\s*"data":\s*([^,]+),\s*\n\s*\}\)'
    content = re.sub(pattern1, r'utils.Success(c, "\1", \2)', content, flags=re.MULTILINE)
    
    # Pattern 2: Success responses with data (single line compact)
    # c.JSON(http.StatusOK, gin.H{"message": "...", "data": ...})
    pattern2 = r'c\.JSON\(http\.StatusOK,\s*gin\.H\{"message":\s*"([^"]+)",\s*"data":\s*([^}]+)\}\)'
    content = re.sub(pattern2, r'utils.Success(c, "\1", \2)', content)
    
    # Pattern 3: Success responses without data
    # c.JSON(http.StatusOK, gin.H{"message": "..."})
    pattern3 = r'c\.JSON\(http\.StatusOK,\s*gin\.H\{"message":\s*"([^"]+)"\}\)'
    content = re.sub(pattern3, r'utils.SuccessWithoutData(c, "\1")', content)
    
    # Pattern 4: BadRequest responses  
    # c.JSON(http.StatusBadRequest, gin.H{"error": "..."})
    pattern4 = r'c\.JSON\(http\.StatusBadRequest,\s*gin\.H\{"error":\s*"([^"]+)"\}\)'
    content = re.sub(pattern4, r'utils.BadRequest(c, "\1")', content)
    
    # Pattern 5: Unauthorized responses
    # c.JSON(http.StatusUnauthorized, gin.H{"error": "..."})
    pattern5 = r'c\.JSON\(http\.StatusUnauthorized,\s*gin\.H\{"error":\s*"([^"]+)"\}\)'
    content = re.sub(pattern5, r'utils.Unauthorized(c, "\1")', content)
    
    # Pattern 6: NotFound responses
    # c.JSON(http.StatusNotFound, gin.H{"error": "..."})
    pattern6 = r'c\.JSON\(http\.StatusNotFound,\s*gin\.H\{"error":\s*"([^"]+)"\}\)'
    content = re.sub(pattern6, r'utils.NotFound(c, "\1")', content)
    
    # Pattern 7: InternalServerError responses
    # c.JSON(http.StatusInternalServerError, gin.H{"error": "..."})
    pattern7 = r'c\.JSON\(http\.StatusInternalServerError,\s*gin\.H\{"error":\s*"([^"]+)"\}\)'
    content = re.sub(pattern7, r'utils.InternalError(c, "\1")', content)
    
    if content != original_content:
        with open(filepath, 'w') as f:
            f.write(content)
        return True
    return False

# List of handler files to update
handler_files = [
    "handlers/audit_trail_handler.go",
    "handlers/category_handler.go",
    "handlers/change_password_handler.go",
    "handlers/dev_handler.go",
    "handlers/order_handler.go",
    "handlers/payment_handler.go",
    "handlers/pin_handler.go",
    "handlers/product_handler.go",
    "handlers/profile_handler.go",
    "handlers/superadmin_handler.go",
    "handlers/user_handler.go",
]

updated_count = 0
for filepath in handler_files:
    if os.path.exists(filepath):
        if update_file(filepath):
            print(f"✓ Updated: {filepath}")
            updated_count += 1
        else:
            print(f"- No changes: {filepath}")
    else:
        print(f"✗ Not found: {filepath}")

print(f"\nTotal files updated: {updated_count}")
