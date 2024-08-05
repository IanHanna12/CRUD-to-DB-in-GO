document.addEventListener('DOMContentLoaded', function () {
    const loginForm = document.querySelector('form');

    if (loginForm) {
        loginForm.addEventListener('submit', function (e) {
            e.preventDefault();
            const username = document.querySelector('input[name="username"]').value;
            const password = document.querySelector('input[name="password"]').value;
            //TODO carry over username to main page
            localStorage.setItem('username', username);

            fetch('http://localhost:8080/login', {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({username, password}),
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('HTTP status ' + response.status);
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        console.log('Login successful:', data);

                        localStorage.setItem('sessionToken', data.session_token);
                        localStorage.setItem('username', username);

                        if (data.isAdmin) {
                            window.location.href = 'http://localhost:8080/static/admin/admin_view.html';
                        } else {
                            window.location.href = 'http://localhost:8080/static/user/user_view.html';
                        }
                    } else {
                        console.error('Login failed:', data.error || 'Unknown error');
                        alert('Login failed: ' + (data.error || 'Unknown error'));
                    }
                })
                .catch(error => {
                    console.error('Login error:', error);
                    alert('An error occurred during login: ' + error.message);
                });
        });
    }
});
