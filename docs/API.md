# Uriel Virtual Office Platform - API Design

**Base URL:** `/api/v1`

## Overview

Uriel is a virtual office platform that enables distributed teams to collaborate, communicate, and feel connected in persistent digital spaces. This API supports real-time presence, proximity-based interactions, video conferencing, screen sharing, and seamless integrations with productivity tools.

---

## I. Authentication & User Management

### 1. User Registration
* **Endpoint:** `/api/v1/auth/register`
* **Method:** `POST`
* **Purpose:** Register a new user account
* **Authentication:** None (public)
* **Request Body:**
```json
{
    "email": "user@company.com",
    "username": "johndoe",
    "password": "secure_password",
    "full_name": "John Doe",
}
```
* **Response (Success - 201):**
```json
{
    "message": "User registered successfully",
    "user_id": "uuid-string",
    "workspace_id": "workspace-uuid"
}
```

### 2. User Login
* **Endpoint:** `/api/v1/auth/login`
* **Method:** `POST`
* **Purpose:** Authenticate user and get access token
* **Request Body:**
```json
{
    "email": "user@company.com",
    "password": "secure_password"
}
```
* **Response (Success - 200):**
```json
{
    "message": "Login successful",
    "access_token": "jwt-token",
    "refresh_token": "refresh-jwt",
    "user": {
        "user_id": "uuid-string",
        "username": "johndoe",
        "full_name": "John Doe",
        "email": "user@company.com",
        "avatar_url": "https://cdn.uriel.com/avatars/user.png",
        "role": "member",
        "workspace_id": "workspace-uuid"
    }
}
```

### 3. Token Refresh
* **Endpoint:** `/api/v1/auth/refresh`
* **Method:** `POST`
* **Purpose:** Refresh access token using refresh token
* **Request Body:**
```json
{
    "refresh_token": "refresh-jwt"
}
```

### 4. User Profile
* **Endpoint:** `/api/v1/users/profile`
* **Method:** `GET|PUT`
* **Purpose:** Get or update user profile
* **Authentication:** Required
* **PUT Request Body:**
```json
{
    "full_name": "John Smith",
    "avatar_url": "https://cdn.uriel.com/avatars/new-avatar.png",
    "status_message": "Working on Q4 planning",
    "timezone": "America/New_York"
}
```

---

## II. Workspace Management

### 1. Get Workspace Details
* **Endpoint:** `/api/v1/workspaces/:workspace_id`
* **Method:** `GET`
* **Purpose:** Get workspace information and settings
* **Response:**
```json
{
    "workspace_id": "workspace-uuid",
    "name": "Tech Corp Office",
    "description": "Our virtual headquarters",
    "settings": {
        "allow_guests": true,
        "require_approval": false,
        "default_room": "main-office",
        "max_users": 100
    },
    "owner_id": "owner-uuid",
    "created_at": "2025-01-01T00:00:00Z"
}
```

### 2. Workspace Members
* **Endpoint:** `/api/v1/workspaces/:workspace_id/members`
* **Method:** `GET`
* **Purpose:** List all workspace members
* **Query Parameters:**
  * `status`: Filter by online status (`online`, `offline`, `away`)
  * `room_id`: Filter by current room
* **Response:**
```json
{
    "members": [
        {
            "user_id": "uuid-1",
            "username": "johndoe",
            "full_name": "John Doe",
            "avatar_url": "https://cdn.uriel.com/avatars/john.png",
            "role": "admin",
            "status": "online",
            "current_room_id": "main-office",
            "position": {"x": 150, "y": 200},
            "last_seen": "2025-01-16T14:30:00Z"
        }
    ],
    "total_count": 25,
    "online_count": 8
}
```

### 3. Invite Users
* **Endpoint:** `/api/v1/workspaces/:workspace_id/invites`
* **Method:** `POST`
* **Purpose:** Send workspace invitations
* **Request Body:**
```json
{
    "invites": [
        {
            "email": "newuser@company.com",
            "role": "member",
            "personal_message": "Welcome to our team!"
        }
    ]
}
```

---

## III. Room Management

### 1. List Rooms
* **Endpoint:** `/api/v1/workspaces/:workspace_id/rooms`
* **Method:** `GET`
* **Purpose:** Get all rooms in workspace
* **Response:**
```json
{
    "rooms": [
        {
            "room_id": "main-office",
            "name": "Main Office",
            "type": "office",
            "capacity": 50,
            "current_users": 12,
            "background_url": "https://cdn.uriel.com/backgrounds/office.jpg",
            "is_private": false,
            "created_at": "2025-01-01T00:00:00Z"
        },
        {
            "room_id": "meeting-room-a",
            "name": "Conference Room A",
            "type": "meeting",
            "capacity": 10,
            "current_users": 0,
            "background_url": "https://cdn.uriel.com/backgrounds/meeting.jpg",
            "is_private": false
        }
    ]
}
```

### 2. Create Room
* **Endpoint:** `/api/v1/workspaces/:workspace_id/rooms`
* **Method:** `POST`
* **Purpose:** Create a new room
* **Request Body:**
```json
{
    "name": "Design Team Room",
    "type": "team",
    "capacity": 15,
    "background_url": "https://cdn.uriel.com/backgrounds/creative.jpg",
    "is_private": false,
    "description": "Dedicated space for design team collaboration"
}
```

### 3. Room Details
* **Endpoint:** `/api/v1/rooms/:room_id`
* **Method:** `GET`
* **Purpose:** Get detailed room information including current users
* **Response:**
```json
{
    "room_id": "main-office",
    "name": "Main Office",
    "type": "office",
    "capacity": 50,
    "background_url": "https://cdn.uriel.com/backgrounds/office.jpg",
    "current_users": [
        {
            "user_id": "uuid-1",
            "username": "johndoe",
            "full_name": "John Doe",
            "avatar_url": "https://cdn.uriel.com/avatars/john.png",
            "position": {"x": 150, "y": 200},
            "is_speaking": false,
            "is_sharing_screen": false,
            "joined_at": "2025-01-16T14:00:00Z"
        }
    ],
    "objects": [
        {
            "object_id": "whiteboard-1",
            "type": "whiteboard",
            "position": {"x": 300, "y": 100},
            "size": {"width": 200, "height": 150},
            "data": {"content_url": "https://cdn.uriel.com/whiteboards/board1.json"}
        }
    ]
}
```

### 4. Join Room
* **Endpoint:** `/api/v1/rooms/:room_id/join`
* **Method:** `POST`
* **Purpose:** Join a specific room
* **Request Body:**
```json
{
    "position": {"x": 100, "y": 150}
}
```

### 5. Leave Room
* **Endpoint:** `/api/v1/rooms/:room_id/leave`
* **Method:** `POST`
* **Purpose:** Leave current room

---

## IV. Real-time Presence & Movement

### 1. Update Position
* **Endpoint:** `/api/v1/users/position`
* **Method:** `PUT`
* **Purpose:** Update user's position in current room
* **Request Body:**
```json
{
    "room_id": "main-office",
    "position": {"x": 250, "y": 180},
    "facing_direction": "right"
}
```

### 2. Update Status
* **Endpoint:** `/api/v1/users/status`
* **Method:** `PUT`
* **Purpose:** Update user's availability status
* **Request Body:**
```json
{
    "status": "busy",
    "message": "In a meeting until 3 PM",
    "auto_expire_at": "2025-01-16T15:00:00Z"
}
```

---

## V. Communication & Meetings

### 1. Start Proximity Chat
* **Endpoint:** `/api/v1/communication/proximity/start`
* **Method:** `POST`
* **Purpose:** Initiate proximity-based voice chat
* **Request Body:**
```json
{
    "room_id": "main-office",
    "target_user_ids": ["uuid-1", "uuid-2"]
}
```

### 2. Create Meeting
* **Endpoint:** `/api/v1/meetings`
* **Method:** `POST`
* **Purpose:** Schedule or start an instant meeting
* **Request Body:**
```json
{
    "title": "Q4 Planning Review",
    "description": "Quarterly planning discussion",
    "room_id": "meeting-room-a",
    "scheduled_start": "2025-01-16T16:00:00Z",
    "duration_minutes": 60,
    "invitees": ["user-1", "user-2", "user-3"],
    "is_instant": false
}
```

### 3. Join Meeting
* **Endpoint:** `/api/v1/meetings/:meeting_id/join`
* **Method:** `POST`
* **Purpose:** Join an existing meeting
* **Response:**
```json
{
    "meeting_id": "meeting-uuid",
    "webrtc_config": {
        "ice_servers": [
            {"urls": "stun:stun.uriel.com:3478"},
            {"urls": "turn:turn.uriel.com:3478", "username": "user", "credential": "pass"}
        ]
    },
    "room_url": "https://meet.uriel.com/room/meeting-uuid"
}
```

---

## VI. Screen Sharing & Collaboration

### 1. Start Screen Share
* **Endpoint:** `/api/v1/collaboration/screen-share/start`
* **Method:** `POST`
* **Purpose:** Initiate screen sharing session
* **Request Body:**
```json
{
    "room_id": "main-office",
    "share_type": "window",
    "target_audience": "proximity"
}
```

### 2. Whiteboard Operations
* **Endpoint:** `/api/v1/collaboration/whiteboards/:whiteboard_id`
* **Method:** `GET|PUT|DELETE`
* **Purpose:** Manage whiteboard content
* **PUT Request Body:**
```json
{
    "content": {
        "elements": [
            {
                "type": "text",
                "position": {"x": 10, "y": 10},
                "content": "Project Timeline",
                "style": {"font_size": 16, "color": "#000000"}
            }
        ]
    }
}
```

---

## VII. Integrations

### 1. Calendar Integration
* **Endpoint:** `/api/v1/integrations/calendar/connect`
* **Method:** `POST`
* **Purpose:** Connect external calendar (Google, Outlook)
* **Request Body:**
```json
{
    "provider": "google",
    "auth_code": "google-oauth-code",
    "sync_settings": {
        "import_meetings": true,
        "create_room_for_meetings": true
    }
}
```

### 2. Slack Integration
* **Endpoint:** `/api/v1/integrations/slack/connect`
* **Method:** `POST`
* **Purpose:** Connect Slack workspace
* **Request Body:**
```json
{
    "workspace_id": "slack-workspace-id",
    "bot_token": "slack-bot-token",
    "settings": {
        "sync_status": true,
        "notifications": true,
        "channel_mapping": {
            "general": "main-office",
            "dev-team": "dev-room"
        }
    }
}
```

---

## VIII. WebSocket Endpoints

### 1. Real-time Updates
* **Endpoint:** `/ws/presence`
* **Method:** WebSocket Upgrade
* **Purpose:** Real-time presence, movement, and communication updates
* **Authentication:** JWT via query parameter or header

**Incoming Message Types:**
```json
// Position update
{
    "type": "position_update",
    "room_id": "main-office",
    "position": {"x": 150, "y": 200},
    "facing_direction": "right"
}

// Voice chat state
{
    "type": "voice_state",
    "is_muted": false,
    "is_speaking": true,
    "proximity_range": 50
}
```

**Outgoing Message Types:**
```json
// User joined room
{
    "type": "user_joined",
    "user": {
        "user_id": "uuid-1",
        "username": "johndoe",
        "avatar_url": "...",
        "position": {"x": 100, "y": 150}
    },
    "room_id": "main-office"
}

// User moved
{
    "type": "user_moved",
    "user_id": "uuid-1",
    "position": {"x": 200, "y": 180},
    "room_id": "main-office"
}

// Proximity chat started
{
    "type": "proximity_chat_started",
    "participants": ["uuid-1", "uuid-2"],
    "chat_id": "chat-uuid"
}
```

---

## IX. File & Asset Management

### 1. Upload Avatar
* **Endpoint:** `/api/v1/uploads/avatar`
* **Method:** `POST`
* **Content-Type:** `multipart/form-data`
* **Purpose:** Upload user avatar image

### 2. Upload Room Background
* **Endpoint:** `/api/v1/uploads/room-background`
* **Method:** `POST`
* **Content-Type:** `multipart/form-data`
* **Purpose:** Upload custom room background

### 3. File Sharing
* **Endpoint:** `/api/v1/files/share`
* **Method:** `POST`
* **Purpose:** Share files in room or meeting
* **Request Body:**
```json
{
    "file_url": "https://storage.uriel.com/files/document.pdf",
    "file_name": "Q4_Report.pdf",
    "room_id": "main-office",
    "share_type": "broadcast"
}
```

---

## X. Analytics & Reporting

### 1. Workspace Analytics
* **Endpoint:** `/api/v1/analytics/workspace/:workspace_id`
* **Method:** `GET`
* **Purpose:** Get workspace usage analytics
* **Query Parameters:**
  * `period`: `day`, `week`, `month`
  * `metric`: `active_users`, `meeting_time`, `room_usage`

### 2. User Activity
* **Endpoint:** `/api/v1/analytics/users/:user_id/activity`
* **Method:** `GET`
* **Purpose:** Get individual user activity metrics

---

## General API Conventions

### Authentication
- All authenticated endpoints require `Authorization: Bearer <jwt_token>` header
- WebSocket connections authenticate via `token` query parameter

### Rate Limiting
- Position updates: 10 requests per second
- API calls: 100 requests per minute per user
- File uploads: 5 per minute

### Error Responses
All errors follow this format:
```json
{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid position coordinates",
        "details": {
            "field": "position.x",
            "value": -100,
            "constraint": "must be non-negative"
        }
    }
}
```

### WebRTC Configuration
For voice/video features, the API provides WebRTC configuration:
```json
{
    "ice_servers": [
        {"urls": "stun:stun.uriel.com:3478"},
        {
            "urls": "turn:turn.uriel.com:3478",
            "username": "turn_user",
            "credential": "turn_password"
        }
    ]
}
```