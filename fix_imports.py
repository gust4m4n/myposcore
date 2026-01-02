#!/usr/bin/env python3
import re

files_to_fix = {
    "handlers/admin_change_password_handler.go": {"remove_http": True},
    "handlers/admin_change_pin_handler.go": {"remove_http": True},
    "handlers/dev_handler.go": {"remove_utils": True},
    "handlers/faq_handler.go": {"remove_utils": True},
    "handlers/health_handler.go": {"remove_utils": True},
    "handlers/tnc_handler.go": {"remove_utils": True},
}

for file_path, actions in files_to_fix.items():
    with open(file_path, 'r') as f:
        content = f.read()
    
    if actions.get("remove_http"):
        content = re.sub(r'\t"net/http"\n', '', content)
    
    if actions.get("remove_utils"):
        content = re.sub(r'\t"myposcore/utils"\n', '', content)
    
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Fixed: {file_path}")

print("Done!")
