# Authentication and Data Handling 


## Authentication Flow

### Login Process
1. User enters credentials on the login page.
2. Client sends a POST request to `/login` with username and password.
3. Server validates credentials against the database.
4. If valid, server generates a unique session token (UUID).
5. Server stores the session token in the database, associated with the user.
6. Server responds with:
    - Session token
    - User information (including username)
    - Admin status
7. Client stores:
    - Session token in an HTTP-only cookie
    - Username in localStorage

### Session Management
- Session tokens are UUIDs.
- Tokens are stored in the database with an association to the user ID.
- HTTP-only cookies prevent client-side JavaScript access to the token.

### Authentication Middleware (AuthMiddleware)
- Extracts session token from the request's cookies.
- Queries the database to find a user with the given session token.
- If found, attaches user information to the request context.
- Checks if admin access is required for the route.
- If admin access is required and user is not admin, returns 401 Unauthorized.
- If session is invalid, returns 401 Unauthorized.

### Client-Side Authentication
- `fetchWithAuth` function:
    - Retrieves session token from cookies.
    - Adds 'Authorization' header with format: `Bearer {sessionToken}`.
    - Sets `credentials: 'include'` to ensure cookies are sent with cross-origin requests.
    - Handles 401 responses, potentially redirecting to login page.

### Session Validation
- `/validate-session` endpoint:
    - Extracts session token from request.
    - Queries database for user with matching session token.
    - If found, returns JSON response with `{valid: true}`.
    - If not found, returns 401 Unauthorized.

## Data Handling and Prefetching

### Prefetching Mechanism
- `prefetchItems()` function:
    - Called on page load.
    - Sends GET request to `${apiUrl}/all`.
    - Server filters items based on authenticated user's ID.
    - Received items are stored in `prefetchedItems` array.

### Item Operations
- View:
    - Renders items from `prefetchedItems` array.
    - If item not in array, fetches from server using `${apiUrl}/single/{id}`.
- Create:
    - Sends POST request to `${apiUrl}/create`.
    - On success, adds new item to `prefetchedItems`.
- Edit:
    - Sends PUT request to `${apiUrl}/update/{id}`.
    - On success, updates item in `prefetchedItems`.
- Delete:
    - Sends DELETE request to `${apiUrl}/delete/{id}`.
    - On success, removes item from `prefetchedItems`.

### Rendering
- `renderPosts()` function:
    - Maps over `prefetchedItems` array.
    - Generates HTML for each item using `renderPost()` function.
    - Joins HTML strings and sets as innerHTML of `postsContainer`.

## Protection of User-Specific Content

### User ID Association
- Each Item struct includes a UserID field (UUID).

### Server-Side Filtering
- In `GetAllItemsHandler`:
    - `DB.Where("user_id = ?", userID).Find(&items)`
    - Only returns items belonging to the authenticated user.

### Authentication Middleware
- Attaches user information to request context.
- Used in all item-related handlers to verify ownership.
- Checks if user is admin.
- If user is not admin, returns 401 Unauthorized.

### Item-Specific Operations
- Update and Delete operations check item ownership:
    - `DB.Where("id = ? AND user_id = ?", itemID, userID).First(&item)`

### Client-Side Restrictions
- JavaScript only displays and allows operations on fetched items.
- Fetched items are inherently filtered by the server based on user ID.

## API Endpoints
- POST `/login`: Authenticate user and create session
- GET `/validate-session`: Check session validity
- GET `/items/all`: Retrieve all user-specific items
- GET `/items/single/:id`: Retrieve a specific item
- POST `/items/create`: Create a new item
- PUT `/items/update/:id`: Update an existing item
- DELETE `/items/delete/:id`: Delete a specific item

## Security Measures
- Password Hashing:
    - Passwords are hashed using bcrypt before storage.
- HTTP-only Cookies:
    - Session tokens stored in HTTP-only cookies.
- Server-Side Validation:
    - All user actions are validated on the server, not trusting client-side data.
- CORS Configuration:
    - Strict CORS policies to prevent unauthorized cross-origin requests.

## Route Security
- Authentication Middleware:
    - All sensitive routes are protected by AuthMiddleware.
    - Verifies user's session token before allowing access.
- Role-Based Access Control:
    - Distinguishes between regular users and admins.
    - Certain routes (e.g., `/items/prefetch/all`) are restricted to admin users.
- User-Specific Content Protection:
    - Server-side filtering ensures users can only access their own data.
- Secure Session Management:
    - Session tokens stored as UUIDs in HTTP-only cookies.

  - Session tokens: validate request access based on session token.
  - Session tokens: stored in database with user ID association.
  - sessionID: UUID for unique session identification.
  - 