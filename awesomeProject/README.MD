# Authentication and Data Handling System

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
- Session tokens are UUID 
- Tokens are stored in the database with an association to the user ID.
- HTTP-only cookies prevent client-side JavaScript access to the token.

### Authentication Middleware (AuthMiddleware)
1. Extracts session token from the request's cookies.
2. Queries the database to find a user with the given session token.
3. If found, attaches user information to the request context.
4. Checks if admin access is required for the route.
5. If admin access is required and user is not admin, returns 401 Unauthorized.
6. If session is invalid, returns 401 Unauthorized.

### Client-Side Authentication
- `fetchWithAuth` function:
   1. Retrieves session token from cookies.
   2. Adds 'Authorization' header with format: `Bearer {sessionToken}`.
   3. Sets `credentials: 'include'` to ensure cookies are sent with cross-origin requests.
   4. Handles 401 responses, potentially redirecting to login page.

### Session Validation
- `/validate-session` endpoint:
   1. Extracts session token from request.
   2. Queries database for user with matching session token.
   3. If found, returns JSON response with `{valid: true}`.
   4. If not found, returns 401 Unauthorized.

## Data Handling and Prefetching

### Prefetching Mechanism
- `prefetchItems()` function:
   1. Called on page load.
   2. Sends GET request to `${apiUrl}/all`.
   3. Server filters items based on authenticated user's ID.
   4. Received items are stored in `prefetchedItems` array.

### Item Operations
- View:
   1. Renders items from `prefetchedItems` array.
   2. If item not in array, fetches from server using `${apiUrl}/single/{id}`.
- Create:
   1. Sends POST request to `${apiUrl}/create`.
   2. On success, adds new item to `prefetchedItems`.
- Edit:
   1. Sends PUT request to `${apiUrl}/update/{id}`.
   2. On success, updates item in `prefetchedItems`.
- Delete:
   1. Sends DELETE request to `${apiUrl}/delete/{id}`.
   2. On success, removes item from `prefetchedItems`.

### Rendering
- `renderPosts()` function:
   1. Maps over `prefetchedItems` array.
   2. Generates HTML for each item using `renderPost()` function.
   3. Joins HTML strings and sets as innerHTML of postsContainer.

## Protection of User-Specific Content

1. User ID Association:
   - Each Item struct includes a UserID field (UUID).

2. Server-Side Filtering:
   - In GetAllItemsHandler:
     go
     DB.Where("user_id = ?", userID).Find(&items)
     
   - Only returns items belonging to the authenticated user.

3. Authentication Middleware:
   - Attaches user information to request context.
   - Used in all item-related handlers to verify ownership.
   - Checks if user is admin.
   - If user is not admin, returns 401 Unauthorized.

4. Item-Specific Operations:
   - Update and Delete operations check item ownership:
     go
     DB.Where("id = ? AND user_id = ?", itemID, userID).First(&item)
     

5. Client-Side Restrictions:
   - JavaScript only displays and allows operations on fetched items.
   - Fetched items are inherently filtered by the server based on user ID.

## API Endpoints

- `POST /login`: Authenticate user and create session
- `GET /validate-session`: Check session validity
- `GET /items/all`: Retrieve all user-specific items
- `GET /items/single/:id`: Retrieve a specific item
- `POST /items/create`: Create a new item
- `PUT /items/update/:id`: Update an existing item
- `DELETE /items/delete/:id`: Delete a specific item

## Security Measures

1. Password Hashing:
   - Passwords are hashed using bcrypt before storage.


3. HTTP-only Cookies:
   - Session tokens stored in HTTP-only cookies 

4. Server-Side Validation:
   - All user actions are validated on the server, not trusting client-side data.

5. CORS Configuration:
   - Strict CORS policies to prevent unauthorized cross-origin requests.

This system provides a comprehensive, secure, and efficient approach to user authentication and data management, ensuring data privacy and integrity throughout the application.



_____________________________-

prefetch: i load all data associated with the user through user id and store it (prefetchedItems) --> filter it based on
user id

all subsequent requests are made to the fetched items array

thus making it impossible to access data that is not associated with the user