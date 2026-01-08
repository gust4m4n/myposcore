#!/usr/bin/env python3
"""
Update Postman collection with:
1. Add new endpoint "List Products by Category" 
2. Update all product response examples to include category_id and category_detail
3. Update request bodies to include category_id
"""

import json
import sys

print("Loading Postman collection...")
with open('MyPOSCore.postman_collection.json', 'r', encoding='utf-8') as f:
    collection = json.load(f)

# Find the Products section
products_section = None
for item in collection['item']:
    if item.get('name') == 'Products':
        products_section = item
        break

if not products_section:
    print("ERROR: Products section not found!")
    sys.exit(1)

print(f"Found Products section with {len(products_section['item'])} endpoints")

# NEW ENDPOINT: List Products by Category
list_by_category_endpoint = {
    "name": "List Products by Category",
    "request": {
        "auth": {
            "type": "bearer",
            "bearer": [
                {
                    "key": "token",
                    "value": "{{auth_token}}",
                    "type": "string"
                }
            ]
        },
        "method": "GET",
        "header": [],
        "url": {
            "raw": "{{base_url}}/api/v1/products/by-category/13",
            "host": ["{{base_url}}"],
            "path": ["api", "v1", "products", "by-category", "13"],
            "query": [
                {
                    "key": "page",
                    "value": "1",
                    "disabled": True
                },
                {
                    "key": "page_size",
                    "value": "20",
                    "disabled": True
                }
            ]
        },
        "description": "Get products filtered by category ID. Replace '13' with actual category ID. Supports pagination with page and page_size query parameters."
    },
    "response": [
        {
            "name": "List by Category Success",
            "originalRequest": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {{auth_token}}",
                        "type": "text"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/v1/products/by-category/13?page=1&page_size=20",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "products", "by-category", "13"],
                    "query": [
                        {"key": "page", "value": "1"},
                        {"key": "page_size", "value": "20"}
                    ]
                }
            },
            "status": "OK",
            "code": 200,
            "_postman_previewlanguage": "json",
            "header": [{"key": "Content-Type", "value": "application/json"}],
            "cookie": [],
            "body": json.dumps({
                "code": 0,
                "message": "Operation successful",
                "data": [
                    {
                        "id": 73,
                        "tenant_id": 17,
                        "name": "Nasi Goreng Spesial",
                        "description": "Nasi goreng dengan telur, ayam, dan sayuran",
                        "category": "Makanan Utama",
                        "category_id": 13,
                        "category_detail": {
                            "id": 13,
                            "name": "Makanan Utama",
                            "description": "Menu makanan utama dan hidangan pokok",
                            "image": None
                        },
                        "sku": "FOOD-001",
                        "price": 25000,
                        "stock": 50,
                        "image": "/uploads/products/product_73.jpg",
                        "is_active": True,
                        "created_at": "2026-01-08T10:00:00+07:00",
                        "updated_at": "2026-01-08T10:00:00+07:00",
                        "created_by": 1,
                        "created_by_name": "Admin Resto",
                        "updated_by": None,
                        "updated_by_name": None
                    },
                    {
                        "id": 74,
                        "tenant_id": 17,
                        "name": "Mie Goreng",
                        "description": "Mie goreng dengan topping lengkap",
                        "category": "Makanan Utama",
                        "category_id": 13,
                        "category_detail": {
                            "id": 13,
                            "name": "Makanan Utama",
                            "description": "Menu makanan utama dan hidangan pokok",
                            "image": None
                        },
                        "sku": "FOOD-002",
                        "price": 22000,
                        "stock": 45,
                        "image": "/uploads/products/product_74.jpg",
                        "is_active": True,
                        "created_at": "2026-01-08T10:00:00+07:00",
                        "updated_at": "2026-01-08T10:00:00+07:00",
                        "created_by": 1,
                        "created_by_name": "Admin Resto",
                        "updated_by": None,
                        "updated_by_name": None
                    }
                ],
                "pagination": {
                    "page": 1,
                    "page_size": 20,
                    "total_items": 5,
                    "total_pages": 1
                }
            }, indent=2)
        }
    ]
}

# Insert the new endpoint after "Get Categories"
insert_index = None
for idx, endpoint in enumerate(products_section['item']):
    if endpoint.get('name') == 'Get Categories':
        insert_index = idx + 1
        break

if insert_index is not None:
    # Check if already exists
    exists = any(e.get('name') == 'List Products by Category' for e in products_section['item'])
    if not exists:
        products_section['item'].insert(insert_index, list_by_category_endpoint)
        print("✅ Added 'List Products by Category' endpoint")
    else:
        print("⚠️  'List Products by Category' endpoint already exists")

# Update response examples for existing endpoints
updates_made = []

for endpoint in products_section['item']:
    endpoint_name = endpoint.get('name', '')
    
    # Update "List Products" response
    if endpoint_name == 'List Products':
        for response in endpoint.get('response', []):
            if 'body' in response:
                try:
                    body = json.loads(response['body'])
                    if 'data' in body and isinstance(body['data'], list):
                        for product in body['data']:
                            if 'category_id' not in product:
                                product['category_id'] = 13
                                product['category_detail'] = {
                                    "id": 13,
                                    "name": product.get("category", "Makanan Utama"),
                                    "description": "Menu makanan utama dan hidangan pokok",
                                    "image": None
                                }
                        response['body'] = json.dumps(body, indent=2)
                        updates_made.append("List Products response")
                except Exception as e:
                    print(f"⚠️  Error updating List Products response: {e}")
    
    # Update "Get Product by ID" response
    elif endpoint_name == 'Get Product by ID':
        for response in endpoint.get('response', []):
            if 'body' in response:
                try:
                    body = json.loads(response['body'])
                    if 'data' in body and isinstance(body['data'], dict):
                        product = body['data']
                        if 'category_id' not in product:
                            product['category_id'] = 13
                            product['category_detail'] = {
                                "id": 13,
                                "name": product.get("category", "Makanan Utama"),
                                "description": "Menu makanan utama dan hidangan pokok",
                                "image": None
                            }
                        response['body'] = json.dumps(body, indent=2)
                        updates_made.append("Get Product by ID response")
                except Exception as e:
                    print(f"⚠️  Error updating Get Product response: {e}")
    
    # Update "Create Product" request & response
    elif endpoint_name == 'Create Product':
        # Update request body
        if 'request' in endpoint and 'body' in endpoint['request']:
            try:
                body_str = endpoint['request']['body'].get('raw', '')
                if body_str:
                    body = json.loads(body_str)
                    if 'category_id' not in body:
                        body['category_id'] = 13
                        endpoint['request']['body']['raw'] = json.dumps(body, indent=2)
                        updates_made.append("Create Product request body")
            except Exception as e:
                print(f"⚠️  Error updating Create Product request: {e}")
        
        # Update response
        for response in endpoint.get('response', []):
            if 'body' in response:
                try:
                    body = json.loads(response['body'])
                    if 'data' in body and isinstance(body['data'], dict):
                        product = body['data']
                        if 'category_id' not in product:
                            product['category_id'] = 13
                            product['category_detail'] = {
                                "id": 13,
                                "name": product.get("category", "Makanan Utama"),
                                "description": "Menu makanan utama dan hidangan pokok",
                                "image": None
                            }
                        response['body'] = json.dumps(body, indent=2)
                        updates_made.append("Create Product response")
                except Exception as e:
                    print(f"⚠️  Error updating Create Product response: {e}")
    
    # Update "Update Product" request & response
    elif endpoint_name == 'Update Product':
        # Update request body
        if 'request' in endpoint and 'body' in endpoint['request']:
            try:
                body_str = endpoint['request']['body'].get('raw', '')
                if body_str:
                    body = json.loads(body_str)
                    if 'category_id' not in body:
                        body['category_id'] = 13
                        endpoint['request']['body']['raw'] = json.dumps(body, indent=2)
                        updates_made.append("Update Product request body")
            except Exception as e:
                print(f"⚠️  Error updating Update Product request: {e}")
        
        # Update response
        for response in endpoint.get('response', []):
            if 'body' in response:
                try:
                    body = json.loads(response['body'])
                    if 'data' in body and isinstance(body['data'], dict):
                        product = body['data']
                        if 'category_id' not in product:
                            product['category_id'] = 13
                            product['category_detail'] = {
                                "id": 13,
                                "name": product.get("category", "Makanan Utama"),
                                "description": "Menu makanan utama dan hidangan pokok",
                                "image": None
                            }
                        response['body'] = json.dumps(body, indent=2)
                        updates_made.append("Update Product response")
                except Exception as e:
                    print(f"⚠️  Error updating Update Product response: {e}")
    
    # Update "Upload Product Image" response
    elif endpoint_name == 'Upload Product Image':
        for response in endpoint.get('response', []):
            if 'body' in response:
                try:
                    body = json.loads(response['body'])
                    if 'data' in body and isinstance(body['data'], dict):
                        product = body['data']
                        if 'category_id' not in product:
                            product['category_id'] = 13
                            product['category_detail'] = {
                                "id": 13,
                                "name": product.get("category", "Makanan Utama"),
                                "description": "Menu makanan utama dan hidangan pokok",
                                "image": None
                            }
                        response['body'] = json.dumps(body, indent=2)
                        updates_made.append("Upload Product Image response")
                except Exception as e:
                    print(f"⚠️  Error updating Upload Image response: {e}")

# Save updated collection
print("\nSaving updated Postman collection...")
with open('MyPOSCore.postman_collection.json', 'w', encoding='utf-8') as f:
    json.dump(collection, f, indent=2, ensure_ascii=False)

print("\n" + "="*60)
print("✅ Postman collection updated successfully!")
print("="*60)
print(f"\nUpdates made:")
for update in set(updates_made):
    print(f"  ✓ {update}")
print(f"\nTotal updates: {len(updates_made)}")
print("\nBackup file: MyPOSCore.postman_collection.json.backup_*")
