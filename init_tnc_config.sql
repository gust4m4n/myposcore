-- Insert TNC content into config table
INSERT INTO configs (key, value, created_at, updated_at)
VALUES (
    'tnc',
    '# Terms and Conditions - MyPOS Core System

## 1. Introduction
Welcome to MyPOS Core, a comprehensive Point of Sale system designed for multi-tenant and multi-branch businesses.

## 2. Acceptance of Terms
By accessing and using MyPOS Core, you agree to be bound by these Terms and Conditions.

## 3. Service Description
MyPOS Core provides:
- Multi-tenant and multi-branch POS system
- User management with role-based access control
- Product and inventory management
- Order and payment processing
- Dashboard and reporting features

## 4. User Obligations
- Maintain the confidentiality of your account credentials
- Use the service in compliance with applicable laws
- Not attempt to gain unauthorized access to the system
- Report any security vulnerabilities to our team

## 5. Data Privacy
We are committed to protecting your data:
- All sensitive data is encrypted
- We do not share your data with third parties without consent
- You retain ownership of your business data

## 6. Service Availability
- We strive for 99.9% uptime
- Scheduled maintenance will be announced in advance
- We are not liable for service interruptions beyond our control

## 7. Payment Terms
- Subscription fees are billed monthly or annually
- Payment must be made in advance
- Refunds are handled on a case-by-case basis

## 8. Limitation of Liability
MyPOS Core is provided "as is" without warranties. We are not liable for:
- Loss of data or business interruption
- Indirect or consequential damages
- Issues arising from misuse of the system

## 9. Termination
- Either party may terminate the agreement with 30 days notice
- Upon termination, you will have 60 days to export your data
- Outstanding fees must be paid before termination

## 10. Changes to Terms
- We reserve the right to modify these terms
- Users will be notified of material changes
- Continued use constitutes acceptance of modified terms

## 11. Governing Law
These terms are governed by the laws of Indonesia.

## 12. Contact Information
For questions about these terms, contact us at:
- Email: support@myposcore.com
- Phone: +62 xxx xxxx xxxx

Last Updated: January 3, 2026
Version: 1.0.0',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = CURRENT_TIMESTAMP;
