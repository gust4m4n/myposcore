#!/bin/bash

# Batch update all remaining gin.H responses

cd /Users/gustaman/Desktop/GUSTAMAN7/myposcore/handlers

echo "Updating handlers with sed..."

# Update base_handler.go - line 32
sed -i '' '32s|c\.JSON(statusCode, gin\.H{|utils.Error(c, statusCode, code, message); return; /*|' base_handler.go
sed -i '' '33s|"code":    code,|REMOVED|' base_handler.go  
sed -i '' '34s|"message": message,|REMOVED|' base_handler.go
sed -i '' '35s|\t})|\t*/|' base_handler.go

# change_password_handler.go - line 55  
sed -i '' 's/c\.JSON(http\.StatusOK, gin\.H{"message": "Password changed successfully"})/utils.SuccessWithoutData(c, "Password changed successfully")/g' change_password_handler.go

# Update pin_handler.go lines 50, 78, 100
sed -i '' 's/c\.JSON(http\.StatusOK, gin\.H{"message": "PIN created successfully"})/utils.SuccessWithoutData(c, "PIN created successfully")/g' pin_handler.go
sed -i '' 's/c\.JSON(http\.StatusOK, gin\.H{"message": "PIN changed successfully"})/utils.SuccessWithoutData(c, "PIN changed successfully")/g' pin_handler.go

echo "✓ Basic patterns updated"
echo "Note: Complex nested structures in order, payment, product, user, superadmin handlers need manual review"
