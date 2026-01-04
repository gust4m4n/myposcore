#!/usr/bin/env python3
import json

# Read the Postman collection
with open('MyPOSCore.postman_collection.json', 'r') as f:
    collection = json.load(f)

# Find the User folder
user_folder = None
for folder in collection['item']:
    if folder['name'] == 'User':
        user_folder = folder
        break

if not user_folder:
    print("Error: Could not find User folder")
    exit(1)

# Update tenant endpoints
updates_made = 0

for item in user_folder['item']:
    # Update List Tenants
    if item['name'] == 'List Tenants':
        # Update description
        if 'request' in item and 'description' in item['request']:
            item['request']['description'] = "Get list of all tenants in the system with pagination support. All authenticated users can access.\n\nQuery Parameters:\n- page (optional): Page number, default 1\n- page_size (optional): Items per page, default 32\n\nResponse includes pagination metadata:\n- page: Current page number\n- page_size: Items per page\n- total_items: Total number of tenants\n- total_pages: Total number of pages\n- data: Array of tenant items"
            updates_made += 1
        
        # Update response body with pagination format
        if 'response' in item and len(item['response']) > 0:
            item['response'][0]['body'] = '''{
  "page": 1,
  "page_size": 32,
  "total_items": 2,
  "total_pages": 1,
  "data": [
    {
      "id": 1,
      "name": "Food Corner 99",
      "description": "Restoran modern dengan berbagai pilihan menu internasional dan lokal untuk pengalaman kuliner terbaik",
      "address": "Jl. M.H. Thamrin No. 99, Menteng, Jakarta Pusat, DKI Jakarta 10350",
      "website": "https://www.foodcorner99.com",
      "email": "info@foodcorner99.com",
      "phone": "021-3199-8899",
      "image": "/uploads/tenants/tenant_1_1704000000.png",
      "is_active": true,
      "created_at": "2025-12-24 10:00:00",
      "updated_at": "2025-12-24 10:00:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": 1,
      "updated_by_name": "Super Admin"
    },
    {
      "id": 2,
      "name": "Fashion Store 101",
      "description": "Pusat fashion terlengkap dengan koleksi pakaian, sepatu, dan aksesoris trendy untuk pria dan wanita",
      "address": "Grand Indonesia Shopping Town Tower C, Jl. M.H. Thamrin No. 1, Menteng, Jakarta Pusat, DKI Jakarta 10310",
      "website": "https://www.fashionstore101.com",
      "email": "contact@fashionstore101.com",
      "phone": "021-2358-1010",
      "image": "/uploads/tenants/tenant_2_1704001000.png",
      "is_active": true,
      "created_at": "2025-12-24 10:00:00",
      "updated_at": "2025-12-25 14:30:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": 2,
      "updated_by_name": "John Admin"
    }
  ]
}'''
            updates_made += 1
    
    # Update Get Tenant
    elif item['name'] == 'Get Tenant':
        if 'request' in item:
            item['request']['description'] = "Get detailed information about a specific tenant by ID. All authenticated users can access."
            updates_made += 1
    
    # Update Create Tenant
    elif item['name'] == 'Create Tenant':
        if 'request' in item and 'description' in item['request']:
            desc = item['request']['description']
            # Remove superadmin only mentions
            desc = desc.replace('(superadmin only)', '(all users)')
            desc = desc.replace('Only superadmin can access.', 'All authenticated users can access.')
            desc = desc.replace('Hanya superadmin yang bisa akses.', 'Semua user yang terautentikasi dapat mengakses.')
            item['request']['description'] = desc
            updates_made += 1
        
        # Update response body
        if 'response' in item and len(item['response']) > 0:
            item['response'][0]['body'] = '''{
  "code": 0,
  "message": "Tenant created successfully",
  "data": {
    "id": 3,
    "name": "New Tenant Name",
    "description": "Tenant description",
    "address": "Jl. Example No. 123, Jakarta",
    "website": "https://newtenant.com",
    "email": "info@newtenant.com",
    "phone": "021-12345678",
    "image": "/uploads/tenants/tenant_3_1735987654.png",
    "is_active": true,
    "created_at": "2026-01-05 10:00:00",
    "updated_at": "2026-01-05 10:00:00",
    "created_by": 2
  }
}'''
            updates_made += 1
    
    # Update Update Tenant
    elif item['name'] == 'Update Tenant':
        if 'request' in item and 'description' in item['request']:
            desc = item['request']['description']
            desc = desc.replace('(superadmin only)', '(all users)')
            desc = desc.replace('Only superadmin can access.', 'All authenticated users can access.')
            desc = desc.replace('Hanya superadmin yang bisa akses.', 'Semua user yang terautentikasi dapat mengakses.')
            item['request']['description'] = desc
            updates_made += 1
        
        # Update response body
        if 'response' in item and len(item['response']) > 0:
            item['response'][0]['body'] = '''{
  "code": 0,
  "message": "Tenant updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Tenant Name",
    "description": "Updated description",
    "address": "Updated Address, Jakarta",
    "website": "https://updated-tenant.com",
    "email": "updated@tenant.com",
    "phone": "021-99999999",
    "image": "/uploads/tenants/tenant_1_1735987999.png",
    "is_active": true,
    "created_at": "2025-12-24 10:00:00",
    "updated_at": "2026-01-05 15:30:00",
    "updated_by": 2
  }
}'''
            updates_made += 1
    
    # Update Delete Tenant
    elif item['name'] == 'Delete Tenant':
        if 'request' in item and 'description' in item['request']:
            desc = item['request']['description']
            desc = desc.replace('(superadmin only)', '(all users)')
            desc = desc.replace('Only superadmin can access.', 'All authenticated users can access.')
            desc = desc.replace('Hanya superadmin yang bisa akses.', 'Semua user yang terautentikasi dapat mengakses.')
            item['request']['description'] = desc
            updates_made += 1
        
        # Update response body
        if 'response' in item and len(item['response']) > 0:
            item['response'][0]['body'] = '''{
  "code": 0,
  "message": "Tenant deleted successfully"
}'''
            updates_made += 1

# Write the updated collection back
with open('MyPOSCore.postman_collection.json', 'w') as f:
    json.dump(collection, f, indent=2)

print(f"✓ Updated {updates_made} tenant endpoint fields")
print("✓ All tenant endpoints now show they are accessible to all authenticated users")
print("✓ Response examples updated with proper format")
print("✓ List Tenants now uses pagination format")
