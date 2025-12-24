# Income & Expense Tracker

A production-ready full-stack Income & Expense Tracker built with **Golang** backend, **HTML/CSS/JavaScript** frontend, **SQLite** database, and **Docker** support.

## 🎯 Features

- **User Authentication**: Secure registration and login with JWT tokens and bcrypt password hashing
- **Income & Expense Management**: Full CRUD operations for tracking income and expenses
- **Category Support**: Pre-defined categories (Food, Transport, Rent, Salary, etc.)
- **Date Filtering**: Filter transactions by date, date ranges, or view today's transactions
- **Dashboard**: Visual overview with total income, expenses, balance, and daily summaries
- **Charts**: Interactive Chart.js visualizations showing income vs expense trends
- **PDF Export**: Generate and download PDF reports for any date range
- **Responsive Design**: Mobile-friendly UI built with pure HTML/CSS/JavaScript

## 🧩 Tech Stack

- **Backend**: Go 1.21+ with net/http
- **Frontend**: Pure HTML5, CSS3, Vanilla JavaScript
- **Database**: SQLite with proper indexing and foreign keys
- **Authentication**: JWT tokens with bcrypt password hashing
- **Charts**: Chart.js v4.4.0
- **PDF Generation**: gofpdf library
- **Deployment**: Docker & Docker Compose

## 📁 Project Structure

```
myexpress-tracker/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/
│   │   └── auth.go              # JWT & bcrypt authentication
│   ├── database/
│   │   ├── sqlite.go            # Database connection
│   │   └── migrations.go        # Schema & migrations
│   ├── handlers/
│   │   ├── auth.go              # Auth endpoints
│   │   ├── income.go            # Income CRUD
│   │   ├── expense.go           # Expense CRUD
│   │   ├── category.go          # Category endpoints
│   │   ├── dashboard.go         # Dashboard data
│   │   └── export.go            # PDF export
│   ├── middleware/
│   │   └── auth.go              # JWT middleware
│   ├── models/
│   │   └── models.go            # Data models
│   └── repository/
│       ├── user.go              # User repository
│       ├── category.go          # Category repository
│       ├── income.go            # Income repository
│       └── expense.go           # Expense repository
├── configs/
│   └── config.go                # Configuration management
├── web/
│   ├── index.html               # Home redirect
│   ├── login.html               # Login page
│   ├── register.html            # Registration page
│   ├── dashboard.html           # Main dashboard
│   └── static/
│       ├── css/
│       │   └── style.css        # Application styles
│       └── js/
│           ├── auth.js          # Authentication logic
│           └── dashboard.js     # Dashboard functionality
├── data/                        # SQLite database (gitignored)
├── Dockerfile                   # Multi-stage Docker build
├── docker-compose.yml           # Docker Compose configuration
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
└── README.md                    # This file
```

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Categories Table
```sql
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Income Table
```sql
CREATE TABLE income (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    amount REAL NOT NULL CHECK(amount > 0),
    description TEXT,
    income_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);
```

### Expense Table
```sql
CREATE TABLE expense (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    amount REAL NOT NULL CHECK(amount > 0),
    description TEXT,
    expense_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose (for containerized deployment)
- Git

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd myexpress-tracker
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

5. **Access the application**
   - Open browser: `http://localhost:8080`
   - Default page redirects to login
   - Register a new account to get started

### Using Docker

1. **Build and run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

2. **View logs**
   ```bash
   docker-compose logs -f
   ```

3. **Stop the application**
   ```bash
   docker-compose down
   ```

4. **Rebuild after code changes**
   ```bash
   docker-compose up -d --build
   ```

### Using Docker without Compose

1. **Build the Docker image**
   ```bash
   docker build -t expense-tracker:latest .
   ```

2. **Run the container**
   ```bash
   docker run -d \
     -p 8080:8080 \
     -v $(pwd)/data:/app/data \
     -e JWT_SECRET=your-secret-key \
     --name expense-tracker \
     expense-tracker:latest
   ```

## 📡 API Endpoints

### Authentication

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "john_doe",
  "password": "securepassword"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email_or_username": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "john_doe"
  },
  "message": "Login successful"
}
```

### Protected Endpoints (Require JWT Token)

**Authorization Header:**
```
Authorization: Bearer <your-jwt-token>
```

#### Get Categories
```http
GET /api/categories?type=income
```

#### Create Income
```http
POST /api/income
Content-Type: application/json

{
  "category_id": 1,
  "amount": 5000.00,
  "description": "Monthly salary",
  "income_date": "2025-01-15"
}
```

#### Get Incomes
```http
GET /api/income?date=2025-01-15
GET /api/income?start_date=2025-01-01&end_date=2025-01-31
GET /api/income?category_id=1
```

#### Update Income
```http
PUT /api/income/{id}
Content-Type: application/json

{
  "category_id": 1,
  "amount": 5500.00,
  "description": "Updated salary",
  "income_date": "2025-01-15"
}
```

#### Delete Income
```http
DELETE /api/income/{id}
```

#### Create Expense
```http
POST /api/expense
Content-Type: application/json

{
  "category_id": 5,
  "amount": 150.00,
  "description": "Groceries",
  "expense_date": "2025-01-15"
}
```

#### Get Expenses
```http
GET /api/expense?date=2025-01-15
GET /api/expense?start_date=2025-01-01&end_date=2025-01-31
```

#### Update Expense
```http
PUT /api/expense/{id}
```

#### Delete Expense
```http
DELETE /api/expense/{id}
```

#### Get Dashboard Summary
```http
GET /api/dashboard
```

**Response:**
```json
{
  "total_income": 10000.00,
  "total_expense": 5000.00,
  "balance": 5000.00,
  "today_income": 500.00,
  "today_expense": 200.00,
  "monthly_income": 8000.00,
  "monthly_expense": 4000.00,
  "daily_data": [
    {"date": "2025-01-01", "income": 500, "expense": 200},
    ...
  ],
  "category_breakdown": {
    "income_by_category": {"Salary": 5000, "Freelance": 3000},
    "expense_by_category": {"Food": 1500, "Rent": 2000}
  }
}
```

#### Export to PDF
```http
GET /api/export/pdf?start_date=2025-01-01&end_date=2025-01-31
```

## 🐳 Docker Deployment

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Port to run the server | `8080` |
| `DATABASE_PATH` | Path to SQLite database | `./data/tracker.db` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-secret-key-change-in-production` |
| `JWT_EXPIRATION_HOURS` | JWT token expiration time | `24` |
| `ENVIRONMENT` | Application environment | `development` |

### Production Deployment

1. Update `JWT_SECRET` in docker-compose.yml or use environment variables
2. Consider using a reverse proxy (Nginx) for SSL/TLS
3. Set up regular database backups
4. Monitor application logs

## ☁️ AWS EC2 Deployment Guide

### Prerequisites
- AWS Account
- EC2 instance with Ubuntu 22.04 LTS
- Security group allowing ports 22 (SSH), 80 (HTTP), 443 (HTTPS)

### Step 1: Connect to EC2 Instance
```bash
ssh -i your-key.pem ubuntu@your-ec2-public-ip
```

### Step 2: Install Docker
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker ubuntu

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version
```

### Step 3: Deploy Application
```bash
# Clone repository
git clone <your-repo-url>
cd myexpress-tracker

# Create .env file
cp .env.example .env
nano .env  # Update JWT_SECRET and other configs

# Build and run
docker-compose up -d

# Check status
docker-compose ps
docker-compose logs -f
```

### Step 4: Set Up Nginx Reverse Proxy (Optional)

```bash
# Install Nginx
sudo apt install nginx -y

# Create Nginx configuration
sudo nano /etc/nginx/sites-available/expense-tracker
```

Add the following configuration:
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/expense-tracker /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Step 5: Set Up SSL with Let's Encrypt (Optional)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain SSL certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal is configured by default
```

### Step 6: Set Up Automatic Backups

```bash
# Create backup script
mkdir -p ~/backups
nano ~/backup-db.sh
```

Add the following:
```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="$HOME/backups"
cp ~/myexpress-tracker/data/tracker.db "$BACKUP_DIR/tracker_$DATE.db"
# Keep only last 30 days of backups
find "$BACKUP_DIR" -name "tracker_*.db" -mtime +30 -delete
```

```bash
# Make executable
chmod +x ~/backup-db.sh

# Add to crontab (daily at 2 AM)
crontab -e
# Add line:
0 2 * * * /home/ubuntu/backup-db.sh
```

### Step 7: Monitor Application

```bash
# View logs
docker-compose logs -f

# Check container status
docker-compose ps

# Restart application
docker-compose restart

# Update application
cd ~/myexpress-tracker
git pull
docker-compose up -d --build
```

## 🔒 Security Best Practices

1. **JWT Secret**: Use a strong, random secret key in production
2. **HTTPS**: Always use SSL/TLS in production (Let's Encrypt)
3. **Database Backups**: Implement regular automated backups
4. **Firewall**: Configure firewall to allow only necessary ports
5. **Updates**: Keep system and dependencies up to date
6. **Rate Limiting**: Consider adding rate limiting for API endpoints
7. **Input Validation**: All inputs are validated on the backend

## 🧪 Sample Test Data

After registration, you can add sample data:

**Income Examples:**
- Salary: $5000 (monthly)
- Freelance Project: $1500
- Investment Returns: $200

**Expense Examples:**
- Rent: $1200
- Groceries: $400
- Transportation: $150
- Utilities: $200
- Entertainment: $100

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Support

For issues and questions, please open an issue on GitHub.

---

**Built with ❤️ using Go, HTML/CSS/JS, and SQLite**
