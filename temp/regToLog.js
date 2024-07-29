function showMessage(message, isError = false) {
    const messageElement = document.getElementById('message');
    messageElement.textContent = message;
    messageElement.className = 'message ' + (isError ? 'error' : 'success');
    messageElement.style.display = 'block';
}

function showLoginForm() {
    document.getElementById('register-form').style.display = 'none';
    document.getElementById('login-form').style.display = 'block';
}

document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('registration-form').addEventListener('submit', function(e) {
        e.preventDefault();
        fetch('/register', {
            method: 'POST',
            body: new FormData(this)
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

    document.getElementById('login-form-element').addEventListener('submit', function(e) {
        e.preventDefault();
        fetch('/login', {
            method: 'POST',
            body: new FormData(this)
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.text();
        })
        .then(message => {
            showMessage(message);
            // Call loginSuccess function from main.js
            if (window.loginSuccess) {
                window.loginSuccess();
            }
        })
        .catch(error => {
            showMessage(error.message, true);
        });
    });
});