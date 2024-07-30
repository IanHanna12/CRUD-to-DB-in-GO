# Authentication System

## Key Features
- Secure user authentication with admin and regular user roles
- UUID session token generation for all authenticated sessions
- Middleware-based protection for authenticated and admin-only routes
- Session expiration and cleanup mechanism

## Implementation

### Session Structure:


type Session struct {
     UserID    string
     Username  string
     IsAdmin   bool
     ExpiresAt time.Time
}


### Login Process:
- Users authenticate via the `/login` endpoint
- Admin status determined by username "admin" and password check
- UUID session token generated for all successful logins
- Session stored in sessions map with user details and expiration

### Session Management:
- All sessions stored server-side in sessions map
- Client stores session UUID for subsequent requests
- Sessions include user ID, username, admin status, and expiration time

### Route Protection:
- `WithAuthentication` middleware secures authenticated-only routes
- `WithAdminAuthentication` middleware secures admin-only routes
- Verifies presence and validity of session UUID in request headers
- Checks session expiration and admin status for appropriate access

## Usage
- User login: Any non-empty username and password
- Admin login: Username "admin", password "admin"
- Protected routes: Include session UUID in request headers

## Code Snippets

### Login Handler


func LoginHandler(w http.ResponseWriter, r *http.Request) {
     // ... 
     if loginSuccessful {
         sessionID := uuid.New().String()
         sessions[sessionID] = Session{
             UserID:    uuid.New().String(),
             Username:  loginRequest.Username,
             IsAdmin:   isAdmin,
             ExpiresAt: time.Now().Add(24 * time.Hour),
         }
         // ... (return session ID and redirect URL to client)
     }
}



1. Login Initiation:
    - User submits credentials via the login form in login.html.
    - login.js captures the form submission and sends a POST request to '/login'.

2. Server-side Authentication (LoginHandler):
    - Receives login request and validates credentials.
    - For admin, checks against hardcoded "admin" username and password.
    - For regular users, accepts any non-empty password.

3. Session Creation:
    - Upon successful authentication, generates a UUID for the session.
    - Creates a Session struct with user details and expiration time.
    - Stores the session in the server-side 'sessions' map.

4. Client-side Session Storage:
    - Server sends back sessionID and user role (isAdmin).
    - login.js stores these in localStorage for subsequent requests.

5. Protected Route Access:
    - mainPage.js includes sessionID in the header of each API request.
    - Server uses middleware (WithAuthentication/WithAdminAuthentication) to validate the sessionID.
    - Middleware checks if the session exists and hasn't expired.

6. Role-based Access Control:
    - WithAdminAuthentication middleware ensures only admin users access certain routes.
    - Regular authentication allows access to user-level functionalities.

7. Session Management:
    - Sessions have a 24-hour expiration.
    - Sessions are stored in a server-side map.
