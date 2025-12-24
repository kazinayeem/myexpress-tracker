// API Configuration
const API_BASE = window.location.origin + '/api';
let dailyChart = null;

// Check if user is authenticated
function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login.html';
        return false;
    }
    return true;
}

// Helper function to make authenticated requests
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

// Initialize dashboard
async function initDashboard() {
    if (!checkAuth()) return;
    
    // Set username
    const username = localStorage.getItem('username');
    document.getElementById('user-name').textContent = `Welcome, ${username}`;
    
    // Set today's date as default
    const today = new Date().toISOString().split('T')[0];
    document.getElementById('income-date').value = today;
    document.getElementById('expense-date').value = today;
    document.getElementById('date-filter').value = today;
    
    // Load data
    await loadDashboard();
    await loadCategories();
    await loadTransactions();
}

// Load dashboard summary
async function loadDashboard() {
    try {
        const response = await apiRequest('/dashboard');
        const data = await response.json();
        
        // Update summary cards
        document.getElementById('total-income').textContent = `$${data.total_income.toFixed(2)}`;
        document.getElementById('total-expense').textContent = `$${data.total_expense.toFixed(2)}`;
        document.getElementById('balance').textContent = `$${data.balance.toFixed(2)}`;
        document.getElementById('today-income').textContent = `$${data.today_income.toFixed(2)}`;
        document.getElementById('today-expense').textContent = `$${data.today_expense.toFixed(2)}`;
        
        // Update balance color
        const balanceElement = document.getElementById('balance');
        if (data.balance >= 0) {
            balanceElement.className = 'amount income';
        } else {
            balanceElement.className = 'amount expense';
        }
        
        // Load chart
        if (data.daily_data && data.daily_data.length > 0) {
            loadChart(data.daily_data);
        }
    } catch (error) {
        console.error('Failed to load dashboard:', error);
    }
}

// Load chart
function loadChart(dailyData) {
    const ctx = document.getElementById('daily-chart').getContext('2d');
    
    // Destroy existing chart if it exists
    if (dailyChart) {
        dailyChart.destroy();
    }
    
    const labels = dailyData.map(d => d.date);
    const incomeData = dailyData.map(d => d.income);
    const expenseData = dailyData.map(d => d.expense);
    
    dailyChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Income',
                    data: incomeData,
                    borderColor: '#27ae60',
                    backgroundColor: 'rgba(39, 174, 96, 0.1)',
                    tension: 0.4
                },
                {
                    label: 'Expense',
                    data: expenseData,
                    borderColor: '#e74c3c',
                    backgroundColor: 'rgba(231, 76, 60, 0.1)',
                    tension: 0.4
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'top',
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return '$' + value.toFixed(0);
                        }
                    }
                }
            }
        }
    });
}

// Load categories
async function loadCategories() {
    try {
        const response = await apiRequest('/categories');
        const categories = await response.json();
        
        const incomeCategories = categories.filter(c => c.type === 'income');
        const expenseCategories = categories.filter(c => c.type === 'expense');
        
        // Populate income category select
        const incomeSelect = document.getElementById('income-category');
        incomeSelect.innerHTML = incomeCategories.map(c => 
            `<option value="${c.id}">${c.name}</option>`
        ).join('');
        
        // Populate expense category select
        const expenseSelect = document.getElementById('expense-category');
        expenseSelect.innerHTML = expenseCategories.map(c => 
            `<option value="${c.id}">${c.name}</option>`
        ).join('');
        
        // Populate filter
        const categoryFilter = document.getElementById('category-filter');
        const transactionType = document.getElementById('transaction-type').value;
        const filterCategories = transactionType === 'income' ? incomeCategories : expenseCategories;
        categoryFilter.innerHTML = '<option value="">All Categories</option>' + 
            filterCategories.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
    } catch (error) {
        console.error('Failed to load categories:', error);
    }
}

// Load transactions
async function loadTransactions() {
    const transactionType = document.getElementById('transaction-type').value;
    const categoryId = document.getElementById('category-filter').value;
    const date = document.getElementById('date-filter').value;
    
    let endpoint = `/${transactionType}?`;
    if (categoryId) endpoint += `category_id=${categoryId}&`;
    if (date) endpoint += `date=${date}`;
    
    try {
        const response = await apiRequest(endpoint);
        const transactions = await response.json();
        
        const listElement = document.getElementById('transactions-list');
        
        if (!transactions || transactions.length === 0) {
            listElement.innerHTML = '<p>No transactions found.</p>';
            return;
        }
        
        listElement.innerHTML = transactions.map(t => {
            const amount = transactionType === 'income' ? t.amount : t.amount;
            const date = transactionType === 'income' ? t.income_date : t.expense_date;
            
            return `
                <div class="transaction-item ${transactionType}">
                    <div class="transaction-info">
                        <h4>${t.category_name}</h4>
                        <p>${t.description || 'No description'}</p>
                        <p><small>${date}</small></p>
                    </div>
                    <div>
                        <div class="transaction-amount ${transactionType}">
                            ${transactionType === 'income' ? '+' : '-'}$${amount.toFixed(2)}
                        </div>
                        <div class="transaction-actions">
                            <button class="btn-icon" onclick="edit${transactionType.charAt(0).toUpperCase() + transactionType.slice(1)}(${t.id})" title="Edit">✏️</button>
                            <button class="btn-icon" onclick="delete${transactionType.charAt(0).toUpperCase() + transactionType.slice(1)}(${t.id})" title="Delete">🗑️</button>
                        </div>
                    </div>
                </div>
            `;
        }).join('');
    } catch (error) {
        console.error('Failed to load transactions:', error);
    }
}

// Clear filters
function clearFilters() {
    document.getElementById('category-filter').value = '';
    document.getElementById('date-filter').value = '';
    loadTransactions();
}

// Modal functions
function showAddIncomeModal() {
    document.getElementById('income-id').value = '';
    document.getElementById('income-form').reset();
    document.getElementById('income-modal-title').textContent = 'Add Income';
    document.getElementById('income-date').value = new Date().toISOString().split('T')[0];
    document.getElementById('income-modal').classList.remove('hidden');
}

function showAddExpenseModal() {
    document.getElementById('expense-id').value = '';
    document.getElementById('expense-form').reset();
    document.getElementById('expense-modal-title').textContent = 'Add Expense';
    document.getElementById('expense-date').value = new Date().toISOString().split('T')[0];
    document.getElementById('expense-modal').classList.remove('hidden');
}

function closeModals() {
    document.getElementById('income-modal').classList.add('hidden');
    document.getElementById('expense-modal').classList.add('hidden');
}

// Income form handler
document.getElementById('income-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const id = document.getElementById('income-id').value;
    const data = {
        category_id: parseInt(document.getElementById('income-category').value),
        amount: parseFloat(document.getElementById('income-amount').value),
        income_date: document.getElementById('income-date').value,
        description: document.getElementById('income-description').value
    };
    
    try {
        let response;
        if (id) {
            response = await apiRequest(`/income/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            });
        } else {
            response = await apiRequest('/income', {
                method: 'POST',
                body: JSON.stringify(data)
            });
        }
        
        if (response.ok) {
            closeModals();
            await loadDashboard();
            await loadTransactions();
        } else {
            const error = await response.json();
            alert(error.error || 'Failed to save income');
        }
    } catch (error) {
        alert('Failed to save income');
    }
});

// Expense form handler
document.getElementById('expense-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const id = document.getElementById('expense-id').value;
    const data = {
        category_id: parseInt(document.getElementById('expense-category').value),
        amount: parseFloat(document.getElementById('expense-amount').value),
        expense_date: document.getElementById('expense-date').value,
        description: document.getElementById('expense-description').value
    };
    
    try {
        let response;
        if (id) {
            response = await apiRequest(`/expense/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            });
        } else {
            response = await apiRequest('/expense', {
                method: 'POST',
                body: JSON.stringify(data)
            });
        }
        
        if (response.ok) {
            closeModals();
            await loadDashboard();
            await loadTransactions();
        } else {
            const error = await response.json();
            alert(error.error || 'Failed to save expense');
        }
    } catch (error) {
        alert('Failed to save expense');
    }
});

// Edit functions
async function editIncome(id) {
    try {
        const response = await apiRequest(`/income?date=`);
        const incomes = await response.json();
        const income = incomes.find(i => i.id === id);
        
        if (income) {
            document.getElementById('income-id').value = income.id;
            document.getElementById('income-category').value = income.category_id;
            document.getElementById('income-amount').value = income.amount;
            document.getElementById('income-date').value = income.income_date;
            document.getElementById('income-description').value = income.description;
            document.getElementById('income-modal-title').textContent = 'Edit Income';
            document.getElementById('income-modal').classList.remove('hidden');
        }
    } catch (error) {
        alert('Failed to load income details');
    }
}

async function editExpense(id) {
    try {
        const response = await apiRequest(`/expense?date=`);
        const expenses = await response.json();
        const expense = expenses.find(e => e.id === id);
        
        if (expense) {
            document.getElementById('expense-id').value = expense.id;
            document.getElementById('expense-category').value = expense.category_id;
            document.getElementById('expense-amount').value = expense.amount;
            document.getElementById('expense-date').value = expense.expense_date;
            document.getElementById('expense-description').value = expense.description;
            document.getElementById('expense-modal-title').textContent = 'Edit Expense';
            document.getElementById('expense-modal').classList.remove('hidden');
        }
    } catch (error) {
        alert('Failed to load expense details');
    }
}

// Delete functions
async function deleteIncome(id) {
    if (!confirm('Are you sure you want to delete this income?')) return;
    
    try {
        const response = await apiRequest(`/income/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            await loadDashboard();
            await loadTransactions();
        } else {
            alert('Failed to delete income');
        }
    } catch (error) {
        alert('Failed to delete income');
    }
}

async function deleteExpense(id) {
    if (!confirm('Are you sure you want to delete this expense?')) return;
    
    try {
        const response = await apiRequest(`/expense/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            await loadDashboard();
            await loadTransactions();
        } else {
            alert('Failed to delete expense');
        }
    } catch (error) {
        alert('Failed to delete expense');
    }
}

// Export to PDF
async function exportPDF() {
    const startDate = prompt('Start date (YYYY-MM-DD):', new Date(Date.now() - 30*24*60*60*1000).toISOString().split('T')[0]);
    const endDate = prompt('End date (YYYY-MM-DD):', new Date().toISOString().split('T')[0]);
    
    if (startDate && endDate) {
        const token = localStorage.getItem('token');
        window.open(`${API_BASE}/export/pdf?start_date=${startDate}&end_date=${endDate}`, '_blank');
    }
}

// Logout function
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    window.location.href = '/login.html';
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', initDashboard);

// Update category filter when transaction type changes
document.getElementById('transaction-type').addEventListener('change', async () => {
    await loadCategories();
    await loadTransactions();
});
