# Authentication 

## Login Process
1. User submits credentials (username and password) to the `/login` endpoint.
2. Server validates credentials against the database.
3. If valid, a new session is created:
    - `sessionID` and `sessionToken` are generated.
    - `sessionID` is stored in the user's database record.
    - `sessionToken` is set as a cookie in the response.

## Authenticated Requests
1. Client includes the `sessionToken` in the `Authorization` header of each request.
2. Server extracts the token from the `Authorization` header.
3. The `AuthMiddleware` validates the session:
    - Retrieves the user record using the `sessionID`.
    - Checks user permissions if admin access is required.
4. If valid, the request proceeds to the handler.

## Session Management
- Sessions are stored in the database, linked to user records.
- The `sessionToken` in the cookie is used to authenticate subsequent requests.
- Sessions expire after a set time period 

