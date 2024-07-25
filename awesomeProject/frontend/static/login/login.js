document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM fully loaded');
    const loginForm = document.querySelector('form');
    console.log('Login form found:', !!loginForm);

    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            console.log('Login form submitted');
            const username = document.querySelector('input[type="text"]').value;
            const password = document.querySelector('input[type="password"]').value;

            console.log('Username:', username);
            console.log('Password:', password);

            // Simulating successful login for demonstration
            console.log('Login successful, redirecting...');
            localStorage.setItem('isAdmin', 'true');
            window.location.href = 'http://localhost:8080/static/main_page/mainpage.html';
        });
    }
});
