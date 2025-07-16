No problem, here is the API route design in Markdown format.

-----

# API Route Design

**Base URL:** `/api/v1`

## I. Authentication & User Management Routes

These routes handle player registration, login, and token management.

1.  **Player Registration**

      * **Endpoint:** `/api/v1/auth/register`
      * **Method:** `POST`
      * **Purpose:** Allows new players to create an account.
      * **Authentication:** None (public access)
      * **Request Body (JSON):**
        ```json
        {
            "username": "player123",
            "email": "player@example.com",
            "password": "secure_password_123"
        }
        ```
      * **Response (Success - HTTP 201 Created):**
        ```json
        {
            "message": "Player registered successfully",
            "player_id": "uuid-of-new-player"
        }
        ```
      * **Response (Error - HTTP 400 Bad Request, 409 Conflict):**
        ```json
        {
            "error": "Username or email already exists"
        }
        ```
      * **Gin Handler (Conceptual):** `auth.RegisterPlayer`

2.  **Player Login**

      * **Endpoint:** `/api/v1/auth/login`
      * **Method:** `POST`
      * **Purpose:** Authenticates a player and issues a JWT.
      * **Authentication:** None (public access)
      * **Request Body (JSON):**
        ```json
        {
            "username": "player123",
            "password": "secure_password_123"
        }
        ```
      * **Response (Success - HTTP 200 OK):**
        ```json
        {
            "message": "Login successful",
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "player_id": "uuid-of-player"
        }
        ```
      * **Response (Error - HTTP 401 Unauthorized, 400 Bad Request):**
        ```json
        {
            "error": "Invalid credentials"
        }
        ```
      * **Gin Handler (Conceptual):** `auth.LoginPlayer`

3.  **Token Refresh (Optional but Recommended)**

      * **Endpoint:** `/api/v1/auth/refresh`
      * **Method:** `POST`
      * **Purpose:** Allows players to refresh their JWT using a refresh token (if implemented) without re-logging in.
      * **Authentication:** Requires a valid, unexpired refresh token.
      * **Request Body (JSON):**
        ```json
        {
            "refresh_token": "your_refresh_token_here"
        }
        ```
      * **Response (Success - HTTP 200 OK):**
        ```json
        {
            "message": "Token refreshed successfully",
            "token": "new_jwt_token_here"
        }
        ```
      * **Response (Error - HTTP 401 Unauthorized):**
        ```json
        {
            "error": "Invalid or expired refresh token"
        }
        ```
      * **Gin Handler (Conceptual):** `auth.RefreshToken`

-----

## II. Player Location & Data Routes

These routes handle real-time player position updates and fetching player-related information.

1.  **Update Player Location**

      * **Endpoint:** `/api/v1/players/:player_id/location`
      * **Method:** `PUT` or `PATCH` (PUT for full replacement, PATCH for partial update; PATCH is slightly more semantically correct for just updating location, but PUT is often used for simplicity for "set current state").
      * **Purpose:** Allows a player to send their real-time location updates. This will be the **high-throughput, low-latency** endpoint.
      * **Authentication:** Required (JWT from logged-in player). Authorization: Player must be updating their *own* location.
      * **Request Body (JSON):**
        ```json
        {
            "latitude": 34.0522,
            "longitude": -118.2437,
            "timestamp": "2025-07-16T14:30:00Z" // ISO 8601 format
        }
        ```
      * **Response (Success - HTTP 200 OK or 204 No Content):**
        ```json
        {
            "message": "Location updated successfully"
        }
        ```
        (Or an empty body with 204 if no additional information is needed).
      * **Response (Error - HTTP 400 Bad Request, 401 Unauthorized, 403 Forbidden, 500 Internal Server Error):**
        ```json
        {
            "error": "Invalid location data"
        }
        ```
      * **Gin Handler (Conceptual):** `player.UpdatePlayerLocation`

2.  **Get Player Locations (for Map Rendering)**

      * **Endpoint:** `/api/v1/players/locations`
      * **Method:** `GET`
      * **Purpose:** Retrieves the locations of multiple players for map rendering. Could support filtering for nearby players.
      * **Authentication:** Required (JWT from any authenticated player).
      * **Query Parameters (Optional):**
          * `center_lat`: Latitude of the map center.
          * `center_lon`: Longitude of the map center.
          * `radius_km`: Radius in kilometers to search for players around the center.
          * `limit`: Maximum number of players to return.
      * **Response (Success - HTTP 200 OK):**
        ```json
        {
            "players": [
                {
                    "player_id": "uuid-of-player-1",
                    "username": "Alice",
                    "latitude": 34.0530,
                    "longitude": -118.2440,
                    "timestamp": "2025-07-16T14:30:05Z"
                },
                {
                    "player_id": "uuid-of-player-2",
                    "username": "Bob",
                    "latitude": 34.0510,
                    "longitude": -118.2420,
                    "timestamp": "2025-07-16T14:30:02Z"
                }
            ]
        }
        ```
      * **Response (Error - HTTP 401 Unauthorized, 500 Internal Server Error):**
        ```json
        {
            "error": "Authentication required"
        }
        ```
      * **Gin Handler (Conceptual):** `player.GetPlayerLocations`

-----

## III. Points of Interest (POI) Routes

These routes allow for the management and retrieval of POIs. Administrative or Game Master roles would typically manage these.

1.  **Create POI**

      * **Endpoint:** `/api/v1/pois/create`
      * **Method:** `POST`
      * **Purpose:** Adds a new Point of Interest to the game world.
      * **Authentication:** Required (JWT). **Authorization:** Must have `admin` or `game_master` role.
      * **Request Body (JSON):**
        ```json
        {
            "name": "Healing Spring",
            "type": "healing_station",
            "latitude": 34.0600,
            "longitude": -118.2500,
            "description": "A mystical spring that restores health.",
            "metadata": {
                "heal_amount": 50,
                "respawn_time_seconds": 300
            }
        }
        ```
      * **Response (Success - HTTP 201 Created):**
        ```json
        {
            "message": "POI created successfully",
            "poi_id": "uuid-of-new-poi"
        }
        ```
      * **Response (Error - HTTP 400 Bad Request, 401 Unauthorized, 403 Forbidden):**
        ```json
        {
            "error": "Insufficient permissions"
        }
        ```
      * **Gin Handler (Conceptual):** `poi.CreatePOI`

2.  **Get POIs**

      * **Endpoint:** `/api/v1/pois`
      * **Method:** `GET`
      * **Purpose:** Retrieves a list of Points of Interest. Can be filtered by type, proximity, etc.
      * **Authentication:** Required (JWT).
      * **Query Parameters (Optional):**
          * `type`: Filter by POI type (e.g., `healing_station`, `quest_giver`).
          * `center_lat`, `center_lon`, `radius_km`: For geospatial search.
          * `limit`, `offset`: For pagination.
      * **Response (Success - HTTP 200 OK):**
        ```json
        {
            "pois": [
                {
                    "poi_id": "uuid-of-poi-1",
                    "name": "Healing Spring",
                    "type": "healing_station",
                    "latitude": 34.0600,
                    "longitude": -118.2500,
                    "description": "A mystical spring that restores health.",
                    "metadata": { /* ... */ }
                },
                {
                    "poi_id": "uuid-of-poi-2",
                    "name": "Dragon's Lair",
                    "type": "boss_area",
                    "latitude": 34.0700,
                    "longitude": -118.2600,
                    "description": "Beware the mighty beast!",
                    "metadata": { /* ... */ }
                }
            ]
        }
        ```
      * **Response (Error - HTTP 401 Unauthorized, 500 Internal Server Error):**
        ```json
        {
            "error": "Authentication required"
        }
        ```
      * **Gin Handler (Conceptual):** `poi.GetPOIs`

3.  **Get Single POI by ID**

      * **Endpoint:** `/api/v1/pois/:poi_id`
      * **Method:** `GET`
      * **Purpose:** Retrieves detailed information for a specific POI.
      * **Authentication:** Required (JWT).
      * **Response (Success - HTTP 200 OK):**
        ```json
        {
            "poi_id": "uuid-of-poi-1",
            "name": "Healing Spring",
            "type": "healing_station",
            "latitude": 34.0600,
            "longitude": -118.2500,
            "description": "A mystical spring that restores health.",
            "metadata": {
                "heal_amount": 50,
                "respawn_time_seconds": 300
            },
            "created_at": "2025-07-15T10:00:00Z",
            "updated_at": "2025-07-16T14:00:00Z"
        }
        ```
      * **Response (Error - HTTP 401 Unauthorized, 404 Not Found):**
        ```json
        {
            "error": "POI not found"
        }
        ```
      * **Gin Handler (Conceptual):** `poi.GetPOIByID`

4.  **Update POI**

      * **Endpoint:** `/api/v1/pois/:poi_id`
      * **Method:** `PUT` or `PATCH` (PUT for full replacement, PATCH for partial update).
      * **Purpose:** Modifies an existing Point of Interest.
      * **Authentication:** Required (JWT). **Authorization:** Must have `admin` or `game_master` role.
      * **Request Body (JSON - similar to Create POI, but often with optional fields for PATCH):**
        ```json
        {
            "description": "A powerful spring, now heals more!",
            "metadata": {
                "heal_amount": 75
            }
        }
        ```
      * **Response (Success - HTTP 200 OK):**
        ```json
        {
            "message": "POI updated successfully",
            "poi_id": "uuid-of-poi-1"
        }
        ```
      * **Response (Error - HTTP 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found):**
        ```json
        {
            "error": "POI not found or invalid data"
        }
        ```
      * **Gin Handler (Conceptual):** `poi.UpdatePOI`

5.  **Delete POI**

      * **Endpoint:** `/api/v1/pois/:poi_id`
      * **Method:** `DELETE`
      * **Purpose:** Removes a Point of Interest from the game world.
      * **Authentication:** Required (JWT). **Authorization:** Must have `admin` or `game_master` role.
      * **Response (Success - HTTP 204 No Content):** (No body, indicates successful deletion)
      * **Response (Error - HTTP 401 Unauthorized, 403 Forbidden, 404 Not Found):**
        ```json
        {
            "error": "POI not found or insufficient permissions"
        }
        ```
      * **Gin Handler (Conceptual):** `poi.DeletePOI`

-----

## IV. Real-time Communication Endpoint

This is not a traditional REST endpoint but the entry point for WebSocket connections.

1.  **WebSocket Connection for Map Updates**
      * **Endpoint:** `/ws/map-updates`
      * **Method:** `GET` (HTTP Upgrade request)
      * **Purpose:** Establishes a persistent bidirectional WebSocket connection for pushing real-time player location updates, and potentially new POI alerts, directly to connected clients.
      * **Authentication:** Required (JWT). The JWT can be sent as a query parameter or via a custom header during the WebSocket handshake.
      * **Messages Sent by Server (JSON):**
        ```json
        // Player location update
        {
            "type": "player_update",
            "player_id": "uuid-of-player-A",
            "username": "Alice",
            "latitude": 34.0535,
            "longitude": -118.2445,
            "timestamp": "2025-07-16T14:30:10Z"
        }
        // New POI notification
        {
            "type": "poi_created",
            "poi_id": "uuid-of-new-poi",
            "name": "New Resource Node",
            "latitude": 34.0550,
            "longitude": -118.2460,
            "poi_type": "resource_node"
        }
        ```
      * **Messages Sent by Client (Optional):** Clients could send messages to request specific updates or indicate their current viewport for optimized data.
        ```json
        {
            "type": "viewport_update",
            "min_lat": 34.00,
            "min_lon": -118.30,
            "max_lat": 34.10,
            "max_lon": -118.20
        }
        ```
      * **Gin Handler (Conceptual):** `websocket.HandleMapUpdates` (This handler would upgrade the HTTP connection to a WebSocket and manage the client connection).

-----

### General Considerations for Implementation:

  * **Versionining (`/v1`):** Using `/api/v1` in the URL path is a good practice for API versioning, allowing for future changes without breaking existing clients.
  * **Error Handling:** Consistent error response structures (e.g., `{"error": "message"}`) and appropriate HTTP status codes are essential for robust client-side development.
  * **Input Validation:** Every incoming request body and query parameter must be rigorously validated on the server-side (Gin's `ShouldBindJSON` helps, but further semantic validation is needed) to prevent invalid data or security vulnerabilities.
  * **Rate Limiting:** The `PUT /api/v1/players/:player_id/location` endpoint, in particular, will receive frequent requests. Implementing rate limiting is crucial to prevent abuse and ensure server stability.
  * **Security:** As discussed previously, JWT validation via Gin middleware will be applied to all authenticated routes. Role-based authorization will ensure that only authorized users can perform sensitive operations (e.g., POI creation/deletion).