# API Testing Guide

This file contains example curl commands for testing all API endpoints.

## Authentication

### Register a new user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "johndoe",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "john@example.com",
    "password": "password123"
  }'
```

**Save the token from the response for subsequent requests:**
```bash
export TOKEN="your-jwt-token-here"
```

## Categories

### Get all categories
```bash
curl -X GET http://localhost:8080/api/categories \
  -H "Authorization: Bearer $TOKEN"
```

### Get income categories only
```bash
curl -X GET "http://localhost:8080/api/categories?type=income" \
  -H "Authorization: Bearer $TOKEN"
```

### Get expense categories only
```bash
curl -X GET "http://localhost:8080/api/categories?type=expense" \
  -H "Authorization: Bearer $TOKEN"
```

## Income

### Create income
```bash
curl -X POST http://localhost:8080/api/income \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "category_id": 1,
    "amount": 5000.00,
    "description": "Monthly salary",
    "income_date": "2025-12-24"
  }'
```

### Get all income
```bash
curl -X GET http://localhost:8080/api/income \
  -H "Authorization: Bearer $TOKEN"
```

### Get income by date
```bash
curl -X GET "http://localhost:8080/api/income?date=2025-12-24" \
  -H "Authorization: Bearer $TOKEN"
```

### Get income by date range
```bash
curl -X GET "http://localhost:8080/api/income?start_date=2025-12-01&end_date=2025-12-31" \
  -H "Authorization: Bearer $TOKEN"
```

### Get income by category
```bash
curl -X GET "http://localhost:8080/api/income?category_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Update income
```bash
curl -X PUT http://localhost:8080/api/income/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "category_id": 1,
    "amount": 5500.00,
    "description": "Updated salary",
    "income_date": "2025-12-24"
  }'
```

### Delete income
```bash
curl -X DELETE http://localhost:8080/api/income/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Expense

### Create expense
```bash
curl -X POST http://localhost:8080/api/expense \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "category_id": 5,
    "amount": 150.00,
    "description": "Weekly groceries",
    "expense_date": "2025-12-24"
  }'
```

### Get all expenses
```bash
curl -X GET http://localhost:8080/api/expense \
  -H "Authorization: Bearer $TOKEN"
```

### Get expenses by date
```bash
curl -X GET "http://localhost:8080/api/expense?date=2025-12-24" \
  -H "Authorization: Bearer $TOKEN"
```

### Get expenses by date range
```bash
curl -X GET "http://localhost:8080/api/expense?start_date=2025-12-01&end_date=2025-12-31" \
  -H "Authorization: Bearer $TOKEN"
```

### Get expenses by category
```bash
curl -X GET "http://localhost:8080/api/expense?category_id=5" \
  -H "Authorization: Bearer $TOKEN"
```

### Update expense
```bash
curl -X PUT http://localhost:8080/api/expense/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "category_id": 5,
    "amount": 175.00,
    "description": "Updated groceries amount",
    "expense_date": "2025-12-24"
  }'
```

### Delete expense
```bash
curl -X DELETE http://localhost:8080/api/expense/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Dashboard

### Get dashboard summary
```bash
curl -X GET http://localhost:8080/api/dashboard \
  -H "Authorization: Bearer $TOKEN"
```

This returns:
- Total income
- Total expense
- Balance
- Today's income and expense
- Monthly income and expense
- Last 30 days daily data (for charts)
- Category breakdown

## Export

### Export to PDF (last 30 days)
```bash
curl -X GET "http://localhost:8080/api/export/pdf?start_date=2025-11-24&end_date=2025-12-24" \
  -H "Authorization: Bearer $TOKEN" \
  --output report.pdf
```

### Export to PDF (custom date range)
```bash
curl -X GET "http://localhost:8080/api/export/pdf?start_date=2025-12-01&end_date=2025-12-31" \
  -H "Authorization: Bearer $TOKEN" \
  --output december_report.pdf
```

## Testing Workflow Example

```bash
# 1. Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","username":"testuser","password":"test123"}'

# 2. Login and save token
TOKEN=$(curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email_or_username":"test@test.com","password":"test123"}' \
  | jq -r '.token')

# 3. Create some income
curl -X POST http://localhost:8080/api/income \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"category_id":1,"amount":5000,"description":"Salary","income_date":"2025-12-24"}'

# 4. Create some expenses
curl -X POST http://localhost:8080/api/expense \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"category_id":5,"amount":200,"description":"Food","expense_date":"2025-12-24"}'

curl -X POST http://localhost:8080/api/expense \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"category_id":7,"amount":1200,"description":"Rent","expense_date":"2025-12-24"}'

# 5. View dashboard
curl -X GET http://localhost:8080/api/dashboard \
  -H "Authorization: Bearer $TOKEN" | jq

# 6. Export report
curl -X GET "http://localhost:8080/api/export/pdf?start_date=2025-12-01&end_date=2025-12-31" \
  -H "Authorization: Bearer $TOKEN" \
  --output report.pdf
```

## Error Responses

All endpoints return JSON error responses:

```json
{
  "error": "error message here"
}
```

Common HTTP status codes:
- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists (e.g., duplicate email)
- `500 Internal Server Error`: Server error
