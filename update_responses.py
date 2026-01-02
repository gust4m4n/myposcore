#!/usr/bin/env python3
import re
import glob
import os

def update_response(file_path):
    """Update response patterns in a Go handler file."""
    
    with open(file_path, 'r') as f:
        content = f.read()
    
    original = content
    
    # Simple single-line error responses
    content = re.sub(
        r'c\.JSON\(http\.StatusBadRequest,\s*gin\.H\{"error":\s*([^}]+)\}\)',
        r'utils.BadRequest(c, \1)',
        content
    )
    content = re.sub(
        r'c\.JSON\(http\.StatusUnauthorized,\s*gin\.H\{"error":\s*([^}]+)\}\)',
        r'utils.Unauthorized(c, \1)',
        content
    )
    content = re.sub(
        r'c\.JSON\(http\.StatusForbidden,\s*gin\.H\{"error":\s*([^}]+)\}\)',
        r'utils.Forbidden(c, \1)',
        content
    )
    content = re.sub(
        r'c\.JSON\(http\.StatusNotFound,\s*gin\.H\{"error":\s*([^}]+)\}\)',
        r'utils.NotFound(c, \1)',
        content
    )
    content = re.sub(
        r'c\.JSON\(http\.StatusInternalServerError,\s*gin\.H\{"error":\s*([^}]+)\}\)',
        r'utils.InternalError(c, \1)',
        content
    )
    content = re.sub(
        r'c\.JSON\(http\.StatusConflict,\s*gin\.H\{"error":\s*([^}]+)\}\)',
        r'utils.Conflict(c, \1)',
        content
    )
    
    # Simple single-line success with data
    content = re.sub(
        r'c\.JSON\(http\.StatusOK,\s*gin\.H\{"data":\s*([^}]+)\}\)',
        r'utils.Success(c, "Success", \1)',
        content
    )
    
    # Simple single-line success with message only
    content = re.sub(
        r'c\.JSON\(http\.StatusOK,\s*gin\.H\{"message":\s*([^}]+)\}\)',
        r'utils.SuccessWithoutData(c, \1)',
        content
    )
    
    if content != original:
        with open(file_path, 'w') as f:
            f.write(content)
        return True
    return False

def main():
    handlers_dir = "handlers"
    updated = []
    
    for file_path in glob.glob(os.path.join(handlers_dir, "*.go")):
        if "_test.go" in file_path:
            continue
            
        if update_response(file_path):
            updated.append(file_path)
            print(f"Updated: {file_path}")
    
    print(f"\nTotal files updated: {len(updated)}")

if __name__ == "__main__":
    main()
