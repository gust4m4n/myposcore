#!/usr/bin/env python3
"""
Update all Postman collection response examples to use standardized format:
- Add "code" field (0 for success, 1-7 for errors)
- Add "message" field for all responses
- Wrap data in "data" field for success responses
"""

import json
import re

def update_response_body(body_str, status_code):
    """Update response body to include code and message fields"""
    try:
        # Parse existing body
        body = json.loads(body_str)
        
        # Skip if already has code field
        if isinstance(body, dict) and 'code' in body:
            return body_str
        
        # Determine code based on status
        if status_code >= 200 and status_code < 300:
            code = 0
            # Success response
            if isinstance(body, dict):
                # If already has message and data, wrap properly
                if 'message' in body and 'data' in body:
                    return json.dumps({
                        "code": 0,
                        "message": body['message'],
                        "data": body['data']
                    }, indent=2)
                # If has message only
                elif 'message' in body:
                    return json.dumps({
                        "code": 0,
                        "message": body['message']
                    }, indent=2)
                # If has data but no message
                elif 'data' in body:
                    return json.dumps({
                        "code": 0,
                        "message": "Operation successful",
                        "data": body['data']
                    }, indent=2)
                # If has other fields, wrap in data
                else:
                    # Determine message from context
                    message = "Operation successful"
                    if 'token' in body:
                        message = "Login successful"
                    elif 'status' in body:
                        message = "Status retrieved successfully"
                    
                    return json.dumps({
                        "code": 0,
                        "message": message,
                        "data": body
                    }, indent=2)
        else:
            # Error response
            code = 1  # Default bad request
            if status_code == 401:
                code = 2  # Unauthorized
            elif status_code == 403:
                code = 3  # Forbidden
            elif status_code == 404:
                code = 4  # Not found
            elif status_code == 500:
                code = 5  # Internal error
            elif status_code == 409:
                code = 6  # Conflict
            elif status_code == 422:
                code = 7  # Unprocessable
            
            if isinstance(body, dict):
                # If has error field, use it as message
                if 'error' in body:
                    return json.dumps({
                        "code": code,
                        "message": body['error']
                    }, indent=2)
                # If already has message
                elif 'message' in body:
                    return json.dumps({
                        "code": code,
                        "message": body['message']
                    }, indent=2)
        
        return body_str
    except json.JSONDecodeError:
        return body_str
    except Exception as e:
        print(f"Error updating body: {e}")
        return body_str

def update_collection_responses(collection):
    """Recursively update all response examples in collection"""
    updated_count = 0
    
    def process_item(item):
        nonlocal updated_count
        
        # Process responses in this item
        if 'response' in item:
            for response in item['response']:
                if 'code' in response and 'body' in response:
                    status_code = response['code']
                    old_body = response['body']
                    new_body = update_response_body(old_body, status_code)
                    if new_body != old_body:
                        response['body'] = new_body
                        updated_count += 1
        
        # Process nested items
        if 'item' in item:
            for subitem in item['item']:
                process_item(subitem)
    
    # Process all top-level items
    if 'item' in collection:
        for item in collection['item']:
            process_item(item)
    
    return updated_count

# Read Postman collection
with open('MyPOSCore.postman_collection.json', 'r') as f:
    collection = json.load(f)

print("Updating Postman collection responses...")
updated_count = update_collection_responses(collection)

# Write updated collection
with open('MyPOSCore.postman_collection.json', 'w') as f:
    json.dump(collection, f, indent='\t')

print(f"✓ Updated {updated_count} response examples")
print("✓ Postman collection updated successfully!")
