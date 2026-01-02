#!/usr/bin/env python3
"""
Fix all pagination responses to use standardized format
"""

import json

def fix_pagination_response(item):
    """Fix pagination responses recursively"""
    fixed_count = 0
    
    # Fix responses in this item
    if 'response' in item:
        for response in item['response']:
            if 'body' in response:
                try:
                    body = json.loads(response['body'])
                    
                    # Check if this is an old pagination format
                    if isinstance(body, dict) and 'page' in body and 'data' in body and 'code' not in body:
                        # This is old pagination format, convert it
                        items = body.get('data', [])
                        new_body = {
                            "code": 0,
                            "message": response.get('name', 'Operation successful'),
                            "data": {
                                "items": items,
                                "pagination": {
                                    "page": body.get('page', 1),
                                    "limit": body.get('page_size', body.get('limit', 10)),
                                    "total": body.get('total_items', body.get('total', len(items))),
                                    "total_pages": body.get('total_pages', 1)
                                }
                            }
                        }
                        response['body'] = json.dumps(new_body, indent=2)
                        fixed_count += 1
                        print(f"✓ Fixed: {response.get('name', 'Unnamed')}")
                        
                except json.JSONDecodeError:
                    print(f"✗ Invalid JSON in: {response.get('name', 'Unnamed')}")
                except Exception as e:
                    print(f"✗ Error in {response.get('name', 'Unnamed')}: {e}")
    
    # Process nested items
    if 'item' in item:
        for subitem in item['item']:
            fixed_count += fix_pagination_response(subitem)
    
    return fixed_count

# Read collection
with open('MyPOSCore.postman_collection.json', 'r') as f:
    collection = json.load(f)

print("Fixing pagination responses...")
print("=" * 60)

total_fixed = 0
if 'item' in collection:
    for item in collection['item']:
        total_fixed += fix_pagination_response(item)

# Write back
with open('MyPOSCore.postman_collection.json', 'w') as f:
    json.dump(collection, f, indent='\t')

print("=" * 60)
print(f"✓ Fixed {total_fixed} responses")
print("✓ Postman collection updated")
