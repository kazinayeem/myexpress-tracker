# AWS Deployment Guide

This guide covers deploying your Income & Expense Tracker to AWS using different methods.

---

## Option 1: AWS EC2 (Recommended for Production)

### Prerequisites
- AWS Account
- AWS CLI installed and configured
- SSH key pair for EC2

### Step 1: Launch EC2 Instance

1. **Go to AWS Console** → EC2 → Launch Instance

2. **Choose AMI:**
   - Select **Ubuntu Server 22.04 LTS** (Free tier eligible)

3. **Instance Type:**
   - Select **t2.micro** (Free tier) or **t3.small** for better performance

4. **Configure Instance:**
   - Auto-assign Public IP: **Enable**
   - Add storage: **20 GB** (minimum)

5. **Security Group:**
   Create rules:
   ```
   Type            Protocol    Port Range    Source
   SSH             TCP         22            Your IP / 0.0.0.0/0
   HTTP            TCP         80            0.0.0.0/0
   HTTPS           TCP         443           0.0.0.0/0
   Custom TCP      TCP         8080          0.0.0.0/0 (temporary)
   ```

6. **Launch** and download the `.pem` key file

### Step 2: Connect to EC2

```bash
# Set key permissions
chmod 400 your-key.pem

# Connect to instance
ssh -i your-key.pem ubuntu@your-ec2-public-ip
```

### Step 3: Install Dependencies

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Git
sudo apt install git -y

# Install Nginx (reverse proxy)
sudo apt install nginx -y

# Install SQLite
sudo apt install sqlite3 -y
```

### Step 4: Deploy Application

```bash
# Create app directory
sudo mkdir -p /opt/expense-tracker
sudo chown ubuntu:ubuntu /opt/expense-tracker
cd /opt/expense-tracker

# Clone or upload your code
# Option A: Upload via SCP from local machine
# scp -i your-key.pem -r /path/to/myexpress-tracker/* ubuntu@your-ec2-ip:/opt/expense-tracker/

# Option B: Upload as zip and extract
# On local machine: zip -r app.zip myexpress-tracker/
# scp -i your-key.pem app.zip ubuntu@your-ec2-ip:/opt/expense-tracker/
# On EC2: unzip app.zip && mv myexpress-tracker/* . && rm -rf myexpress-tracker app.zip

# Build the application
cd /opt/expense-tracker
go mod download
go build -o app ./cmd/server

# Create data directory
mkdir -p /opt/expense-tracker/data

# Test the application
./app
# Press Ctrl+C after verifying it starts
```

### Step 5: Create Systemd Service

```bash
sudo nano /etc/systemd/system/expense-tracker.service
```

Add this content:

```ini
[Unit]
Description=Income & Expense Tracker
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/opt/expense-tracker
ExecStart=/opt/expense-tracker/app
Restart=always
RestartSec=5
Environment="PORT=8080"
Environment="ENV=production"

[Install]
WantedBy=multi-user.target
```

Save and enable:

```bash
sudo systemctl daemon-reload
sudo systemctl enable expense-tracker
sudo systemctl start expense-tracker
sudo systemctl status expense-tracker
```

### Step 6: Configure Nginx Reverse Proxy

```bash
sudo nano /etc/nginx/sites-available/expense-tracker
```

Add this configuration:

```nginx
server {
    listen 80;
    server_name your-domain.com;  # or use EC2 public IP

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

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/expense-tracker /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Step 7: Configure Domain (Optional)

1. Go to your domain registrar (GoDaddy, Namecheap, etc.)
2. Add an **A Record**:
   - Type: `A`
   - Name: `@` or `tracker`
   - Value: Your EC2 public IP
   - TTL: 3600

3. Wait for DNS propagation (5-30 minutes)

### Step 8: Setup SSL with Let's Encrypt (Optional but Recommended)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Get SSL certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal is configured automatically
# Test renewal
sudo certbot renew --dry-run
```

### Step 9: Setup Automatic Backups

```bash
# Create backup script
sudo nano /opt/expense-tracker/backup.sh
```

Add this content:

```bash
#!/bin/bash
BACKUP_DIR="/opt/expense-tracker/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_FILE="/opt/expense-tracker/data/tracker.db"

mkdir -p $BACKUP_DIR
cp $DB_FILE $BACKUP_DIR/tracker_$DATE.db

# Keep only last 7 days of backups
find $BACKUP_DIR -name "tracker_*.db" -mtime +7 -delete
```

Make executable and add to cron:

```bash
chmod +x /opt/expense-tracker/backup.sh
crontab -e
```

Add this line (daily backup at 2 AM):

```
0 2 * * * /opt/expense-tracker/backup.sh
```

### Application is now live! 🎉

Access your app at: `http://your-ec2-ip` or `https://your-domain.com`

---

## Option 2: AWS Elastic Beanstalk

### Step 1: Prepare Application

Create `Procfile`:

```bash
web: ./app
```

Create `.ebextensions/01_go.config`:

```yaml
commands:
  01_download_go:
    command: |
      wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
      tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
      export PATH=$PATH:/usr/local/go/bin
```

### Step 2: Deploy

```bash
# Install EB CLI
pip install awsebcli

# Initialize
eb init -p docker expense-tracker --region us-east-1

# Create environment
eb create expense-tracker-env

# Deploy
eb deploy

# Open in browser
eb open
```

---

## Option 3: AWS ECS with Docker

### Step 1: Push Docker Image to ECR

```bash
# Authenticate Docker to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com

# Create ECR repository
aws ecr create-repository --repository-name expense-tracker --region us-east-1

# Build and tag image
docker build -t expense-tracker .
docker tag expense-tracker:latest YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/expense-tracker:latest

# Push to ECR
docker push YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/expense-tracker:latest
```

### Step 2: Create ECS Task Definition

```json
{
  "family": "expense-tracker",
  "containerDefinitions": [
    {
      "name": "expense-tracker",
      "image": "YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/expense-tracker:latest",
      "memory": 512,
      "cpu": 256,
      "essential": true,
      "portMappings": [
        {
          "containerPort": 8080,
          "hostPort": 8080
        }
      ]
    }
  ]
}
```

### Step 3: Create ECS Service

Use AWS Console or CLI to create:
1. ECS Cluster
2. Task Definition (use JSON above)
3. Service with Load Balancer

---

## Option 4: AWS Lightsail (Easiest)

### Step 1: Create Instance

1. Go to **AWS Lightsail** Console
2. Click **Create Instance**
3. Select **Linux/Unix** → **OS Only** → **Ubuntu 22.04**
4. Choose plan ($3.50/month minimum)
5. Name your instance and create

### Step 2: Configure and Deploy

Same as EC2 steps 2-9 above, but using Lightsail's browser SSH.

---

## Cost Estimates

### EC2 Option:
- t2.micro (free tier): **$0/month** (first year)
- t3.small: **~$15/month**
- Domain: **~$12/year**
- SSL: **Free** (Let's Encrypt)

### Lightsail:
- $3.50/month to $20/month (includes everything)

### ECS:
- Fargate: **~$30/month** (0.25 vCPU, 0.5 GB)

---

## Monitoring and Maintenance

### Check Application Status

```bash
# Check service status
sudo systemctl status expense-tracker

# View logs
sudo journalctl -u expense-tracker -f

# Restart service
sudo systemctl restart expense-tracker
```

### Database Management

```bash
# Backup database
cp /opt/expense-tracker/data/tracker.db ~/backup-$(date +%Y%m%d).db

# Check database
sqlite3 /opt/expense-tracker/data/tracker.db "SELECT * FROM users;"
```

### Update Application

```bash
# Upload new build
scp -i your-key.pem app ubuntu@your-ec2-ip:/opt/expense-tracker/

# Restart service
sudo systemctl restart expense-tracker
```

---

## Troubleshooting

### Application won't start:
```bash
# Check logs
sudo journalctl -u expense-tracker -n 50

# Check port availability
sudo netstat -tulpn | grep 8080
```

### Database locked:
```bash
# Stop application
sudo systemctl stop expense-tracker

# Check database integrity
sqlite3 /opt/expense-tracker/data/tracker.db "PRAGMA integrity_check;"

# Restart
sudo systemctl start expense-tracker
```

### Nginx errors:
```bash
# Check configuration
sudo nginx -t

# View error logs
sudo tail -f /var/log/nginx/error.log
```

---

## Security Best Practices

1. **Firewall:** Only allow necessary ports
2. **SSH:** Use key authentication only, disable password login
3. **SSL:** Always use HTTPS in production
4. **Backups:** Automate daily backups
5. **Updates:** Keep system packages updated
6. **Monitoring:** Setup CloudWatch alarms
7. **Database:** Encrypt sensitive data
8. **Environment Variables:** Use for secrets, not hardcoded

---

## Need Help?

- AWS Documentation: https://docs.aws.amazon.com/
- AWS Free Tier: https://aws.amazon.com/free/
- AWS Support: https://console.aws.amazon.com/support/

Your application is production-ready! 🚀
