# Test Scenarios

This document provides test scenarios to verify the application functionality.

## Test Scenario 1: User Registration and Login

### Steps:
1. Open `http://localhost:8080`
2. Click "Register here"
3. Fill in:
   - Email: `alice@example.com`
   - Username: `alice`
   - Password: `alice123`
4. Click "Register"

### Expected Result:
- Success message appears
- Automatically redirected to dashboard
- Welcome message shows "Welcome, alice"

### Verify:
- Dashboard loads successfully
- Summary cards show $0.00 for all values
- Chart is empty

---

## Test Scenario 2: Add Income

### Steps:
1. Click "Add Income" button
2. Fill in modal:
   - Category: Salary
   - Amount: 5000
   - Date: Today's date
   - Description: "December salary"
3. Click "Save"

### Expected Result:
- Modal closes
- Total Income card updates to $5,000.00
- Balance updates to $5,000.00
- Transaction appears in list

### Verify:
- Chart shows data point for today
- Today's Income shows $5,000.00

---

## Test Scenario 3: Add Multiple Expenses

### Steps:
1. Click "Add Expense"
2. Add Expense 1:
   - Category: Rent
   - Amount: 1200
   - Date: Today
   - Description: "Monthly rent"
   
3. Click "Add Expense" again
4. Add Expense 2:
   - Category: Food
   - Amount: 350
   - Date: Today
   - Description: "Groceries and dining"

5. Click "Add Expense" again
6. Add Expense 3:
   - Category: Transport
   - Amount: 150
   - Date: Today
   - Description: "Gas and public transit"

### Expected Result:
- Total Expense: $1,700.00
- Balance: $3,300.00 (5000 - 1700)
- All three expenses appear in transaction list

### Verify:
- Today's Expense: $1,700.00
- Chart shows expense line for today
- Transaction list shows all entries

---

## Test Scenario 4: Filter Transactions

### Steps:
1. Select "Income" from transaction type dropdown
2. Observe filtered list

### Expected Result:
- Only income transactions visible (Salary entry)

### Steps:
3. Select "Expense" from transaction type dropdown
4. Select "Food" from category filter

### Expected Result:
- Only food expense visible ($350)

### Steps:
5. Click "Clear Filters"

### Expected Result:
- All expense transactions visible again

---

## Test Scenario 5: Edit Transaction

### Steps:
1. Find a transaction in the list
2. Click the ✏️ (edit) icon
3. Change amount from 350 to 400
4. Change description
5. Click "Save"

### Expected Result:
- Modal closes
- Transaction updates in list
- Total amounts recalculated
- Balance updated

---

## Test Scenario 6: Delete Transaction

### Steps:
1. Find a transaction
2. Click 🗑️ (delete) icon
3. Confirm deletion in popup

### Expected Result:
- Transaction removed from list
- Totals recalculated
- Balance updated

---

## Test Scenario 7: Date Filtering

### Steps:
1. Add transactions for different dates:
   - Income: 2000 on Dec 20
   - Expense: 500 on Dec 21
   - Income: 1500 on Dec 22

2. Use date filter to select Dec 21

### Expected Result:
- Only Dec 21 transactions visible
- Total shows only for that day

---

## Test Scenario 8: Export PDF Report

### Steps:
1. Click "Export Report (PDF)"
2. When prompted:
   - Start date: First day of current month
   - End date: Today's date
3. Browser prompts to download PDF

### Expected Result:
- PDF file downloads
- PDF contains:
  - Report header with date range
  - Summary section (total income, expense, balance)
  - Income details table
  - Expense details table

### Verify:
- All transactions within date range included
- Amounts calculated correctly
- Categories displayed properly

---

## Test Scenario 9: Dashboard Chart

### Steps:
1. Add transactions across multiple days
2. Observe the chart

### Expected Result:
- Chart displays last 30 days
- Green line shows income trend
- Red line shows expense trend
- Hovering shows exact values
- X-axis shows dates
- Y-axis shows amounts with $ symbol

---

## Test Scenario 10: Logout and Re-login

### Steps:
1. Click "Logout" button
2. Verify redirect to login page
3. Login again with same credentials

### Expected Result:
- Successfully logged out
- Redirected to login page
- Can login again
- All data persists (income/expense still there)
- Dashboard shows same totals as before logout

---

## Test Scenario 11: Category Breakdown

### Steps:
1. Add various expenses in different categories:
   - Food: $400
   - Rent: $1200
   - Transport: $150
   - Entertainment: $100

2. Check dashboard data (via API or inspect network)

### Expected Result:
- Dashboard API returns category breakdown
- Each category shows correct total
- Categories with $0 not included

---

## Test Scenario 12: Authentication Security

### Steps:
1. Logout
2. Try to access: `http://localhost:8080/dashboard.html` directly
3. Try to access API: `http://localhost:8080/api/dashboard` without token

### Expected Result:
- Accessing dashboard.html without login redirects to login
- API requests without token return 401 Unauthorized

---

## Test Scenario 13: Input Validation

### Steps:
1. Try to add income with:
   - Amount: 0 (should fail)
   - Amount: -100 (should fail)
   - Amount: empty (should fail)
   - No category selected (should fail)
   - No date selected (should fail)

### Expected Result:
- Form validation prevents submission
- Error messages displayed
- Transaction not created

---

## Test Scenario 14: Concurrent Users

### Steps:
1. Open application in two different browsers
2. Register two different users
3. Add transactions in both accounts
4. Verify data isolation

### Expected Result:
- Each user sees only their own transactions
- Totals are independent
- No data leakage between accounts

---

## Test Scenario 15: Mobile Responsiveness

### Steps:
1. Open application in browser
2. Open developer tools (F12)
3. Switch to mobile device view (iPhone, Android)
4. Test all features

### Expected Result:
- Layout adjusts for mobile screen
- All buttons accessible
- Forms usable
- Charts responsive
- No horizontal scrolling

---

## API Testing with curl

### Register and get token:
```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","username":"testuser","password":"test123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email_or_username":"test@test.com","password":"test123"}'
```

### Test protected endpoints:
```bash
# Save token
export TOKEN="your-token-here"

# Get categories
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/categories

# Add income
curl -X POST http://localhost:8080/api/income \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"category_id":1,"amount":5000,"income_date":"2025-12-24","description":"Test"}'

# Get dashboard
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/dashboard
```

---

## Performance Testing

### Steps:
1. Add 100+ transactions
2. Apply various filters
3. Generate PDF report
4. Check response times

### Expected Result:
- Database queries remain fast (< 100ms)
- UI remains responsive
- PDF generates successfully
- No memory leaks

---

## Error Handling Testing

### Test these scenarios:
1. Invalid JSON in API request
2. Missing required fields
3. Invalid date format
4. Expired JWT token
5. Invalid category ID
6. Attempting to edit other user's transaction

### Expected Result:
- Clear error messages
- Appropriate HTTP status codes
- No server crashes
- Errors logged properly

---

## Browser Compatibility

Test in:
- Chrome/Edge (Chromium)
- Firefox
- Safari (if on Mac)

### Expected Result:
- Works in all modern browsers
- Chart.js renders correctly
- CSS styles consistent
- JavaScript functions properly

---

**All tests passing = Production Ready! ✅**
