#!/usr/bin/env python3
"""
Simplify TNC endpoints in Postman collection - only keep GET /tnc
"""

import json

# Read collection
with open('MyPOSCore.postman_collection.json', 'r') as f:
    collection = json.load(f)

# Find TNC section and replace with simplified version
def simplify_tnc(item):
    if item.get('name') == 'Terms & Conditions':
        # Replace entire TNC section with just one GET endpoint
        item['item'] = [{
            "name": "Get Terms & Conditions",
            "request": {
                "method": "GET",
                "header": [],
                "url": {
                    "raw": "{{base_url}}/api/v1/tnc",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "tnc"]
                },
                "description": "Public endpoint untuk mendapatkan Terms & Conditions. Mengembalikan static content. Tidak memerlukan authentication."
            },
            "response": [{
                "name": "Get TnC Success",
                "originalRequest": {
                    "method": "GET",
                    "header": [],
                    "url": {
                        "raw": "{{base_url}}/api/v1/tnc",
                        "host": ["{{base_url}}"],
                        "path": ["api", "v1", "tnc"]
                    }
                },
                "status": "OK",
                "code": 200,
                "_postman_previewlanguage": "json",
                "header": [{"key": "Content-Type", "value": "application/json"}],
                "cookie": [],
                "body": json.dumps({
                    "code": 0,
                    "message": "Terms and conditions retrieved successfully",
                    "data": {
                        "title": "Terms and Conditions - MyPOS Core System",
                        "content": "# Terms and Conditions - MyPOS Core System\n\n## 1. Introduction\nWelcome to MyPOS Core...\n\n## 2. Acceptance of Terms\nBy using MyPOS Core, you agree to these terms...",
                        "version": "1.0.0",
                        "updated_at": "2026-01-03T00:00:00+07:00"
                    }
                }, indent=2)
            }]
        }]
        print(f"✓ Simplified Terms & Conditions section")
        return True
    
    # Process nested items
    if 'item' in item:
        for subitem in item['item']:
            if simplify_tnc(subitem):
                return True
    return False

# Process collection
if 'item' in collection:
    for item in collection['item']:
        simplify_tnc(item)

# Write back
with open('MyPOSCore.postman_collection.json', 'w') as f:
    json.dump(collection, f, indent='\t')

print("=" * 60)
print("✓ Postman collection updated")
print("✓ TNC now has only 1 endpoint: GET /api/v1/tnc")
