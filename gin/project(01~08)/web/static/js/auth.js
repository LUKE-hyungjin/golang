// Authentication JavaScript

// Handle login form
const loginForm = document.getElementById('loginForm');
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const data = await window.BlogPlatform.apiCall('/api/v1/auth/login', {
                method: 'POST',
                body: JSON.stringify({ username, password }),
            });

            // Store token and user info
            window.BlogPlatform.setToken(data.token);
            localStorage.setItem('user', JSON.stringify(data.user));

            window.BlogPlatform.showNotification('Login successful!', 'success');

            // Redirect to home page
            setTimeout(() => {
                window.location.href = '/';
            }, 1000);
        } catch (error) {
            window.BlogPlatform.showNotification(error.message || 'Login failed', 'error');
        }
    });
}

// Handle registration form
const registerForm = document.getElementById('registerForm');
if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        try {
            const data = await window.BlogPlatform.apiCall('/api/v1/auth/register', {
                method: 'POST',
                body: JSON.stringify({ username, email, password }),
            });

            window.BlogPlatform.showNotification('Registration successful! Please login.', 'success');

            // Redirect to login page
            setTimeout(() => {
                window.location.href = '/login';
            }, 1500);
        } catch (error) {
            window.BlogPlatform.showNotification(error.message || 'Registration failed', 'error');
        }
    });
}
