# Authentication, Session Management, and Route Security
DEPRECATED
## Authentication

~~~~~~The system uses a username and password-based authentication mechanism. Passwords are securely hashed using bcrypt before storage in the database.

### User Creation
- New users can be created through the `/login` endpoint.
- The `CreateUser` function handles user creation, setting the `IsAdmin` field to 0 for admins and 1 for regular users.

### Login Process
1. Users submit credentials to the `/login` endpoint.
2. The system checks if the user exists:
    - If not, a new user is created.
    - If yes, the password is verified against the stored hash.
3. Upon successful authentication, the isAdmin (bool) flag is returned in header and used for all subsequent requests.




## Route Security

Routes are secured using middleware functions that check for valid isAdmin status 

### AuthMiddleware
- Verifies the presence of a valid admin status in the request header.

- For admin-only routes, it verifies the user's admin status and allows the request to proceed if the user is an admin, else blocked.

### Route Protection

- Authenticated routes: Require valid login credentials and login.
- Admin routes: Require valid admin status to proceed.

### Implementation
- The `AuthMiddleware` function is used as a wrapper for  routes in routing.

## Security 
- Passwords are never stored in plain text.
- Passwords are hashed using bcrypt before storage in the database.
- Admin status is determined by the `IsAdmin` field in the sent request header.
- All routes require valid login credentials and login.
- All admin routes require valid admin status to proceed.



-- 31.07.2024
user view works as intended

user_view: Post blog and edit the posts name and content works (editing author name of blog is blocked,
view by id works)


TODO: Test admin view