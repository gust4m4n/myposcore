#!/usr/bin/env python3
import re

files = [
    "handlers/admin_change_password_handler.go",
    "handlers/admin_change_pin_handler.go",
    "handlers/audit_trail_handler.go",
    "handlers/category_handler.go",
    "handlers/change_password_handler.go",
    "handlers/dev_handler.go",
    "handlers/faq_handler.go",
    "handlers/health_handler.go",
    "handlers/order_handler.go",
    "handlers/payment_handler.go",
    "handlers/pin_handler.go",
    "handlers/product_handler.go",
    "handlers/profile_handler.go",
    "handlers/superadmin_handler.go",
    "handlers/tnc_handler.go",
    "handlers/user_handler.go"
]

for file_path in files:
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Check if utils is already imported
    if '"myposcore/utils"' not in content:
        # Add utils import after package declaration and before the first import
        content = re.sub(
            r'(package handlers\n\nimport \()',
            r'\1\n\t"myposcore/utils"',
            content
        )
        
        with open(file_path, 'w') as f:
            f.write(content)
        print(f"Added utils import to {file_path}")

print("Done!")
