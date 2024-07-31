function showMessage(message, isError = false) {
    const messageElement = document.getElementById('message');
    if (messageElement) {
        messageElement.textContent = message;
        messageElement.className = 'message ' + (isError ? 'error' : 'success');
        messageElement.style.display = 'block';
    }
}

function showLoginForm() {
    const registerForm = document.getElementById('register-form');
    const loginForm = document.getElementById('login-form-element');
    if (registerForm && loginForm) {
        registerForm.style.display = 'none';
        loginForm.style.display = 'block';
    }
}

function showRegisterForm() {
    const registerForm = document.getElementById('register-form');
    const loginForm = document.getElementById('login-form-element');
    if (registerForm && loginForm) {
        registerForm.style.display = 'block';
        loginForm.style.display = 'none';
    }
}

document.addEventListener('DOMContentLoaded', function() {
    const registrationForm = document.getElementById('registration-form');
    const loginFormElement = document.getElementById('login-form');

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
                }
            })
            .catch(error => {
                showMessage(error.message, true);
            });
        });
    }
});
