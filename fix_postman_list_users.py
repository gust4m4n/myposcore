#!/usr/bin/env python3
"""
Fix the invalid JSON in List Users response
"""

import json
import re

# Read file
with open('MyPOSCore.postman_collection.json', 'r') as f:
    content = f.read()

# Fix the typo: "role: staff | branchadmin" -> "role": "staff"
content = content.replace('"role: staff | branchadmin"', '"role": "staff"')

# Find and fix List Users Success response body
# Pattern to find the problematic response
pattern = r'("name": "List Users Success".*?"body": ")(\{[^}]*"page": 1.*?"data": \[[^\]]*\]\})'

def fix_list_users_response(match):
    prefix = match.group(1)
    old_body_escaped = match.group(2)
    
    # Unescape the JSON
    old_body = old_body_escaped.replace('\\n', '\n').replace('\\"', '"')
    
    try:
        # Parse the old body
        old_data = json.loads(old_body)
        
        # Create new standardized format
        items = old_data.get('data', [])
        new_body = {
            "code": 0,
            "message": "Users retrieved successfully",
            "data": {
                "items": items,
                "pagination": {
                    "page": old_data.get('page', 1),
                    "limit": old_data.get('page_size', 32),
                    "total": old_data.get('total_items', len(items)),
                    "total_pages": old_data.get('total_pages', 1)
                }
            }
        }
        
        # Convert back to escaped JSON string
        new_body_str = json.dumps(new_body, indent=2)
        new_body_escaped = new_body_str.replace('"', '\\"').replace('\n', '\\n')
        
        return prefix + new_body_escaped
    except:
        return match.group(0)

# Try to fix with regex
content = re.sub(pattern, fix_list_users_response, content, flags=re.DOTALL)

# Write back
with open('MyPOSCore.postman_collection.json', 'w') as f:
    f.write(content)

print("✓ Fixed List Users Success response")

# Verify it's valid JSON
try:
    with open('MyPOSCore.postman_collection.json', 'r') as f:
        collection = json.load(f)
    print("✓ Postman collection is valid JSON")
except json.JSONDecodeError as e:
    print(f"✗ JSON error: {e}")
