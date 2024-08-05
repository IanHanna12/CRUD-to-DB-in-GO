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





# Handlers

## LoginHandler
- Authenticates users
- Validates credentials against the database
- Creates and manages user sessions

## GetAllItemsHandler
- Retrieves all items for the authenticated user
- Filters items based on user ID

## CreateItemHandler
- Processes new item creation requests
- Associates items with the current user
- Validates input data

## GetItemByIDHandler
- Fetches a specific item by its unique identifier
- Ensures the item belongs to the requesting user

## UpdateItemHandler
- Handles updates to existing items
- Verifies user ownership before allowing changes

## DeleteItemByIDHandler
- Removes a specific item from the database
- Checks user authorization before deletion

## DeleteAllItemsHandler
- Deletes all items associated with the current user
- Requires user authentication

## PrefetchItemsHandler
- Retrieves items for client-side caching
- Optimizes performance for frequent data access

## ValidateSessionHandler
- Verifies the validity of user sessions
- Ensures secure and authenticated user states

## AuthMiddleware
- Implements authentication checks for protected routes
- Manages user permissions and access control
