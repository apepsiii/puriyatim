// Main JavaScript file for Panti App

// DOM Content Loaded
document.addEventListener('DOMContentLoaded', function() {
    // Initialize components
    initDropdowns();
    initModals();
    initTooltips();
    initFormValidation();
    initFileUpload();
    initDatePickers();
    initDataTables();
    initCharts();
    initConfirmDialogs();
    initAutoLogout();
});

// Dropdown functionality
function initDropdowns() {
    const dropdowns = document.querySelectorAll('[data-dropdown]');
    
    dropdowns.forEach(dropdown => {
        const trigger = dropdown.querySelector('[data-dropdown-trigger]');
        const menu = dropdown.querySelector('[data-dropdown-menu]');
        
        if (trigger && menu) {
            trigger.addEventListener('click', function(e) {
                e.preventDefault();
                menu.classList.toggle('hidden');
                
                // Close other dropdowns
                dropdowns.forEach(otherDropdown => {
                    if (otherDropdown !== dropdown) {
                        const otherMenu = otherDropdown.querySelector('[data-dropdown-menu]');
                        if (otherMenu) {
                            otherMenu.classList.add('hidden');
                        }
                    }
                });
            });
        }
    });
    
    // Close dropdowns when clicking outside
    document.addEventListener('click', function(e) {
        dropdowns.forEach(dropdown => {
            const menu = dropdown.querySelector('[data-dropdown-menu]');
            if (menu && !dropdown.contains(e.target)) {
                menu.classList.add('hidden');
            }
        });
    });
}

// Modal functionality
function initModals() {
    const modals = document.querySelectorAll('[data-modal]');
    
    modals.forEach(modal => {
        const triggers = document.querySelectorAll(`[data-modal-trigger="${modal.id}"]`);
        const closeButtons = modal.querySelectorAll('[data-modal-close]');
        
        // Open modal
        triggers.forEach(trigger => {
            trigger.addEventListener('click', function(e) {
                e.preventDefault();
                modal.classList.remove('hidden');
                modal.classList.add('flex');
                document.body.style.overflow = 'hidden';
            });
        });
        
        // Close modal
        closeButtons.forEach(button => {
            button.addEventListener('click', function(e) {
                e.preventDefault();
                modal.classList.add('hidden');
                modal.classList.remove('flex');
                document.body.style.overflow = 'auto';
            });
        });
        
        // Close modal when clicking outside
        modal.addEventListener('click', function(e) {
            if (e.target === modal) {
                modal.classList.add('hidden');
                modal.classList.remove('flex');
                document.body.style.overflow = 'auto';
            }
        });
    });
}

// Tooltip functionality
function initTooltips() {
    const tooltips = document.querySelectorAll('[data-tooltip]');
    
    tooltips.forEach(tooltip => {
        tooltip.addEventListener('mouseenter', function() {
            const text = tooltip.getAttribute('data-tooltip');
            const tooltipElement = document.createElement('div');
            tooltipElement.className = 'absolute z-50 px-2 py-1 text-xs text-white bg-gray-900 rounded shadow-lg';
            tooltipElement.textContent = text;
            tooltipElement.style.bottom = '100%';
            tooltipElement.style.left = '50%';
            tooltipElement.style.transform = 'translateX(-50%) translateY(-4px)';
            tooltipElement.style.whiteSpace = 'nowrap';
            tooltipElement.setAttribute('data-tooltip-element', '');
            
            tooltip.style.position = 'relative';
            tooltip.appendChild(tooltipElement);
        });
        
        tooltip.addEventListener('mouseleave', function() {
            const tooltipElement = tooltip.querySelector('[data-tooltip-element]');
            if (tooltipElement) {
                tooltipElement.remove();
            }
        });
    });
}

// Form validation
function initFormValidation() {
    const forms = document.querySelectorAll('[data-validate]');
    
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            let isValid = true;
            const errors = [];
            
            // Validate required fields
            const requiredFields = form.querySelectorAll('[required]');
            requiredFields.forEach(field => {
                if (!field.value.trim()) {
                    isValid = false;
                    showFieldError(field, 'Field ini wajib diisi');
                    errors.push(field);
                } else {
                    clearFieldError(field);
                }
            });
            
            // Validate email fields
            const emailFields = form.querySelectorAll('[type="email"]');
            emailFields.forEach(field => {
                if (field.value && !isValidEmail(field.value)) {
                    isValid = false;
                    showFieldError(field, 'Format email tidak valid');
                    errors.push(field);
                }
            });
            
            // Validate phone fields
            const phoneFields = form.querySelectorAll('[data-validate-phone]');
            phoneFields.forEach(field => {
                if (field.value && !isValidPhone(field.value)) {
                    isValid = false;
                    showFieldError(field, 'Format nomor telepon tidak valid');
                    errors.push(field);
                }
            });
            
            if (!isValid) {
                e.preventDefault();
                // Focus on first error field
                if (errors.length > 0) {
                    errors[0].focus();
                }
            }
        });
    });
}

// Show field error
function showFieldError(field, message) {
    clearFieldError(field);
    
    field.classList.add('border-red-500');
    
    const errorElement = document.createElement('div');
    errorElement.className = 'text-red-500 text-sm mt-1';
    errorElement.textContent = message;
    errorElement.setAttribute('data-field-error', '');
    
    field.parentNode.appendChild(errorElement);
}

// Clear field error
function clearFieldError(field) {
    field.classList.remove('border-red-500');
    
    const errorElement = field.parentNode.querySelector('[data-field-error]');
    if (errorElement) {
        errorElement.remove();
    }
}

// Email validation
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Phone validation
function isValidPhone(phone) {
    const phoneRegex = /^[0-9+\-\s()]+$/;
    return phoneRegex.test(phone) && phone.replace(/\D/g, '').length >= 10;
}

// File upload functionality
function initFileUpload() {
    const fileInputs = document.querySelectorAll('input[type="file"]');
    
    fileInputs.forEach(input => {
        input.addEventListener('change', function() {
            const file = this.files[0];
            const preview = this.parentNode.querySelector('[data-file-preview]');
            
            if (file && preview) {
                if (file.type.startsWith('image/')) {
                    const reader = new FileReader();
                    reader.onload = function(e) {
                        preview.src = e.target.result;
                        preview.classList.remove('hidden');
                    };
                    reader.readAsDataURL(file);
                } else {
                    preview.classList.add('hidden');
                }
            }
            
            // Update file name display
            const fileName = this.parentNode.querySelector('[data-file-name]');
            if (fileName) {
                fileName.textContent = file ? file.name : 'Tidak ada file yang dipilih';
            }
        });
    });
}

// Date picker functionality
function initDatePickers() {
    const dateInputs = document.querySelectorAll('[data-date-picker]');
    
    dateInputs.forEach(input => {
        // Set max date to today for birth dates
        if (input.hasAttribute('data-max-today')) {
            input.max = new Date().toISOString().split('T')[0];
        }
        
        // Set min date to today for future dates
        if (input.hasAttribute('data-min-today')) {
            input.min = new Date().toISOString().split('T')[0];
        }
    });
}

// Data table functionality
function initDataTables() {
    const tables = document.querySelectorAll('[data-table]');
    
    tables.forEach(table => {
        const searchInput = table.parentNode.querySelector('[data-table-search]');
        const pagination = table.parentNode.querySelector('[data-table-pagination]');
        
        if (searchInput) {
            searchInput.addEventListener('input', function() {
                const searchTerm = this.value.toLowerCase();
                const rows = table.querySelectorAll('tbody tr');
                
                rows.forEach(row => {
                    const text = row.textContent.toLowerCase();
                    row.style.display = text.includes(searchTerm) ? '' : 'none';
                });
            });
        }
        
        // Sort functionality
        const headers = table.querySelectorAll('th[data-sort]');
        headers.forEach(header => {
            header.addEventListener('click', function() {
                const column = this.getAttribute('data-sort');
                const tbody = table.querySelector('tbody');
                const rows = Array.from(tbody.querySelectorAll('tr'));
                
                rows.sort((a, b) => {
                    const aText = a.querySelector(`td[data-column="${column}"]`).textContent;
                    const bText = b.querySelector(`td[data-column="${column}"]`).textContent;
                    return aText.localeCompare(bText);
                });
                
                rows.forEach(row => tbody.appendChild(row));
            });
        });
    });
}

// Chart functionality (placeholder for chart library integration)
function initCharts() {
    const charts = document.querySelectorAll('[data-chart]');
    
    charts.forEach(chart => {
        const type = chart.getAttribute('data-chart');
        const data = JSON.parse(chart.getAttribute('data-chart-data'));
        
        // This is a placeholder - integrate with Chart.js or similar library
        console.log(`Chart type: ${type}`, data);
    });
}

// Confirm dialog functionality
function initConfirmDialogs() {
    const confirmButtons = document.querySelectorAll('[data-confirm]');
    
    confirmButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            const message = this.getAttribute('data-confirm');
            if (!confirm(message)) {
                e.preventDefault();
            }
        });
    });
}

// Auto logout functionality
function initAutoLogout() {
    const autoLogoutElements = document.querySelectorAll('[data-auto-logout]');
    
    if (autoLogoutElements.length > 0) {
        let timeout;
        
        function resetTimeout() {
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                if (confirm('Anda akan keluar karena tidak ada aktivitas. Lanjutkan?')) {
                    window.location.href = '/admin/logout';
                }
            }, 30 * 60 * 1000); // 30 minutes
        }
        
        // Reset timeout on user activity
        document.addEventListener('mousemove', resetTimeout);
        document.addEventListener('keypress', resetTimeout);
        document.addEventListener('click', resetTimeout);
        document.addEventListener('scroll', resetTimeout);
        
        resetTimeout();
    }
}

// Utility functions
function formatCurrency(amount) {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR'
    }).format(amount);
}

function formatDate(date) {
    return new Intl.DateTimeFormat('id-ID', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    }).format(new Date(date));
}

function formatDateTime(date) {
    return new Intl.DateTimeFormat('id-ID', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    }).format(new Date(date));
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// API helper functions
async function apiRequest(url, options = {}) {
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    
    const mergedOptions = { ...defaultOptions, ...options };
    
    try {
        const response = await fetch(url, mergedOptions);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API request failed:', error);
        throw error;
    }
}

// Toast notification system
function showToast(message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `fixed bottom-4 right-4 z-50 p-4 rounded-md shadow-lg ${
        type === 'success' ? 'bg-green-500' : 
        type === 'error' ? 'bg-red-500' : 
        type === 'warning' ? 'bg-yellow-500' : 'bg-blue-500'
    } text-white`;
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    setTimeout(() => {
        toast.remove();
    }, 5000);
}

// Export functions for use in other scripts
window.PantiApp = {
    formatCurrency,
    formatDate,
    formatDateTime,
    debounce,
    apiRequest,
    showToast,
    showFieldError,
    clearFieldError
};