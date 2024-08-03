function showMessage(message, isError = false) {
    const messageElement = document.getElementById('message');
    if (messageElement) {
        messageElement.textContent = message;
        messageElement.className = 'message ' + (isError ? 'error' : 'success');
        messageElement.style.display = 'block';
    } else {
        console.error('Message element not found');
    }
}

function showLoginForm() {
    const registerForm = document.getElementById('register-form');
    const loginForm = document.getElementById('login-form-element');
    if (registerForm && loginForm) {
        registerForm.style.display = 'none';
        loginForm.style.display = 'block';
    } else {
        console.error('Register or login form not found');
    }
}

function showRegisterForm() {
    const registerForm = document.getElementById('register-form');
    const loginForm = document.getElementById('login-form-element');
    if (registerForm && loginForm) {
        registerForm.style.display = 'block';
        loginForm.style.display = 'none';
    } else {
        console.error('Register or login form not found');
    }
}

document.addEventListener('DOMContentLoaded', function() {
    const registrationForm = document.getElementById('registration-form');
    const loginFormElement = document.getElementById('login-form');
    const goToLoginButton = document.getElementById('go-to-login');
    const goToRegisterButton = document.getElementById('go-to-register');

    if (registrationForm) {
        registrationForm.addEventListener('submit', function(e) {
            e.preventDefault();
            fetch('/register', {
                method: 'POST',
                body: new FormData(registrationForm)
            })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => { throw new Error(text) });
                }
                return response.text();
            })
            .then(message => {
                showMessage(message);
                setTimeout(() => {
                    showLoginForm();
                }, 2000);
            })
            .catch(error => {
                showMessage(error.message, true);
            });
        });
    } else {
        console.error('Registration form not found');
    }

    if (loginFormElement) {
        loginFormElement.addEventListener('submit', function(e) {
            e.preventDefault();
            fetch('/login', {
                method: 'POST',
                body: new FormData(loginFormElement)
            })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => { throw new Error(text) });
                }
                return response.text();
            })
            .then(message => {
                showMessage(message);
                if (window.loginSuccess) {
                    window.loginSuccess();
                } else {
                    console.error('loginSuccess function not found');
                }
            })
            .catch(error => {
                showMessage(error.message, true);
            });
        });
    } else {
        console.error('Login form not found');
    }

    // Add event listeners for switching between forms
    if (goToLoginButton) {
        goToLoginButton.addEventListener('click', showLoginForm);
    }
    if (goToRegisterButton) {
        goToRegisterButton.addEventListener('click', showRegisterForm);
    }
});