#!/bin/bash

# AWS EC2 Deployment Script for Income & Expense Tracker
# Run this script on your local machine after setting up EC2 instance

set -e

echo "==================================="
echo "AWS Deployment Script"
echo "==================================="

# Configuration
read -p "Enter your EC2 public IP: " EC2_IP
read -p "Enter path to your .pem key file: " PEM_KEY
read -p "Enter domain name (or press Enter to skip): " DOMAIN

SSH_USER="ubuntu"
APP_DIR="/opt/expense-tracker"

echo ""
echo "Preparing deployment package..."

# Create deployment package
cd "$(dirname "$0")"
rm -f deploy.tar.gz
tar -czf deploy.tar.gz \
    --exclude='*.db' \
    --exclude='*.log' \
    --exclude='.git' \
    --exclude='node_modules' \
    --exclude='deploy.tar.gz' \
    .

echo "Uploading to EC2..."

# Copy deployment package
scp -i "$PEM_KEY" deploy.tar.gz ${SSH_USER}@${EC2_IP}:/tmp/

echo "Installing application on EC2..."

# Execute remote commands
ssh -i "$PEM_KEY" ${SSH_USER}@${EC2_IP} << 'ENDSSH'
# Update system
sudo apt-get update

# Install dependencies if not already installed
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
fi

# Create app directory
sudo mkdir -p /opt/expense-tracker
sudo chown ubuntu:ubuntu /opt/expense-tracker

# Extract application
cd /opt/expense-tracker
tar -xzf /tmp/deploy.tar.gz
rm /tmp/deploy.tar.gz

# Build application
echo "Building application..."
go mod download
go build -o app ./cmd/server

# Create data directory
mkdir -p data

# Create systemd service
echo "Creating systemd service..."
sudo tee /etc/systemd/system/expense-tracker.service > /dev/null << 'EOF'
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
EOF

# Reload systemd and start service
sudo systemctl daemon-reload
sudo systemctl enable expense-tracker
sudo systemctl restart expense-tracker

# Install and configure Nginx
if ! command -v nginx &> /dev/null; then
    echo "Installing Nginx..."
    sudo apt-get install -y nginx
fi

# Wait for service to start
sleep 3

# Check service status
if sudo systemctl is-active --quiet expense-tracker; then
    echo "✅ Application is running!"
else
    echo "❌ Application failed to start. Check logs with: sudo journalctl -u expense-tracker"
    exit 1
fi

ENDSSH

echo ""
echo "==================================="
echo "✅ Deployment Complete!"
echo "==================================="
echo ""
echo "Application is running on EC2!"
echo "Access at: http://${EC2_IP}:8080"
echo ""
echo "Next steps:"
echo "1. Setup Nginx reverse proxy (see AWS_DEPLOYMENT.md)"
echo "2. Configure domain name (if provided)"
echo "3. Setup SSL certificate with Let's Encrypt"
echo ""
echo "Useful commands:"
echo "  Check status: ssh -i $PEM_KEY ubuntu@${EC2_IP} 'sudo systemctl status expense-tracker'"
echo "  View logs:    ssh -i $PEM_KEY ubuntu@${EC2_IP} 'sudo journalctl -u expense-tracker -f'"
echo "  Restart app:  ssh -i $PEM_KEY ubuntu@${EC2_IP} 'sudo systemctl restart expense-tracker'"
echo ""

# Cleanup
rm -f deploy.tar.gz

# Optional: Setup Nginx
read -p "Do you want to setup Nginx reverse proxy now? (y/n): " SETUP_NGINX

if [ "$SETUP_NGINX" = "y" ]; then
    echo "Setting up Nginx..."
    
    NGINX_CONFIG="server {
    listen 80;
    server_name ${DOMAIN:-$EC2_IP};

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}"

    ssh -i "$PEM_KEY" ${SSH_USER}@${EC2_IP} << ENDSSH2
echo "$NGINX_CONFIG" | sudo tee /etc/nginx/sites-available/expense-tracker
sudo ln -sf /etc/nginx/sites-available/expense-tracker /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl restart nginx
ENDSSH2

    echo ""
    echo "✅ Nginx configured!"
    if [ -n "$DOMAIN" ]; then
        echo "Access at: http://${DOMAIN}"
    else
        echo "Access at: http://${EC2_IP}"
    fi
    
    # Optional: Setup SSL
    if [ -n "$DOMAIN" ]; then
        read -p "Do you want to setup SSL certificate? (y/n): " SETUP_SSL
        if [ "$SETUP_SSL" = "y" ]; then
            echo "Setting up SSL with Let's Encrypt..."
            ssh -i "$PEM_KEY" ${SSH_USER}@${EC2_IP} << 'ENDSSH3'
sudo apt-get install -y certbot python3-certbot-nginx
sudo certbot --nginx -d $DOMAIN --non-interactive --agree-tos -m admin@$DOMAIN
ENDSSH3
            echo "✅ SSL certificate installed!"
            echo "Access at: https://${DOMAIN}"
        fi
    fi
fi

echo ""
echo "🎉 All done! Your application is live!"
