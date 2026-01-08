#!/usr/bin/env python3
"""Parse registered routes from server logs and compare with handlers"""

import re
import subprocess
import os

# Parse routes from server startup log
routes_output = """
[GIN-debug] GET    /api/v1/tnc               --> myposcore/handlers.(*TnCHandler).GetTnC-fm
[GIN-debug] GET    /api/v1/faq               --> myposcore/handlers.(*FAQHandler).GetAllFAQ-fm
[GIN-debug] GET    /api/v1/faq/:id           --> myposcore/handlers.(*FAQHandler).GetFAQByID-fm
[GIN-debug] POST   /api/v1/config/set        --> myposcore/handlers.(*ConfigHandler).SetConfig-fm
[GIN-debug] GET    /api/v1/config/get/:key   --> myposcore/handlers.(*ConfigHandler).GetConfig-fm
[GIN-debug] POST   /api/v1/logout            --> myposcore/handlers.(*LogoutHandler).Handle-fm
[GIN-debug] GET    /api/v1/profile           --> myposcore/handlers.(*ProfileHandler).Handle-fm
[GIN-debug] PUT    /api/v1/profile           --> myposcore/handlers.(*ProfileHandler).UpdateProfile-fm
[GIN-debug] PUT    /api/v1/change-password   --> myposcore/handlers.(*ChangePasswordHandler).Handle-fm
[GIN-debug] PUT    /api/v1/admin/change-password --> myposcore/handlers.(*AdminChangePasswordHandler).Handle-fm
[GIN-debug] POST   /api/v1/profile/photo     --> myposcore/handlers.(*ProfileHandler).UploadProfileImage-fm
[GIN-debug] DELETE /api/v1/profile/photo     --> myposcore/handlers.(*ProfileHandler).DeleteProfileImage-fm
[GIN-debug] POST   /api/v1/pin/create        --> myposcore/handlers.(*PINHandler).CreatePIN-fm
[GIN-debug] PUT    /api/v1/pin/change        --> myposcore/handlers.(*PINHandler).ChangePIN-fm
[GIN-debug] GET    /api/v1/pin/check         --> myposcore/handlers.(*PINHandler).CheckPIN-fm
[GIN-debug] PUT    /api/v1/admin/change-pin  --> myposcore/handlers.(*AdminChangePINHandler).Handle-fm
[GIN-debug] GET    /api/v1/branches          --> myposcore/handlers.(*BranchHandler).GetBranches-fm
[GIN-debug] GET    /api/v1/branches/:id      --> myposcore/handlers.(*BranchHandler).GetBranch-fm
[GIN-debug] POST   /api/v1/branches          --> myposcore/handlers.(*BranchHandler).CreateBranch-fm
[GIN-debug] PUT    /api/v1/branches/:id      --> myposcore/handlers.(*BranchHandler).UpdateBranch-fm
[GIN-debug] DELETE /api/v1/branches/:id      --> myposcore/handlers.(*BranchHandler).DeleteBranch-fm
[GIN-debug] GET    /api/v1/branches/:id/users --> myposcore/handlers.(*BranchHandler).GetBranchUsers-fm
[GIN-debug] GET    /api/v1/categories        --> myposcore/handlers.(*CategoryHandler).ListCategories-fm
[GIN-debug] GET    /api/v1/categories/:id    --> myposcore/handlers.(*CategoryHandler).GetCategory-fm
[GIN-debug] POST   /api/v1/categories        --> myposcore/handlers.(*CategoryHandler).CreateCategory-fm
[GIN-debug] PUT    /api/v1/categories/:id    --> myposcore/handlers.(*CategoryHandler).UpdateCategory-fm
[GIN-debug] DELETE /api/v1/categories/:id    --> myposcore/handlers.(*CategoryHandler).DeleteCategory-fm
[GIN-debug] GET    /api/v1/products/categories --> myposcore/handlers.(*ProductHandler).GetCategories-fm
[GIN-debug] GET    /api/v1/products/by-category/:category_id --> myposcore/handlers.(*ProductHandler).ListProductsByCategoryID-fm
[GIN-debug] GET    /api/v1/products          --> myposcore/handlers.(*ProductHandler).ListProducts-fm
[GIN-debug] GET    /api/v1/products/:id      --> myposcore/handlers.(*ProductHandler).GetProduct-fm
[GIN-debug] POST   /api/v1/products          --> myposcore/handlers.(*ProductHandler).CreateProduct-fm
[GIN-debug] DELETE /api/v1/products/:id/photo --> myposcore/handlers.(*ProductHandler).DeleteProductImage-fm
[GIN-debug] POST   /api/v1/orders            --> myposcore/handlers.(*OrderHandler).CreateOrder-fm
[GIN-debug] GET    /api/v1/orders            --> myposcore/handlers.(*OrderHandler).ListOrders-fm
[GIN-debug] GET    /api/v1/orders/:id        --> myposcore/handlers.(*OrderHandler).GetOrder-fm
[GIN-debug] GET    /api/v1/orders/:id/payments --> myposcore/handlers.(*PaymentHandler).GetPaymentsByOrder-fm
[GIN-debug] POST   /api/v1/payments          --> myposcore/handlers.(*PaymentHandler).CreatePayment-fm
[GIN-debug] GET    /api/v1/payments          --> myposcore/handlers.(*PaymentHandler).ListPayments-fm
[GIN-debug] GET    /api/v1/payments/:id      --> myposcore/handlers.(*PaymentHandler).GetPayment-fm
[GIN-debug] GET    /api/v1/payments/performance --> myposcore/handlers.(*PaymentHandler).GetPaymentPerformance-fm
[GIN-debug] GET    /api/v1/users             --> myposcore/handlers.(*UserHandler).ListUsers-fm
[GIN-debug] GET    /api/v1/users/:id         --> myposcore/handlers.(*UserHandler).GetUser-fm
[GIN-debug] POST   /api/v1/users             --> myposcore/handlers.(*UserHandler).CreateUser-fm
[GIN-debug] PUT    /api/v1/users/:id         --> myposcore/handlers.(*UserHandler).UpdateUser-fm
[GIN-debug] DELETE /api/v1/users/:id         --> myposcore/handlers.(*UserHandler).DeleteUser-fm
[GIN-debug] GET    /api/v1/tenants           --> myposcore/handlers.(*TenantHandler).ListTenants-fm
[GIN-debug] GET    /api/v1/tenants/:id       --> myposcore/handlers.(*TenantHandler).GetTenant-fm
[GIN-debug] POST   /api/v1/tenants           --> myposcore/handlers.(*TenantHandler).CreateTenant-fm
[GIN-debug] PUT    /api/v1/tenants/:id       --> myposcore/handlers.(*TenantHandler).UpdateTenant-fm
[GIN-debug] DELETE /api/v1/tenants/:id       --> myposcore/handlers.(*TenantHandler).DeleteTenant-fm
[GIN-debug] GET    /api/v1/audit-trails      --> myposcore/handlers.(*AuditTrailHandler).ListAuditTrails-fm
[GIN-debug] GET    /api/v1/audit-trails/user/:user_id --> myposcore/handlers.(*AuditTrailHandler).GetUserActivityLog-fm
[GIN-debug] GET    /api/v1/audit-trails/entity/:entity_type/:entity_id --> myposcore/handlers.(*AuditTrailHandler).GetEntityAuditHistory-fm
[GIN-debug] GET    /api/v1/audit-trails/:id  --> myposcore/handlers.(*AuditTrailHandler).GetAuditTrailByID-fm
[GIN-debug] GET    /api/v1/dashboard         --> myposcore/handlers.(*SuperAdminHandler).Dashboard-fm
[GIN-debug] POST   /api/v1/faq               --> myposcore/handlers.(*FAQHandler).CreateFAQ-fm
[GIN-debug] PUT    /api/v1/faq/:id           --> myposcore/handlers.(*FAQHandler).UpdateFAQ-fm
[GIN-debug] DELETE /api/v1/faq/:id           --> myposcore/handlers.(*FAQHandler).DeleteFAQ-fm
"""

registered_routes = {}
for line in routes_output.strip().split('\n'):
    match = re.search(r'\[GIN-debug\] (\w+)\s+(.+?)\s+-->', line)
    if match:
        method, path = match.groups()
        path = path.strip()
        key = f"{method} {path}"
        registered_routes[key] = True

# Handler methods inventory
handlers_info = {
    "ProductHandler": [
        "ListProducts",
        "ListProductsByCategoryID",
        "GetProduct",
        "CreateProduct",
        "UpdateProduct",
        "DeleteProduct",
        "GetCategories",
        "UploadProductImage",
        "DeleteProductImage"
    ],
    "CategoryHandler": [
        "ListCategories",
        "GetCategory",
        "CreateCategory",
        "UpdateCategory",
        "DeleteCategory"
    ],
    "BranchHandler": [
        "GetBranches",
        "GetBranch",
        "CreateBranch",
        "UpdateBranch",
        "DeleteBranch",
        "GetBranchUsers"
    ],
    "TenantHandler": [
        "ListTenants",
        "GetTenant",
        "CreateTenant",
        "UpdateTenant",
        "DeleteTenant"
    ],
    "OrderHandler": [
        "CreateOrder",
        "ListOrders",
        "GetOrder"
    ],
    "PaymentHandler": [
        "CreatePayment",
        "ListPayments",
        "GetPayment",
        "GetPaymentsByOrder",
        "GetPaymentPerformance"
    ],
    "UserHandler": [
        "ListUsers",
        "GetUser",
        "CreateUser",
        "UpdateUser",
        "DeleteUser"
    ],
    "ProfileHandler": [
        "Handle",  # GET
        "UpdateProfile",  # PUT
        "UploadProfileImage",
        "DeleteProfileImage"
    ],
    "FAQHandler": [
        "GetAllFAQ",
        "GetFAQByID",
        "CreateFAQ",
        "UpdateFAQ",
        "DeleteFAQ"
    ],
    "TnCHandler": [
        "GetTnC"
    ],
    "AuditTrailHandler": [
        "ListAuditTrails",
        "GetAuditTrailByID",
        "GetUserActivityLog",
        "GetEntityAuditHistory"
    ]
}

# Expected routes based on RESTful conventions
expected_routes = {
    "ProductHandler": {
        "GET /api/v1/products": "ListProducts",
        "GET /api/v1/products/:id": "GetProduct",
        "POST /api/v1/products": "CreateProduct",
        "PUT /api/v1/products/:id": "UpdateProduct",
        "DELETE /api/v1/products/:id": "DeleteProduct",
        "GET /api/v1/products/by-category/:category_id": "ListProductsByCategoryID",
        "GET /api/v1/products/categories": "GetCategories",
        "POST /api/v1/products/:id/photo": "UploadProductImage",
        "DELETE /api/v1/products/:id/photo": "DeleteProductImage"
    },
    "CategoryHandler": {
        "GET /api/v1/categories": "ListCategories",
        "GET /api/v1/categories/:id": "GetCategory",
        "POST /api/v1/categories": "CreateCategory",
        "PUT /api/v1/categories/:id": "UpdateCategory",
        "DELETE /api/v1/categories/:id": "DeleteCategory"
    },
    "BranchHandler": {
        "GET /api/v1/branches": "GetBranches",
        "GET /api/v1/branches/:id": "GetBranch",
        "POST /api/v1/branches": "CreateBranch",
        "PUT /api/v1/branches/:id": "UpdateBranch",
        "DELETE /api/v1/branches/:id": "DeleteBranch",
        "GET /api/v1/branches/:id/users": "GetBranchUsers"
    },
    "TenantHandler": {
        "GET /api/v1/tenants": "ListTenants",
        "GET /api/v1/tenants/:id": "GetTenant",
        "POST /api/v1/tenants": "CreateTenant",
        "PUT /api/v1/tenants/:id": "UpdateTenant",
        "DELETE /api/v1/tenants/:id": "DeleteTenant"
    }
}

# Check missing routes
print("=" * 80)
print("ROUTE REGISTRATION STATUS")
print("=" * 80)
print()

missing_routes = []
for handler_name, routes in expected_routes.items():
    print(f"\n{handler_name}:")
    print("-" * 80)
    for route, method_name in routes.items():
        if route in registered_routes:
            print(f"  ✅ {route:<60} -> {method_name}")
        else:
            print(f"  ❌ {route:<60} -> {method_name} MISSING!")
            missing_routes.append((handler_name, route, method_name))

if missing_routes:
    print("\n" + "=" * 80)
    print("MISSING ROUTES SUMMARY")
    print("=" * 80)
    for handler, route, method in missing_routes:
        print(f"  {handler}: {route} -> {method}")
    print(f"\nTotal missing: {len(missing_routes)} route(s)")
else:
    print("\n" + "=" * 80)
    print("✅ ALL ROUTES ARE REGISTERED!")
    print("=" * 80)

print(f"\nTotal registered routes: {len(registered_routes)}")
