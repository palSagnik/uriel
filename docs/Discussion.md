# Uriel Virtual Office Platform - Routes & Architecture

This document provides detailed route specifications, architectural decisions, and implementation considerations for the Uriel virtual office platform.

---

## Route Specifications by Module

### Authentication Module (`/api/v1/auth`)

```
POST   /register           - User registration
POST   /login              - User authentication  
POST   /refresh            - Token refresh
POST   /logout             - User logout
POST   /forgot-password    - Password reset request
POST   /reset-password     - Password reset confirmation
GET    /verify-email       - Email verification
```

**Rate Limiting:**
- Login attempts: 5 per minute per IP
- Registration: 3 per minute per IP
- Password reset: 1 per minute per email

### User Management (`/api/v1/users`)

```
GET    /profile                    - Get current user profile
PUT    /profile                    - Update user profile
POST   /profile/avatar            - Upload avatar image
GET    /profile/activity          - Get user activity history
PUT    /status                    - Update presence status
PUT    /position                  - Update position in room
GET    /settings                  - Get user preferences
PUT    /settings                  - Update user preferences
DELETE /account                   - Delete user account
```

### Workspace Management (`/api/v1/workspaces`)

```
POST   /                          - Create new workspace
GET    /:workspace_id             - Get workspace details
PUT    /:workspace_id             - Update workspace settings
DELETE /:workspace_id             - Delete workspace
GET    /:workspace_id/members     - List workspace members
POST   /:workspace_id/invites     - Send invitations
GET    /:workspace_id/invites     - List pending invitations
DELETE /:workspace_id/invites/:id - Cancel invitation
POST   /:workspace_id/members/:user_id/role - Update member role
DELETE /:workspace_id/members/:user_id      - Remove member
```

### Room Management (`/api/v1/rooms`, `/api/v1/workspaces/:workspace_id/rooms`)

```
GET    /                     - List all accessible rooms
POST   /                     - Create new room
GET    /:room_id            - Get room details
PUT    /:room_id            - Update room settings
DELETE /:room_id            - Delete room
POST   /:room_id/join       - Join room
POST   /:room_id/leave      - Leave room
GET    /:room_id/users      - Get users in room
POST   /:room_id/objects    - Add object to room
PUT    /:room_id/objects/:object_id    - Update room object
DELETE /:room_id/objects/:object_id    - Remove room object
POST   /:room_id/background             - Update room background
```

### Communication (`/api/v1/communication`)

```
POST   /proximity/start             - Start proximity chat
POST   /proximity/stop              - End proximity chat
GET    /proximity/active            - Get active proximity chats
POST   /voice/mute                  - Mute/unmute microphone
POST   /voice/deafen                - Deafen/undeafen audio
GET    /voice/status                - Get voice connection status
```

### Meeting Management (`/api/v1/meetings`)

```
GET    /                     - List meetings (upcoming, ongoing)
POST   /                     - Create/schedule meeting
GET    /:meeting_id         - Get meeting details
PUT    /:meeting_id         - Update meeting
DELETE /:meeting_id         - Cancel meeting
POST   /:meeting_id/join    - Join meeting
POST   /:meeting_id/leave   - Leave meeting
POST   /:meeting_id/start   - Start scheduled meeting
POST   /:meeting_id/end     - End ongoing meeting
GET    /:meeting_id/participants    - Get meeting participants
POST   /:meeting_id/invite          - Invite users to meeting
```

### Screen Sharing & Collaboration (`/api/v1/collaboration`)

```
POST   /screen-share/start          - Start screen sharing
POST   /screen-share/stop           - Stop screen sharing
GET    /screen-share/active         - Get active screen shares
POST   /whiteboards                 - Create whiteboard
GET    /whiteboards/:id            - Get whiteboard content
PUT    /whiteboards/:id            - Update whiteboard
DELETE /whiteboards/:id            - Delete whiteboard
POST   /whiteboards/:id/collaborate - Join whiteboard session
GET    /files/shared               - List shared files
POST   /files/share                - Share file
DELETE /files/:id                  - Remove shared file
```

### Integrations (`/api/v1/integrations`)

```
GET    /                           - List workspace integrations
POST   /calendar/connect           - Connect calendar service
POST   /slack/connect              - Connect Slack workspace
POST   /google/connect             - Connect Google Workspace
POST   /outlook/connect            - Connect Outlook
GET    /:integration_id           - Get integration details
PUT    /:integration_id           - Update integration settings
DELETE /:integration_id           - Remove integration
POST   /:integration_id/sync      - Manual sync trigger
GET    /:integration_id/logs      - Get sync logs
```

### File Management (`/api/v1/uploads`, `/api/v1/files`)

```
POST   /uploads/avatar             - Upload user avatar
POST   /uploads/room-background    - Upload room background
POST   /uploads/workspace-logo     - Upload workspace logo
POST   /files                      - Upload file for sharing
GET    /files/:id                  - Download file
DELETE /files/:id                  - Delete file
GET    /files/:id/metadata        - Get file metadata
```

### Analytics & Reporting (`/api/v1/analytics`)

```
GET    /workspace/:workspace_id    - Workspace analytics
GET    /users/:user_id/activity   - User activity analytics
GET    /rooms/:room_id/usage      - Room usage statistics
GET    /meetings/summary          - Meeting summary statistics
GET    /engagement/daily          - Daily engagement metrics
GET    /engagement/weekly         - Weekly engagement metrics
```

---

## WebSocket Event Specifications

### Connection Endpoint: `/ws/presence`

**Authentication:** JWT token via query parameter or header

#### Incoming Events (Client → Server)

```javascript
// Position update
{
    "type": "position_update",
    "room_id": "main-office",
    "position": {"x": 150, "y": 200},
    "facing_direction": "right",
    "timestamp": "2025-01-16T14:30:00Z"
}

// Voice state change
{
    "type": "voice_state",
    "is_muted": false,
    "is_speaking": true,
    "proximity_range": 50
}

// Room join/leave
{
    "type": "room_action",
    "action": "join|leave",
    "room_id": "main-office",
    "position": {"x": 100, "y": 150}
}

// Typing indicator
{
    "type": "typing",
    "is_typing": true,
    "context": "chat|whiteboard",
    "context_id": "whiteboard-123"
}

// Interaction with room objects
{
    "type": "object_interaction",
    "object_id": "whiteboard-1",
    "action": "start_editing|stop_editing",
    "data": {...}
}
```

#### Outgoing Events (Server → Client)

```javascript
// User presence updates
{
    "type": "user_joined",
    "user": {
        "user_id": "uuid-1",
        "username": "johndoe",
        "full_name": "John Doe",
        "avatar_url": "...",
        "position": {"x": 100, "y": 150},
        "status": "online"
    },
    "room_id": "main-office"
}

{
    "type": "user_left",
    "user_id": "uuid-1",
    "room_id": "main-office"
}

{
    "type": "user_moved",
    "user_id": "uuid-1",
    "position": {"x": 200, "y": 180},
    "facing_direction": "left",
    "room_id": "main-office"
}

// Voice/communication events
{
    "type": "proximity_chat_started",
    "participants": ["uuid-1", "uuid-2"],
    "chat_id": "chat-uuid",
    "position": {"x": 150, "y": 200}
}

{
    "type": "voice_state_changed",
    "user_id": "uuid-1",
    "is_muted": true,
    "is_speaking": false
}

// Screen sharing events
{
    "type": "screen_share_started",
    "user_id": "uuid-1",
    "stream_id": "stream-uuid",
    "room_id": "main-office"
}

// Meeting events
{
    "type": "meeting_started",
    "meeting_id": "meeting-uuid",
    "room_id": "meeting-room-a",
    "organizer_id": "uuid-1"
}

// Whiteboard collaboration
{
    "type": "whiteboard_updated",
    "whiteboard_id": "board-uuid",
    "changes": [...],
    "user_id": "uuid-1"
}

// System notifications
{
    "type": "notification",
    "category": "meeting|mention|system",
    "title": "Meeting Starting Soon",
    "message": "Q4 Planning starts in 5 minutes",
    "action_url": "/meetings/meeting-uuid"
}
```

---

## Error Handling Specifications

### Standard Error Response Format

```json
{
    "error": {
        "code": "ERROR_CODE",
        "message": "Human-readable error message",
        "details": {
            "field": "field_name",
            "value": "invalid_value",
            "constraint": "validation_rule"
        },
        "request_id": "req_uuid",
        "timestamp": "2025-01-16T14:30:00Z"
    }
}
```

### Common Error Codes

```
VALIDATION_ERROR       - Input validation failed
AUTHENTICATION_FAILED  - Invalid credentials
AUTHORIZATION_FAILED   - Insufficient permissions
RESOURCE_NOT_FOUND     - Requested resource doesn't exist
RESOURCE_CONFLICT      - Resource already exists or conflicting state
RATE_LIMIT_EXCEEDED    - Too many requests
WORKSPACE_FULL         - Workspace at capacity
ROOM_FULL             - Room at capacity  
MEETING_NOT_ACTIVE    - Meeting not in progress
INTEGRATION_ERROR     - External service error
UPLOAD_SIZE_EXCEEDED  - File too large
INVALID_FILE_TYPE     - Unsupported file format
```

---

## Authentication & Authorization

### JWT Token Structure

```json
{
    "user_id": "uuid-string",
    "username": "johndoe",
    "workspace_id": "workspace-uuid",
    "role": "member|admin|owner",
    "permissions": ["read_workspace", "write_rooms", "manage_users"],
    "iat": 1642348800,
    "exp": 1642352400
}
```

### Permission Levels

**Member:**
- Join public rooms
- Create personal rooms
- Participate in meetings
- Use integrations

**Admin:**
- All member permissions
- Manage workspace settings
- Create/delete rooms
- Invite/remove users
- View analytics

**Owner:**
- All admin permissions
- Delete workspace
- Manage billing
- Transfer ownership

---

## Performance Considerations

### Caching Strategy

```
Redis Cache Layers:
- User sessions: 24 hours TTL
- Room state: 1 hour TTL, real-time invalidation
- Workspace metadata: 6 hours TTL
- Integration tokens: Token expiry TTL
- Analytics data: 30 minutes TTL
```

### Database Query Optimization

```javascript
// Optimized room user query
db.users.find({
    "workspace_id": "workspace-uuid",
    "presence.current_room_id": "room-id",
    "session.is_online": true
}).limit(50)

// Efficient position updates (upsert)
db.users.updateOne(
    {"user_id": "user-uuid"},
    {
        "$set": {
            "presence.position": {"x": 150, "y": 200},
            "presence.last_position_update": new Date()
        }
    }
)
```

### WebSocket Connection Management

```
Connection Pooling:
- Max connections per workspace: 1000
- Connection timeout: 30 seconds
- Heartbeat interval: 30 seconds
- Reconnection attempts: 3 with exponential backoff

Message Queuing:
- Position updates: Throttled to 10/second per user
- Voice state: Immediate delivery
- Presence updates: Batched every 2 seconds
```

---

## Integration Specifications

### Slack Integration Flow

1. **OAuth Setup:**
   ```
   POST /api/v1/integrations/slack/connect
   → Redirect to Slack OAuth
   → Callback with authorization code
   → Exchange for access token
   → Store encrypted credentials
   ```

2. **Status Sync:**
   ```
   User status in Uriel ↔ Slack status
   - Online → 🟢 In virtual office
   - Busy → 🔴 In meeting
   - Away → 🟡 Away from desk
   ```

3. **Channel Mapping:**
   ```json
   {
       "channel_mappings": {
           "#general": "main-office",
           "#dev-team": "development-room", 
           "#design": "design-studio"
       }
   }
   ```

### Calendar Integration Flow

1. **Google Calendar:**
   ```
   - Sync meetings as room events
   - Auto-create meeting rooms
   - Send room links in calendar invites
   - Import busy/free status
   ```

2. **Outlook Integration:**
   ```
   - Similar to Google Calendar
   - Teams meeting links converted to Uriel rooms
   - Exchange availability data
   ```

---

## Security Specifications

### Data Protection

```
Encryption:
- Passwords: bcrypt with salt rounds 12
- API tokens: AES-256 encryption
- Database: Field-level encryption for PII
- Transport: TLS 1.3 for all connections

Privacy:
- Position data retention: 30 days
- Activity logs: 90 days
- Deleted user data: 7 days grace period
- GDPR compliance: Data export/deletion APIs
```

### Rate Limiting

```yaml
Global Limits:
  requests_per_minute: 100
  requests_per_hour: 5000

Endpoint Specific:
  auth_login: 5/minute
  position_update: 600/minute  
  file_upload: 5/minute
  meeting_create: 10/hour

WebSocket:
  connections_per_user: 3
  messages_per_second: 10
  position_updates_per_second: 10
```

This comprehensive route and architecture specification provides the foundation for implementing a robust virtual office platform similar to Gather.town, with detailed API endpoints, real-time communication protocols, and production-ready considerations.

