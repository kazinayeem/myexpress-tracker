// API Configuration
const API_BASE = window.location.origin + '/api';

// Theme initialization
function initializeTheme() {
    const savedTheme = localStorage.getItem('theme');
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const theme = savedTheme || (systemPrefersDark ? 'dark' : 'light');
    
    if (theme === 'dark') {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }
    
    if (!savedTheme) {
        localStorage.setItem('theme', theme);
    }
}

// Initialize theme on page load
initializeTheme();

// Helper function to get token
function getToken() {
    return localStorage.getItem('token');
}

// Helper function to set token
function setToken(token) {
    localStorage.setItem('token', token);
}

// Helper function to remove token
function removeToken() {
    localStorage.removeItem('token');
}

// Helper function to show error
function showError(elementId, message) {
    const errorElement = document.getElementById(elementId);
    if (errorElement) {
        errorElement.textContent = message;
        errorElement.classList.remove('hidden');
    }
}

// Helper function to show success
function showSuccess(elementId, message) {
    const successElement = document.getElementById(elementId);
    if (successElement) {
        successElement.textContent = message;
        successElement.classList.remove('hidden');
    }
}

// Login Form Handler
const loginForm = document.getElementById('login-form');
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const emailOrUsername = document.getElementById('email_or_username').value;
        const password = document.getElementById('password').value;
        
        try {
            const response = await fetch(`${API_BASE}/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email_or_username: emailOrUsername,
                    password: password
                })
            });
            
            const data = await response.json();
            
            if (response.ok) {
                setToken(data.token);
                localStorage.setItem('username', data.user.username);
                window.location.href = '/dashboard.html';
            } else {
                showError('error-message', data.error || 'Login failed');
            }
        } catch (error) {
            showError('error-message', 'Network error. Please try again.');
        }
    });
}

// Register Form Handler
const registerForm = document.getElementById('register-form');
if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const email = document.getElementById('email').value;
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        
        try {
            const response = await fetch(`${API_BASE}/auth/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: email,
                    username: username,
                    password: password
                })
            });
            
            const data = await response.json();
            
            if (response.ok) {
                showSuccess('success-message', 'Registration successful! Redirecting to dashboard...');
                setToken(data.token);
                localStorage.setItem('username', data.user.username);
                setTimeout(() => {
                    window.location.href = '/dashboard.html';
                }, 1500);
            } else {
                showError('error-message', data.error || 'Registration failed');
            }
        } catch (error) {
            showError('error-message', 'Network error. Please try again.');
        }
    });
}

// Logout Function
function logout() {
    removeToken();
    localStorage.removeItem('username');
    localStorage.removeItem('currency');
    // Keep theme preference
    window.location.href = '/login.html';
}
