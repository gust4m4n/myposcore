#!/usr/bin/env python3
"""
Update ALL remaining gin.H{} responses to use utils response functions.
Handles complex multi-line and nested structures.
"""

import re
import os

def process_handler_file(filepath):
    """Process a handler file to update all gin.H responses."""
    with open(filepath, 'r') as f:
        content = f.read()
    
    original = content
    
    # Pattern: c.JSON(http.StatusOK, gin.H{ with message and data on multiple lines
    # Replace with utils.Success
    content = re.sub(
        r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\s*\n\s*"message":\s*"([^"]+)",\s*\n\s*"data":\s*([^}]+(?:\{[^}]*\})?[^}]*),\s*\n\s*\}\)',
        r'utils.Success(c, "\1", \2)',
        content,
        flags=re.MULTILINE | re.DOTALL
    )
    
    # Pattern: c.JSON(http.StatusOK, gin.H{ "data": ... }) without message
    # Replace with utils.Success with generic message
    def replace_data_only(match):
        data_content = match.group(1).strip()
        # Extract meaningful message from context if possible
        return f'utils.Success(c, "Operation successful", {data_content})'
    
    content = re.sub(
        r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\s*\n\s*"data":\s*([^}]+(?:\{[^}]*\})?[^}]*),\s*\n\s*\}\)',
        replace_data_only,
        content,
        flags=re.MULTILINE | re.DOTALL
    )
    
    # Pattern: c.JSON(http.StatusBadRequest, gin.H{"error": "..."})
    content = re.sub(
        r'c\.JSON\(http\.StatusBadRequest,\s*gin\.H\{\s*"error":\s*"([^"]+)"\s*\}\)',
        r'utils.BadRequest(c, "\1")',
        content
    )
    
    # Pattern: c.JSON(http.StatusInternalServerError, gin.H{"error": "..."})
    content = re.sub(
        r'c\.JSON\(http\.StatusInternalServerError,\s*gin\.H\{\s*"error":\s*"([^"]+)"\s*\}\)',
        r'utils.InternalError(c, "\1")',
        content
    )
    
    return content if content != original else None

# Process all handler files
handlers_dir = "handlers"
handler_files = [
    "base_handler.go",
    "change_password_handler.go",
    "order_handler.go",
    "payment_handler.go",
    "pin_handler.go",
    "product_handler.go",
    "superadmin_handler.go",
    "user_handler.go",
]

updated = []
for filename in handler_files:
    filepath = os.path.join(handlers_dir, filename)
    if os.path.exists(filepath):
        result = process_handler_file(filepath)
        if result:
            with open(filepath, 'w') as f:
                f.write(result)
            updated.append(filename)
            print(f"✓ Updated: {filename}")
        else:
            print(f"- No changes: {filename}")
    else:
        print(f"✗ Not found: {filename}")

print(f"\n✓ Updated {len(updated)} files")
print("Note: Some complex nested structures may need manual review")
