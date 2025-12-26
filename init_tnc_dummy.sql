-- Dummy Terms and Conditions for MyPOS Core
-- Run this after table creation

INSERT INTO terms_and_conditions (title, content, version, is_active, created_at, updated_at) VALUES
('Terms and Conditions - MyPOS Core System', 
'# Terms and Conditions for MyPOS Core

## 1. Introduction
Welcome to MyPOS Core, a comprehensive Point of Sale management system. By using this system, you agree to be bound by these Terms and Conditions.

## 2. Service Description
MyPOS Core provides:
- Multi-tenant POS system management
- Branch and user management
- Product catalog and inventory tracking
- Order processing and payment management
- Reporting and analytics

## 3. User Responsibilities
Users are responsible for:
- Maintaining the confidentiality of their login credentials
- All activities conducted under their account
- Ensuring accurate data entry
- Complying with applicable laws and regulations

## 4. Data Privacy
We are committed to protecting your data:
- All transactions are encrypted
- Personal information is stored securely
- We do not share your data with third parties without consent
- Regular backups are performed to prevent data loss

## 5. Payment Terms
- Payment processing fees may apply
- All prices are in Indonesian Rupiah (IDR)
- Refunds are subject to our refund policy
- Subscription fees are non-refundable

## 6. System Availability
- We strive for 99.9% uptime
- Scheduled maintenance will be announced in advance
- We are not liable for service interruptions beyond our control

## 7. Limitation of Liability
MyPOS Core is provided "as is" without warranties. We are not liable for:
- Loss of business or profits
- Data loss due to user error
- Unauthorized access to accounts
- Third-party service failures

## 8. Modifications
We reserve the right to modify these terms at any time. Users will be notified of significant changes.

## 9. Termination
We may terminate or suspend access to our service:
- For violation of these terms
- For non-payment
- At our sole discretion with notice

## 10. Contact Information
For questions about these terms, please contact our support team.

Last Updated: December 2025',
'1.0.0',
true,
NOW(),
NOW()
),

('Privacy Policy', 
'# Privacy Policy

## Information We Collect
- Business information (company name, address, tax ID)
- User information (name, email, phone, role)
- Transaction data (orders, payments, products)
- Usage data (login times, feature usage)

## How We Use Your Information
- To provide and maintain the service
- To process transactions
- To improve user experience
- To communicate important updates
- To provide customer support

## Data Security
- SSL/TLS encryption for data transmission
- Encrypted database storage
- Regular security audits
- Access controls and authentication
- Regular backups

## Data Retention
- Transaction data: 7 years (tax compliance)
- User data: Duration of account + 1 year
- System logs: 90 days

## Your Rights
- Access your personal data
- Request data correction
- Request data deletion
- Export your data
- Opt-out of marketing communications

## Third-Party Services
We may use third-party services for:
- Payment processing
- Cloud hosting
- Analytics

## Changes to Privacy Policy
We will notify users of any material changes to this policy.

Last Updated: December 2025',
'1.0.0',
true,
NOW(),
NOW()
),

('Refund Policy',
'# Refund Policy

## Subscription Refunds
- Subscription fees are generally non-refundable
- Pro-rata refunds may be considered for annual plans (within 14 days)
- No refunds for partial months

## Transaction Refunds
- Customer refunds are handled at the merchant''s discretion
- Processing fees are non-refundable
- Refunds typically process within 5-7 business days

## Cancellation
- You may cancel your subscription at any time
- Access continues until the end of the billing period
- No refund for the current billing period

## Special Circumstances
Refunds may be considered for:
- System downtime exceeding 48 hours
- Critical bugs preventing system use
- Duplicate charges

## How to Request a Refund
Contact our support team with:
- Account information
- Reason for refund request
- Supporting documentation

Last Updated: December 2025',
'1.0.0',
true,
NOW(),
NOW()
);
