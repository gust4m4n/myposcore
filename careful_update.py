#!/usr/bin/env python3
"""
Careful script to update handler files to use standardized response utilities.
Preserves net/http imports needed for constants.
"""

import re
import os

HANDLERS_DIR = "/Users/gustaman/Desktop/GUSTAMAN7/myposcore/handlers"

FILES_TO_UPDATE = [
    "logout_handler.go",
    "profile_handler.go",
    "change_password_handler.go",
    "admin_change_password_handler.go",
    "pin_handler.go",
    "admin_change_pin_handler.go",
    "user_handler.go",
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
    """Add utils import if not present, keep net/http"""
    if '"myposcore/utils"' in content:
        return content, False
    
    # Find the import block and add utils
    if 'import (' in content:
        # Add utils import in the import block
        pattern = r'(import \([^)]+)(\"github\.com/gin-gonic/gin\"\))'
        replacement = r'\1"myposcore/utils"\n\n\t\2'
        if re.search(pattern, content):
            content = re.sub(pattern, replacement, content)
            return content, True
    
    return content, False

def update_responses(content):
    """Update all response patterns carefully"""
    changes = 0
    
    # Pattern 1: Message + Data in single line
    pattern = r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\"message\":\s*\"([^\"]+)\",\s*\"data\":\s*([^\}]+)\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.Success(c, "\1", \2)', content)
        changes += 1
    
    # Pattern 2: Message only
    pattern = r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\"message\":\s*\"([^\"]+)\"\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.SuccessWithoutData(c, "\1")', content)
        changes += 1
    
    # Pattern 3: Data only
    pattern = r'c\.JSON\(http\.StatusOK,\s*gin\.H\{\"data\":\s*([^\}]+)\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.Success(c, "Success", \1)', content)
        changes += 1
    
    # Pattern 4: BadRequest
    pattern = r'c\.JSON\(http\.StatusBadRequest,\s*gin\.H\{\"error\":\s*([^\}]+)\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.BadRequest(c, \1)', content)
        changes += 1
    
    # Pattern 5: Unauthorized
    pattern = r'c\.JSON\(http\.StatusUnauthorized,\s*gin\.H\{\"error\":\s*([^\}]+)\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.Unauthorized(c, \1)', content)
        changes += 1
    
    # Pattern 6: Forbidden
    pattern = r'c\.JSON\(http\.StatusForbidden,\s*gin\.H\{\"error\":\s*([^\}]+)\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.Forbidden(c, \1)', content)
        changes += 1
    
    # Pattern 7: NotFound
    pattern = r'c\.JSON\(http\.StatusNotFound,\s*gin\.H\{\"error\":\s*([^\}]+)\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.NotFound(c, \1)', content)
        changes += 1
    
    # Pattern 8: InternalServerError
    pattern = r'c\.JSON\(http\.StatusInternalServerError,\s*gin\.H\{\"error\":\s*([^\}]+)\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.InternalError(c, \1)', content)
        changes += 1
    
    # Pattern 9: Conflict
    pattern = r'c\.JSON\(http\.StatusConflict,\s*gin\.H\{\"error\":\s*([^\}]+)\}\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.Conflict(c, \1)', content)
        changes += 1
    
    # Pattern 10: Simple StatusOK with variable (no gin.H)
    pattern = r'c\.JSON\(http\.StatusOK,\s*([a-zA-Z_][a-zA-Z0-9_]*)\)'
    if re.search(pattern, content):
        content = re.sub(pattern, r'utils.Success(c, "Success", \1)', content)
        changes += 1
    
    return content, changes

def process_file(filepath):
    """Process a single file"""
    try:
        with open(filepath, 'r') as f:
            content = f.read()
        
        original_content = content
        content, import_added = add_utils_import(content)
        content, response_changes = update_responses(content)
        
        if content != original_content:
            with open(filepath, 'w') as f:
                f.write(content)
            return True, response_changes
        return False, 0
    except Exception as e:
        print(f"Error processing {filepath}: {e}")
        return False, 0

def main():
    """Main function"""
    print("Starting handler file updates...")
    updated_count = 0
    total_changes = 0
    
    for filename in FILES_TO_UPDATE:
        filepath = os.path.join(HANDLERS_DIR, filename)
        if os.path.exists(filepath):
            print(f"Processing {filename}...")
            updated, changes = process_file(filepath)
            if updated:
                updated_count += 1
                total_changes += changes
                print(f"  ✓ Updated {filename} ({changes} response statements)")
            else:
                print(f"  - No changes needed for {filename}")
        else:
            print(f"  ✗ File not found: {filename}")
    
    print(f"\nCompleted!")
    print(f"  Files updated: {updated_count}")
    print(f"  Total response statements changed: {total_changes}")
    print("\nNext steps:")
    print("1. Run: go fmt ./...")
    print("2. Run: go build")
    print("3. Review changes with: git diff")

if __name__ == "__main__":
    main()
