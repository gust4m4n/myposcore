#!/usr/bin/env python3
"""
Verify and fix Postman collection response format.
Check that all responses have proper code and message fields in body.
"""

import json
import re

def check_response_body(body_str, status_code, name):
    """Check if response body has proper format"""
    try:
        body = json.loads(body_str)
        
        if not isinstance(body, dict):
            return False, f"Body is not a dict: {name}"
        
        # Check if has code field
        if 'code' not in body:
            return False, f"Missing 'code' field: {name}"
        
        # Check if has message field (except for health check which has status)
        if 'message' not in body and 'status' not in body:
            return False, f"Missing 'message' field: {name}"
        
        # Verify code matches status
        expected_code = 0 if (status_code >= 200 and status_code < 300) else 1
        if status_code == 401:
            expected_code = 2
        elif status_code == 403:
            expected_code = 3
        elif status_code == 404:
            expected_code = 4
        elif status_code == 500:
            expected_code = 5
        elif status_code == 409:
            expected_code = 6
        elif status_code == 422:
            expected_code = 7
        
        actual_code = body.get('code')
        if actual_code != expected_code:
            return False, f"Code mismatch: {name} (expected {expected_code}, got {actual_code})"
        
        return True, "OK"
        
    except json.JSONDecodeError:
        return False, f"Invalid JSON: {name}"
    except Exception as e:
        return False, f"Error: {name} - {str(e)}"

def check_collection_responses(collection):
    """Check all response examples in collection"""
    issues = []
    success_count = 0
    total_count = 0
    
    def process_item(item, parent_name=""):
        nonlocal success_count, total_count
        
        item_name = item.get('name', 'Unnamed')
        full_name = f"{parent_name} > {item_name}" if parent_name else item_name
        
        # Process responses in this item
        if 'response' in item:
            for response in item['response']:
                if 'body' in response:
                    total_count += 1
                    response_name = response.get('name', 'Unnamed response')
                    status_code = response.get('code', 200)
                    
                    is_valid, message = check_response_body(
                        response['body'], 
                        status_code,
                        f"{full_name} > {response_name}"
                    )
                    
                    if is_valid:
                        success_count += 1
                    else:
                        issues.append(message)
        
        # Process nested items
        if 'item' in item:
            for subitem in item['item']:
                process_item(subitem, full_name)
    
    # Process all top-level items
    if 'item' in collection:
        for item in collection['item']:
            process_item(item)
    
    return total_count, success_count, issues

# Read Postman collection
with open('MyPOSCore.postman_collection.json', 'r') as f:
    collection = json.load(f)

print("Checking Postman collection response format...")
print("=" * 60)

total, success, issues = check_collection_responses(collection)

print(f"\n✓ Total responses checked: {total}")
print(f"✓ Valid responses: {success}")
print(f"✗ Invalid responses: {len(issues)}")

if issues:
    print("\n❌ Issues found:")
    for issue in issues[:10]:  # Show first 10 issues
        print(f"  - {issue}")
    if len(issues) > 10:
        print(f"  ... and {len(issues) - 10} more")
else:
    print("\n🎉 All responses are properly formatted!")
    print("✓ All responses have 'code' and 'message' fields")
    print("✓ All response codes match HTTP status codes")
