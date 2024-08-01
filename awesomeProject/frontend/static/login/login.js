document.addEventListener('DOMContentLoaded', function () {
    const loginForm = document.querySelector('form');

    if (loginForm) {
        loginForm.addEventListener('submit', function (e) {
            e.preventDefault();
            const username = document.querySelector('input[name="username"]').value;
            const password = document.querySelector('input[name="password"]').value;

            fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            })
                .then(response => response.json())
                .then(data => {
                    console.log('Login response:', data);
                    if (data.success) {
                        localStorage.setItem('isAdmin', data.isAdmin);
                        console.log('Redirecting to:', data.redirectURL);
                        window.location.href = data.redirectURL;
                    } else {
                        alert('Login failed: ' + (data.error || 'Unknown error'));
                    }
                })
                .catch(error => {
                    console.error('Login error:', error);
                    alert('An error occurred during login:' + error);
                });
        });
    }
});