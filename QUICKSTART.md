# Quick Start Guide

## 🚀 Start the Application (3 Methods)

### Method 1: Using Go directly
```bash
go run cmd/server/main.go
```

### Method 2: Build and Run
**Windows:**
```bash
run.bat
```

**Linux/Mac:**
```bash
chmod +x run.sh
./run.sh
```

### Method 3: Docker
```bash
docker-compose up -d
```

## 📱 Access the Application

Open your browser and go to: **http://localhost:8080**

## 🔐 First Steps

1. **Register**: Click "Register here" on the login page
2. **Create Account**: Enter email, username, and password (min 6 characters)
3. **Login**: Use your credentials to access the dashboard
4. **Start Tracking**: Add your first income or expense!

## 💡 Quick Tips

### Default Categories

**Income:**
- Salary
- Freelance
- Investment
- Other Income

**Expense:**
- Food
- Transport
- Rent
- Utilities
- Entertainment
- Healthcare
- Shopping
- Other Expense

### Adding Transactions
1. Click "Add Income" or "Add Expense" button
2. Select a category
3. Enter amount (must be greater than 0)
4. Choose a date
5. Add optional description
6. Click "Save"

### Filtering Transactions
- **By Type**: Select Income or Expense
- **By Category**: Choose from dropdown
- **By Date**: Pick a specific date
- **Clear All**: Click "Clear Filters"

### Viewing Reports
- Dashboard shows summary cards with totals
- Chart displays last 30 days of data
- Export PDF for custom date ranges

### Editing/Deleting
- Click the ✏️ icon to edit a transaction
- Click the 🗑️ icon to delete (confirmation required)

## 🛠️ Troubleshooting

### Port Already in Use
If port 8080 is busy, change it in `.env`:
```
SERVER_PORT=3000
```

### Database Issues
Delete the database and restart:
```bash
rm -rf data/tracker.db
go run cmd/server/main.go
```

### Docker Issues
```bash
docker-compose down
docker-compose up -d --build
```

### Can't Login
- Ensure you registered first
- Check that password is at least 6 characters
- Clear browser cache/cookies

## 📊 Sample Data for Testing

After registration, try adding:

**Income:**
```
Salary - $5,000 - Monthly salary
Freelance - $1,500 - Web design project
```

**Expenses:**
```
Rent - $1,200 - Monthly apartment rent
Food - $400 - Weekly groceries
Transport - $150 - Gas and public transit
Utilities - $200 - Electric, water, internet
```

## 🔒 Security Notes

- Never share your JWT token
- Use strong passwords
- In production, always use HTTPS
- Change JWT_SECRET in production

## 📚 Additional Resources

- Full API documentation: See `API_TESTING.md`
- Deployment guide: See `README.md` AWS section
- Docker deployment: See `docker-compose.yml`

## 🆘 Need Help?

- Check the logs: `docker-compose logs -f` (Docker)
- View terminal output (Direct run)
- Open an issue on GitHub

---

**Happy Tracking! 💰📊**
