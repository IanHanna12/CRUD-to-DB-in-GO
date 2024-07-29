featureBranch02 brought up to date with featureBranch01

26.07.2024:
all functionalities of role based CRUD to database implemented, db used: MariaDB.



# Admin Authentication Enhancement


# Admin Authentication System

## Overview
This project implements a robust admin authentication system for a blog CRUD application, using UUID-based sessions and server-side verification.

## Key Features
- Secure admin login with predefined credentials
- UUID session token generation for admin sessions
- Middleware-based protection for admin-only routes

## Implementation
1. Login Process:
    - Users authenticate via the `/login` endpoint
    - Admin status determined by server-side credential check
   username: "admin", password: "admin" (gets hashed server-side)
    - UUID session token generated for admin logins
    - Admin session stored in `adminSessions` map

2. Admin Session Management:
    - Admin sessions stored server-side in `adminSessions` map
    - Client stores session UUID in local storage

3. Route Protection:
    - `WithAdminAuth` middleware secures admin-only routes
    - Verifies presence and validity of admin session UUID in request headers
    - Admin session UUID stored in request headers and checked against `adminSessions` map
## Usage
- Admin login: Use predefined credentials (username: "admin", password: "admin")
- Protected routes: Include admin session UUID in request headers
- Non-admin users: Automatically restricted from accessing admin functions