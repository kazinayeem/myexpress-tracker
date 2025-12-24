// Use API_BASE from auth.js (already declared there)
// No need to redeclare

// Global variables
let categoryChart = null;
let trendChart = null;
let currentCurrency = 'USD';
let currencySymbol = '$';

// Currency symbols mapping
const currencySymbols = {
    'USD': '$',
    'EUR': '€',
    'GBP': '£',
    'JPY': '¥',
    'CNY': '¥',
    'INR': '₹',
    'BDT': '৳',
    'AUD': 'A$',
    'CAD': 'C$',
    'CHF': 'CHF',
    'BRL': 'R$',
    'KRW': '₩',
    'MXN': 'MX$',
    'SGD': 'S$',
    'PKR': '₨',
    'SAR': 'SR',
    'AED': 'د.إ',
    'THB': '฿',
    'MYR': 'RM',
    'IDR': 'Rp'
};

// Check authentication
function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login.html';
        return false;
    }
    return true;
}

// Theme Functions
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

function toggleTheme() {
    const isDark = document.documentElement.classList.contains('dark');
    const newTheme = isDark ? 'light' : 'dark';
    
    if (newTheme === 'dark') {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }
    
    localStorage.setItem('theme', newTheme);
    
    // Update theme on server
    updateUserTheme(newTheme);
    
    // Reload charts with new theme
    loadDashboardData();
}

function setTheme(theme) {
    if (theme === 'dark') {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }
    
    localStorage.setItem('theme', theme);
}

// Mobile Menu
function toggleMobileMenu() {
    const sidebar = document.getElementById('sidebar');
    const overlay = document.getElementById('mobileMenuOverlay');
    
    sidebar.classList.toggle('-translate-x-full');
    overlay.classList.toggle('hidden');
}

// Section Navigation
function showSection(section, event) {
    console.log('showSection called:', section);
    
    const sections = ['dashboardSection', 'incomeSection', 'expenseSection', 'reportsSection'];
    const titles = {
        'dashboard': 'Dashboard',
        'income': 'Income Management',
        'expense': 'Expense Management',
        'reports': 'Reports & Analytics'
    };
    
    sections.forEach(s => {
        const element = document.getElementById(s);
        if (element) {
            element.classList.add('hidden');
        }
    });
    
    const sectionElement = document.getElementById(section + 'Section');
    if (sectionElement) {
        sectionElement.classList.remove('hidden');
        console.log('Section shown:', section);
    } else {
        console.error('Section element not found:', section + 'Section');
    }
    
    const titleElement = document.getElementById('sectionTitle');
    if (titleElement) {
        titleElement.textContent = titles[section] || 'Dashboard';
    }
    
    // Update nav items
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('bg-blue-50', 'dark:bg-gray-700', 'text-blue-600', 'dark:text-blue-400');
        item.classList.add('text-gray-700', 'dark:text-gray-300');
    });
    
    if (event && event.target) {
        const navItem = event.target.closest('.nav-item');
        if (navItem) {
            navItem.classList.add('bg-blue-50', 'dark:bg-gray-700', 'text-blue-600', 'dark:text-blue-400');
            navItem.classList.remove('text-gray-700', 'dark:text-gray-300');
        }
    }
    
    // Close mobile menu
    if (window.innerWidth < 768) {
        toggleMobileMenu();
    }
}

// Settings Modal
function showSettings() {
    document.getElementById('settingsModal').classList.remove('hidden');
    document.getElementById('currencySelect').value = currentCurrency;
}

function closeSettings() {
    document.getElementById('settingsModal').classList.add('hidden');
}

async function saveSettings() {
    const newCurrency = document.getElementById('currencySelect').value;
    const theme = document.documentElement.classList.contains('dark') ? 'dark' : 'light';
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE}/user/settings`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                currency: newCurrency,
                theme: theme
            })
        });
        
        if (response.ok) {
            currentCurrency = newCurrency;
            currencySymbol = currencySymbols[newCurrency];
            localStorage.setItem('currency', newCurrency);
            
            updateCurrencyDisplay();
            loadDashboardData();
            closeSettings();
            
            // Show success message
            showNotification('Settings saved successfully!', 'success');
        } else {
            showNotification('Failed to save settings', 'error');
        }
    } catch (error) {
        console.error('Error saving settings:', error);
        showNotification('Network error', 'error');
    }
}

async function updateUserTheme(theme) {
    try {
        const token = localStorage.getItem('token');
        await fetch(`${API_BASE}/user/settings`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ theme })
        });
    } catch (error) {
        console.error('Error updating theme:', error);
    }
}

// Currency Functions
function updateCurrencyDisplay() {
    // Update all currency symbols
    document.querySelectorAll('[id*="currencySymbol"], [id*="Currency"]').forEach(el => {
        el.textContent = currencySymbol;
    });
    
    document.getElementById('currentCurrency').textContent = currentCurrency;
}

// API Helper
async function apiRequest(endpoint, options = {}) {
    const token = localStorage.getItem('token');
    
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    };
    
    const mergedOptions = {
        ...defaultOptions,
        ...options,
        headers: {
            ...defaultOptions.headers,
            ...options.headers
        }
    };
    
    const response = await fetch(`${API_BASE}${endpoint}`, mergedOptions);
    
    if (response.status === 401) {
        localStorage.removeItem('token');
        window.location.href = '/login.html';
        throw new Error('Unauthorized');
    }
    
    return response;
}

// Load User Profile
async function loadUserProfile() {
    try {
        console.log('Loading user profile...');
        const response = await apiRequest('/user/profile');
        
        if (!response.ok) {
            console.error('Profile API failed:', response.status);
            // Set defaults if API fails
            const username = localStorage.getItem('username') || 'User';
            document.getElementById('usernameDisplay').textContent = username;
            document.getElementById('emailDisplay').textContent = 'user@example.com';
            document.getElementById('userInitial').textContent = username.charAt(0).toUpperCase();
            return;
        }
        
        const user = await response.json();
        console.log('User profile loaded:', user);
        
        // Update user info
        document.getElementById('usernameDisplay').textContent = user.username;
        document.getElementById('emailDisplay').textContent = user.email;
        document.getElementById('userInitial').textContent = user.username.charAt(0).toUpperCase();
        
        // Update currency
        currentCurrency = user.currency || 'USD';
        currencySymbol = currencySymbols[currentCurrency];
        localStorage.setItem('currency', currentCurrency);
        updateCurrencyDisplay();
        
        // Update theme
        if (user.theme) {
            setTheme(user.theme);
        }
    } catch (error) {
        console.error('Failed to load user profile:', error);
        // Set defaults on error
        const username = localStorage.getItem('username') || 'User';
        document.getElementById('usernameDisplay').textContent = username;
        document.getElementById('emailDisplay').textContent = '';
        document.getElementById('userInitial').textContent = username.charAt(0).toUpperCase();
    }
}

// Load Dashboard Data
async function loadDashboardData() {
    try {
        const response = await apiRequest('/dashboard');
        const data = await response.json();
        
        // Update summary cards
        document.getElementById('totalIncome').textContent = data.total_income.toFixed(2);
        document.getElementById('totalExpense').textContent = data.total_expense.toFixed(2);
        document.getElementById('balance').textContent = data.balance.toFixed(2);
        
        // Load charts
        if (data.daily_data && data.daily_data.length > 0) {
            loadTrendChart(data.daily_data);
        }
        
        // Load category breakdown
        await loadCategoryChart();
        
        // Load recent transactions
        await loadRecentTransactions();
    } catch (error) {
        console.error('Failed to load dashboard:', error);
    }
}

// Load Trend Chart
function loadTrendChart(dailyData) {
    const ctx = document.getElementById('trendChart');
    if (!ctx) return;
    
    if (trendChart) {
        trendChart.destroy();
    }
    
    const isDark = document.documentElement.classList.contains('dark');
    const textColor = isDark ? '#e5e7eb' : '#374151';
    const gridColor = isDark ? '#374151' : '#e5e7eb';
    
    const labels = dailyData.map(d => d.date);
    const incomeData = dailyData.map(d => d.income);
    const expenseData = dailyData.map(d => d.expense);
    
    trendChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Income',
                    data: incomeData,
                    borderColor: '#10b981',
                    backgroundColor: 'rgba(16, 185, 129, 0.1)',
                    tension: 0.4,
                    fill: true
                },
                {
                    label: 'Expense',
                    data: expenseData,
                    borderColor: '#ef4444',
                    backgroundColor: 'rgba(239, 68, 68, 0.1)',
                    tension: 0.4,
                    fill: true
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'top',
                    labels: {
                        color: textColor,
                        font: { size: 12 }
                    }
                }
            },
            scales: {
                x: {
                    ticks: { color: textColor },
                    grid: { color: gridColor }
                },
                y: {
                    beginAtZero: true,
                    ticks: {
                        color: textColor,
                        callback: function(value) {
                            return currencySymbol + value.toFixed(0);
                        }
                    },
                    grid: { color: gridColor }
                }
            }
        }
    });
}

// Load Category Chart
async function loadCategoryChart() {
    try {
        const response = await apiRequest('/expense');
        const expenses = await response.json();
        
        const ctx = document.getElementById('categoryChart');
        if (!ctx) return;
        
        if (categoryChart) {
            categoryChart.destroy();
        }
        
        // Group expenses by category
        const categoryTotals = {};
        expenses.forEach(exp => {
            if (!categoryTotals[exp.category_name]) {
                categoryTotals[exp.category_name] = 0;
            }
            categoryTotals[exp.category_name] += exp.amount;
        });
        
        const labels = Object.keys(categoryTotals);
        const data = Object.values(categoryTotals);
        
        if (labels.length === 0) {
            ctx.getContext('2d').fillText('No data available', 10, 50);
            return;
        }
        
        const isDark = document.documentElement.classList.contains('dark');
        const textColor = isDark ? '#e5e7eb' : '#374151';
        
        const colors = [
            '#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6',
            '#ec4899', '#14b8a6', '#f97316', '#6366f1', '#84cc16'
        ];
        
        categoryChart = new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: labels,
                datasets: [{
                    data: data,
                    backgroundColor: colors,
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: {
                            color: textColor,
                            font: { size: 11 },
                            padding: 10
                        }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                return `${label}: ${currencySymbol}${value.toFixed(2)}`;
                            }
                        }
                    }
                }
            }
        });
    } catch (error) {
        console.error('Failed to load category chart:', error);
    }
}

// Load Recent Transactions
async function loadRecentTransactions() {
    try {
        const [incomeRes, expenseRes] = await Promise.all([
            apiRequest('/income'),
            apiRequest('/expense')
        ]);
        
        const income = await incomeRes.json();
        const expenses = await expenseRes.json();
        
        // Combine and sort
        const all = [
            ...income.map(i => ({ ...i, type: 'income' })),
            ...expenses.map(e => ({ ...e, type: 'expense' }))
        ].sort((a, b) => new Date(b.date) - new Date(a.date)).slice(0, 10);
        
        const container = document.getElementById('recentTransactions');
        
        if (all.length === 0) {
            container.innerHTML = '<p class="text-gray-500 dark:text-gray-400 text-center py-4">No transactions yet</p>';
            return;
        }
        
        container.innerHTML = all.map(t => `
            <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg hover:shadow-md transition">
                <div class="flex items-center space-x-3">
                    <div class="w-10 h-10 rounded-full flex items-center justify-center ${
                        t.type === 'income' ? 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400' : 
                        'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400'
                    }">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                                d="${t.type === 'income' ? 'M7 11l5-5m0 0l5 5m-5-5v12' : 'M17 13l-5 5m0 0l-5-5m5 5V6'}"></path>
                        </svg>
                    </div>
                    <div>
                        <p class="font-medium text-gray-800 dark:text-white">${t.description}</p>
                        <p class="text-xs text-gray-500 dark:text-gray-400">${t.category_name} • ${t.date}</p>
                    </div>
                </div>
                <div class="text-right">
                    <p class="font-bold ${t.type === 'income' ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}">
                        ${t.type === 'income' ? '+' : '-'}${currencySymbol}${t.amount.toFixed(2)}
                    </p>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load recent transactions:', error);
    }
}

// Load Categories
async function loadCategories() {
    try {
        const response = await apiRequest('/categories');
        const categories = await response.json();
        
        const incomeCategories = categories.filter(c => c.type === 'income');
        const expenseCategories = categories.filter(c => c.type === 'expense');
        
        // Populate income category select
        const incomeSelect = document.getElementById('incomeCategory');
        if (incomeSelect) {
            incomeSelect.innerHTML = '<option value="">Select category</option>' +
                incomeCategories.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
        }
        
        // Populate expense category select
        const expenseSelect = document.getElementById('expenseCategory');
        if (expenseSelect) {
            expenseSelect.innerHTML = '<option value="">Select category</option>' +
                expenseCategories.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
        }
    } catch (error) {
        console.error('Failed to load categories:', error);
    }
}

// Income Form Handler
const incomeForm = document.getElementById('incomeForm');
if (incomeForm) {
    incomeForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const data = {
            amount: parseFloat(document.getElementById('incomeAmount').value),
            category_id: parseInt(document.getElementById('incomeCategory').value),
            description: document.getElementById('incomeDescription').value,
            income_date: document.getElementById('incomeDate').value
        };
        
        try {
            const response = await apiRequest('/income', {
                method: 'POST',
                body: JSON.stringify(data)
            });
            
            if (response.ok) {
                showNotification('Income added successfully!', 'success');
                incomeForm.reset();
                document.getElementById('incomeDate').value = new Date().toISOString().split('T')[0];
                await loadDashboardData();
                await loadIncomeList();
            } else {
                const error = await response.json();
                showNotification(error.error || 'Failed to add income', 'error');
            }
        } catch (error) {
            showNotification('Network error', 'error');
        }
    });
}

// Expense Form Handler
const expenseForm = document.getElementById('expenseForm');
if (expenseForm) {
    expenseForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const data = {
            amount: parseFloat(document.getElementById('expenseAmount').value),
            category_id: parseInt(document.getElementById('expenseCategory').value),
            description: document.getElementById('expenseDescription').value,
            expense_date: document.getElementById('expenseDate').value
        };
        
        try {
            const response = await apiRequest('/expense', {
                method: 'POST',
                body: JSON.stringify(data)
            });
            
            if (response.ok) {
                showNotification('Expense added successfully!', 'success');
                expenseForm.reset();
                document.getElementById('expenseDate').value = new Date().toISOString().split('T')[0];
                await loadDashboardData();
                await loadExpenseList();
            } else {
                const error = await response.json();
                showNotification(error.error || 'Failed to add expense', 'error');
            }
        } catch (error) {
            showNotification('Network error', 'error');
        }
    });
}

// Load Income List
async function loadIncomeList() {
    try {
        const response = await apiRequest('/income');
        const income = await response.json();
        
        const container = document.getElementById('incomeList');
        if (!container) return;
        
        if (income.length === 0) {
            container.innerHTML = '<p class="text-gray-500 dark:text-gray-400 text-center py-4">No income records yet</p>';
            return;
        }
        
        container.innerHTML = income.map(i => `
            <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <div class="flex-1">
                    <p class="font-medium text-gray-800 dark:text-white">${i.description}</p>
                    <p class="text-sm text-gray-500 dark:text-gray-400">${i.category_name} • ${i.date}</p>
                </div>
                <div class="text-right">
                    <p class="font-bold text-green-600 dark:text-green-400">${currencySymbol}${i.amount.toFixed(2)}</p>
                    <button onclick="deleteIncome(${i.id})" class="text-xs text-red-500 hover:text-red-700">Delete</button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load income list:', error);
    }
}

// Load Expense List
async function loadExpenseList() {
    try {
        const response = await apiRequest('/expense');
        const expenses = await response.json();
        
        const container = document.getElementById('expenseList');
        if (!container) return;
        
        if (expenses.length === 0) {
            container.innerHTML = '<p class="text-gray-500 dark:text-gray-400 text-center py-4">No expense records yet</p>';
            return;
        }
        
        container.innerHTML = expenses.map(e => `
            <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <div class="flex-1">
                    <p class="font-medium text-gray-800 dark:text-white">${e.description}</p>
                    <p class="text-sm text-gray-500 dark:text-gray-400">${e.category_name} • ${e.date}</p>
                </div>
                <div class="text-right">
                    <p class="font-bold text-red-600 dark:text-red-400">${currencySymbol}${e.amount.toFixed(2)}</p>
                    <button onclick="deleteExpense(${e.id})" class="text-xs text-red-500 hover:text-red-700">Delete</button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load expense list:', error);
    }
}

// Delete Functions
async function deleteIncome(id) {
    if (!confirm('Are you sure you want to delete this income?')) return;
    
    try {
        const response = await apiRequest(`/income/${id}`, { method: 'DELETE' });
        if (response.ok) {
            showNotification('Income deleted', 'success');
            await loadDashboardData();
            await loadIncomeList();
        }
    } catch (error) {
        showNotification('Failed to delete income', 'error');
    }
}

async function deleteExpense(id) {
    if (!confirm('Are you sure you want to delete this expense?')) return;
    
    try {
        const response = await apiRequest(`/expense/${id}`, { method: 'DELETE' });
        if (response.ok) {
            showNotification('Expense deleted', 'success');
            await loadDashboardData();
            await loadExpenseList();
        }
    } catch (error) {
        showNotification('Failed to delete expense', 'error');
    }
}

// Export PDF
async function exportPDF() {
    try {
        const response = await apiRequest('/export/pdf', { method: 'GET' });
        const blob = await response.blob();
        
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `expense-report-${new Date().toISOString().split('T')[0]}.pdf`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
        
        showNotification('PDF exported successfully!', 'success');
    } catch (error) {
        showNotification('Failed to export PDF', 'error');
    }
}

// Load Reports
async function loadReports() {
    const startDate = document.getElementById('reportStartDate').value;
    const endDate = document.getElementById('reportEndDate').value;
    
    if (!startDate || !endDate) {
        showNotification('Please select start and end dates', 'error');
        return;
    }
    
    try {
        const [incomeRes, expenseRes] = await Promise.all([
            apiRequest(`/income?start_date=${startDate}&end_date=${endDate}`),
            apiRequest(`/expense?start_date=${startDate}&end_date=${endDate}`)
        ]);
        
        const income = await incomeRes.json();
        const expenses = await expenseRes.json();
        
        // Income summary
        const incomeTotal = income.reduce((sum, i) => sum + i.amount, 0);
        const incomeByCategory = {};
        income.forEach(i => {
            if (!incomeByCategory[i.category_name]) {
                incomeByCategory[i.category_name] = 0;
            }
            incomeByCategory[i.category_name] += i.amount;
        });
        
        // Expense summary
        const expenseTotal = expenses.reduce((sum, e) => sum + e.amount, 0);
        const expenseByCategory = {};
        expenses.forEach(e => {
            if (!expenseByCategory[e.category_name]) {
                expenseByCategory[e.category_name] = 0;
            }
            expenseByCategory[e.category_name] += e.amount;
        });
        
        // Display income report
        const incomeReportEl = document.getElementById('incomeReport');
        incomeReportEl.innerHTML = `
            <div class="mb-4 p-4 bg-green-50 dark:bg-green-900/20 rounded-lg">
                <p class="text-sm text-gray-600 dark:text-gray-400">Total Income</p>
                <p class="text-2xl font-bold text-green-600 dark:text-green-400">${currencySymbol}${incomeTotal.toFixed(2)}</p>
            </div>
            ${Object.entries(incomeByCategory).map(([cat, amount]) => `
                <div class="flex justify-between p-2 border-b border-gray-200 dark:border-gray-700">
                    <span class="text-gray-700 dark:text-gray-300">${cat}</span>
                    <span class="font-semibold text-gray-800 dark:text-white">${currencySymbol}${amount.toFixed(2)}</span>
                </div>
            `).join('')}
        `;
        
        // Display expense report
        const expenseReportEl = document.getElementById('expenseReport');
        expenseReportEl.innerHTML = `
            <div class="mb-4 p-4 bg-red-50 dark:bg-red-900/20 rounded-lg">
                <p class="text-sm text-gray-600 dark:text-gray-400">Total Expenses</p>
                <p class="text-2xl font-bold text-red-600 dark:text-red-400">${currencySymbol}${expenseTotal.toFixed(2)}</p>
            </div>
            ${Object.entries(expenseByCategory).map(([cat, amount]) => `
                <div class="flex justify-between p-2 border-b border-gray-200 dark:border-gray-700">
                    <span class="text-gray-700 dark:text-gray-300">${cat}</span>
                    <span class="font-semibold text-gray-800 dark:text-white">${currencySymbol}${amount.toFixed(2)}</span>
                </div>
            `).join('')}
        `;
        
        showNotification('Report generated successfully!', 'success');
    } catch (error) {
        showNotification('Failed to generate report', 'error');
    }
}

// Notification System
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `fixed top-20 right-4 z-50 px-6 py-3 rounded-lg shadow-lg text-white font-medium transition-all transform translate-x-0 ${
        type === 'success' ? 'bg-green-500' : 
        type === 'error' ? 'bg-red-500' : 
        'bg-blue-500'
    }`;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.transform = 'translateX(400px)';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Initialize
document.addEventListener('DOMContentLoaded', async () => {
    console.log('Dashboard initializing...');
    
    if (!checkAuth()) {
        console.log('Not authenticated, redirecting...');
        return;
    }
    
    initializeTheme();
    console.log('Theme initialized');
    
    // Set default dates
    const today = new Date().toISOString().split('T')[0];
    const firstDay = new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0];
    
    const incomeDateEl = document.getElementById('incomeDate');
    const expenseDateEl = document.getElementById('expenseDate');
    const reportStartEl = document.getElementById('reportStartDate');
    const reportEndEl = document.getElementById('reportEndDate');
    
    if (incomeDateEl) incomeDateEl.value = today;
    if (expenseDateEl) expenseDateEl.value = today;
    if (reportStartEl) reportStartEl.value = firstDay;
    if (reportEndEl) reportEndEl.value = today;
    
    console.log('Loading data...');
    
    // Load data with error handling for each
    try {
        await loadUserProfile();
        console.log('User profile loaded');
    } catch (error) {
        console.error('Error loading user profile:', error);
    }
    
    try {
        await loadCategories();
        console.log('Categories loaded');
    } catch (error) {
        console.error('Error loading categories:', error);
    }
    
    try {
        await loadDashboardData();
        console.log('Dashboard data loaded');
    } catch (error) {
        console.error('Error loading dashboard data:', error);
    }
    
    try {
        await loadIncomeList();
        console.log('Income list loaded');
    } catch (error) {
        console.error('Error loading income list:', error);
    }
    
    try {
        await loadExpenseList();
        console.log('Expense list loaded');
    } catch (error) {
        console.error('Error loading expense list:', error);
    }
    
    console.log('Dashboard initialization complete');
});
