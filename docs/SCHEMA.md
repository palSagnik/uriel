# Uriel Virtual Office Platform - Database Schema Design

This document outlines the comprehensive MongoDB database schema designed to support a virtual office collaboration platform similar to Gather.town. The schema supports workspaces, rooms, real-time user presence, meetings, and collaboration features.

---

## I. `workspaces` Collection

This collection defines organizational units that contain users, rooms, and settings.

**Purpose:** To manage different company/team workspaces, each with their own rooms, members, and configuration.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f7a1a"),
    "workspace_id": "tech-corp-hq",
    "name": "Tech Corp Virtual Office",
    "description": "Our main virtual headquarters",
    "owner_id": "owner-uuid-string",
    "settings": {
        "allow_guests": true,
        "require_approval": false,
        "default_room": "main-office",
        "max_users": 100,
        "working_hours": {
            "timezone": "America/New_York",
            "start": "09:00",
            "end": "17:00",
            "days": ["monday", "tuesday", "wednesday", "thursday", "friday"]
        }
    },
    "subscription": {
        "plan": "business",
        "max_users": 100,
        "features": ["screen_sharing", "integrations", "analytics"],
        "expires_at": ISODate("2025-12-31T23:59:59Z")
    },
    "created_at": ISODate("2025-01-01T00:00:00Z"),
    "updated_at": ISODate("2025-01-16T10:00:00Z")
}
```

**Schema Fields:**
- `_id`: `ObjectId` (MongoDB primary key)
- `workspace_id`: `String` (Unique identifier, slug format)
- `name`: `String` (Display name)
- `description`: `String` (Optional description)
- `owner_id`: `String` (References user_id from users collection)
- `settings`: `Object` (Workspace configuration)
- `subscription`: `Object` (Billing and feature limits)
- `created_at`: `Date`
- `updated_at`: `Date`

**Indexing Strategy:**
- `{ "workspace_id": 1 }`: Unique index
- `{ "owner_id": 1 }`: Index for owner lookups
- `{ "name": "text" }`: Text search index

---

## II. `users` Collection

This collection stores user profiles, authentication, and current presence information.

**Purpose:** To manage user accounts, authentication, real-time presence, and workspace membership.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f6a7b"),
    "user_id": "user-uuid-string",
    "email": "john.doe@techcorp.com",
    "username": "johndoe",
    "full_name": "John Doe",
    "password_hash": "$2a$10$abcdefghijklmnopqrstuvwxyza.1234567890abcdefgHIJKL...",
    "avatar_url": "https://cdn.uriel.com/avatars/john-doe.png",
    "role": "member",
    "workspace_id": "tech-corp-hq",
    "preferences": {
        "timezone": "America/New_York",
        "notification_settings": {
            "proximity_chat": true,
            "mentions": true,
            "meetings": true
        },
        "avatar_settings": {
            "virtual_background": "office",
            "show_name": true
        }
    },
    "presence": {
        "status": "online",
        "status_message": "Working on Q4 planning",
        "auto_expire_at": ISODate("2025-01-16T15:00:00Z"),
        "current_room_id": "main-office",
        "position": {
            "x": 150,
            "y": 200,
            "facing_direction": "right"
        },
        "is_speaking": false,
        "is_sharing_screen": false,
        "last_position_update": ISODate("2025-01-16T14:30:00Z")
    },
    "session": {
        "is_online": true,
        "session_id": "websocket-session-uuid",
        "last_seen": ISODate("2025-01-16T14:30:00Z"),
        "device_info": {
            "platform": "web",
            "browser": "Chrome",
            "version": "120.0.0"
        }
    },
    "created_at": ISODate("2025-01-15T08:00:00Z"),
    "updated_at": ISODate("2025-01-16T14:30:00Z")
}
```

**Schema Fields:**
- `_id`: `ObjectId`
- `user_id`: `String` (UUID for external reference)
- `email`: `String` (Unique email address)
- `username`: `String` (Unique username)
- `full_name`: `String`
- `password_hash`: `String` (bcrypt hash)
- `avatar_url`: `String` (URL to avatar image)
- `role`: `String` (member, admin, owner)
- `workspace_id`: `String` (References workspace)
- `preferences`: `Object` (User settings)
- `presence`: `Object` (Current status and location)
- `session`: `Object` (Connection information)
- `created_at`: `Date`
- `updated_at`: `Date`

**Indexing Strategy:**
- `{ "email": 1 }`: Unique index
- `{ "username": 1 }`: Unique index
- `{ "user_id": 1 }`: Index for UUID lookups
- `{ "workspace_id": 1, "presence.current_room_id": 1 }`: Compound index for room queries
- `{ "workspace_id": 1, "session.is_online": 1 }`: Index for online users
- `{ "workspace_id": 1, "role": 1 }`: Index for role-based queries

---

## III. `rooms` Collection

This collection defines virtual rooms/spaces within workspaces where users can gather.

**Purpose:** To manage different virtual spaces with their layouts, objects, and configurations.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f6c8d"),
    "room_id": "main-office",
    "workspace_id": "tech-corp-hq",
    "name": "Main Office",
    "description": "Our central workspace for daily collaboration",
    "type": "office",
    "capacity": 50,
    "is_private": false,
    "background": {
        "type": "image",
        "url": "https://cdn.uriel.com/backgrounds/modern-office.jpg",
        "dimensions": {
            "width": 1200,
            "height": 800
        }
    },
    "objects": [
        {
            "object_id": "whiteboard-1",
            "type": "whiteboard",
            "position": {"x": 300, "y": 100},
            "size": {"width": 200, "height": 150},
            "data": {
                "content_url": "https://cdn.uriel.com/whiteboards/board1.json",
                "permissions": ["read", "write"]
            }
        },
        {
            "object_id": "meeting-table-1",
            "type": "meeting_area",
            "position": {"x": 500, "y": 300},
            "size": {"width": 150, "height": 100},
            "data": {
                "max_participants": 8,
                "auto_start_voice": true
            }
        }
    ],
    "spawn_points": [
        {"x": 100, "y": 150},
        {"x": 120, "y": 150},
        {"x": 140, "y": 150}
    ],
    "settings": {
        "voice_enabled": true,
        "screen_sharing_enabled": true,
        "max_voice_distance": 100,
        "background_music": null
    },
    "created_by": "owner-user-id",
    "created_at": ISODate("2025-01-01T00:00:00Z"),
    "updated_at": ISODate("2025-01-10T09:15:00Z")
}
```

**Schema Fields:**
- `_id`: `ObjectId`
- `room_id`: `String` (Unique identifier within workspace)
- `workspace_id`: `String` (References workspace)
- `name`: `String`
- `description`: `String`
- `type`: `String` (office, meeting, social, private)
- `capacity`: `Number`
- `is_private`: `Boolean`
- `background`: `Object` (Background configuration)
- `objects`: `Array` (Interactive objects in the room)
- `spawn_points`: `Array` (Entry positions)
- `settings`: `Object` (Room-specific settings)
- `created_by`: `String` (References user_id)
- `created_at`: `Date`
- `updated_at`: `Date`

**Indexing Strategy:**
- `{ "workspace_id": 1, "room_id": 1 }`: Compound unique index
- `{ "workspace_id": 1, "type": 1 }`: Index for filtering by room type
- `{ "workspace_id": 1, "is_private": 1 }`: Index for public/private rooms

---

## IV. `meetings` Collection

This collection manages scheduled and instant meetings within the platform.

**Purpose:** To handle meeting scheduling, participants, and integration with calendar systems.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f6d9e"),
    "meeting_id": "meeting-uuid-string",
    "workspace_id": "tech-corp-hq",
    "title": "Q4 Planning Review",
    "description": "Quarterly planning and review session",
    "organizer_id": "user-uuid-string",
    "room_id": "meeting-room-a",
    "scheduled_start": ISODate("2025-01-16T16:00:00Z"),
    "scheduled_end": ISODate("2025-01-16T17:00:00Z"),
    "actual_start": ISODate("2025-01-16T16:02:00Z"),
    "actual_end": null,
    "status": "in_progress",
    "type": "scheduled",
    "participants": [
        {
            "user_id": "user-1",
            "status": "accepted",
            "joined_at": ISODate("2025-01-16T16:02:00Z"),
            "left_at": null,
            "role": "organizer"
        },
        {
            "user_id": "user-2",
            "status": "accepted",
            "joined_at": ISODate("2025-01-16T16:05:00Z"),
            "left_at": null,
            "role": "participant"
        }
    ],
    "settings": {
        "recording_enabled": false,
        "auto_admit": true,
        "waiting_room": false,
        "screen_sharing_allowed": true
    },
    "integrations": {
        "calendar_event_id": "google-calendar-event-id",
        "slack_channel": "general"
    },
    "created_at": ISODate("2025-01-15T10:00:00Z"),
    "updated_at": ISODate("2025-01-16T16:05:00Z")
}
```

**Schema Fields:**
- `_id`: `ObjectId`
- `meeting_id`: `String` (UUID)
- `workspace_id`: `String`
- `title`: `String`
- `description`: `String`
- `organizer_id`: `String`
- `room_id`: `String`
- `scheduled_start`: `Date`
- `scheduled_end`: `Date`
- `actual_start`: `Date`
- `actual_end`: `Date`
- `status`: `String` (scheduled, in_progress, completed, cancelled)
- `type`: `String` (scheduled, instant)
- `participants`: `Array`
- `settings`: `Object`
- `integrations`: `Object`
- `created_at`: `Date`
- `updated_at`: `Date`

**Indexing Strategy:**
- `{ "meeting_id": 1 }`: Unique index
- `{ "workspace_id": 1, "scheduled_start": 1 }`: Index for upcoming meetings
- `{ "workspace_id": 1, "organizer_id": 1 }`: Index for user's meetings
- `{ "workspace_id": 1, "status": 1 }`: Index for filtering by status

---

## V. `integrations` Collection

This collection manages external service integrations for workspaces.

**Purpose:** To store configuration and credentials for third-party integrations like Slack, Google Calendar, etc.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f6e0f"),
    "integration_id": "integration-uuid-string",
    "workspace_id": "tech-corp-hq",
    "type": "slack",
    "name": "Tech Corp Slack",
    "status": "active",
    "config": {
        "workspace_name": "techcorp",
        "team_id": "T1234567890",
        "bot_user_id": "U1234567890"
    },
    "credentials": {
        "access_token": "encrypted-token",
        "refresh_token": "encrypted-refresh-token",
        "expires_at": ISODate("2025-06-01T00:00:00Z")
    },
    "settings": {
        "sync_status": true,
        "notifications": true,
        "channel_mapping": {
            "general": "main-office",
            "dev-team": "dev-room",
            "design": "design-room"
        },
        "auto_create_rooms": false
    },
    "last_sync": ISODate("2025-01-16T14:00:00Z"),
    "created_by": "user-uuid-string",
    "created_at": ISODate("2025-01-01T00:00:00Z"),
    "updated_at": ISODate("2025-01-16T14:00:00Z")
}
```

**Schema Fields:**
- `_id`: `ObjectId`
- `integration_id`: `String` (UUID)
- `workspace_id`: `String`
- `type`: `String` (slack, google_calendar, outlook, etc.)
- `name`: `String`
- `status`: `String` (active, inactive, error)
- `config`: `Object` (Integration-specific configuration)
- `credentials`: `Object` (Encrypted authentication data)
- `settings`: `Object` (User-defined settings)
- `last_sync`: `Date`
- `created_by`: `String`
- `created_at`: `Date`
- `updated_at`: `Date`

**Indexing Strategy:**
- `{ "integration_id": 1 }`: Unique index
- `{ "workspace_id": 1, "type": 1 }`: Compound index for workspace integrations
- `{ "workspace_id": 1, "status": 1 }`: Index for active integrations

---

## VI. `activities` Collection

This collection tracks user activities and events for analytics and notifications.

**Purpose:** To log user interactions, presence changes, and system events for analytics and audit trails.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f6f10"),
    "activity_id": "activity-uuid-string",
    "workspace_id": "tech-corp-hq",
    "user_id": "user-uuid-string",
    "type": "room_joined",
    "details": {
        "room_id": "main-office",
        "room_name": "Main Office",
        "previous_room_id": "meeting-room-a",
        "position": {"x": 150, "y": 200}
    },
    "metadata": {
        "session_id": "websocket-session-uuid",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0..."
    },
    "timestamp": ISODate("2025-01-16T14:30:00Z")
}
```

**Activity Types:**
- `user_login`, `user_logout`
- `room_joined`, `room_left`
- `meeting_started`, `meeting_joined`, `meeting_left`
- `screen_share_started`, `screen_share_stopped`
- `position_updated`
- `status_changed`

**Indexing Strategy:**
- `{ "workspace_id": 1, "timestamp": -1 }`: Index for recent activities
- `{ "workspace_id": 1, "user_id": 1, "timestamp": -1 }`: Index for user activities
- `{ "workspace_id": 1, "type": 1, "timestamp": -1 }`: Index for activity types
- TTL Index: `{ "timestamp": 1 }` with expiration (e.g., 90 days)

---

## General Schema Considerations

### Data Types and Conventions
- **UUIDs**: All external identifiers use UUID strings for security and interoperability
- **Timestamps**: All dates stored in UTC using MongoDB's ISODate format
- **Coordinates**: Position data uses simple {x, y} objects for 2D room positioning
- **Encryption**: Sensitive data like integration credentials are encrypted at rest

### Performance Optimizations
- **Compound Indexes**: Strategic use of workspace_id as the first field in most indexes
- **Embedded Documents**: User presence and session data embedded to avoid joins
- **TTL Indexes**: Automatic cleanup of old activity logs and temporary data

### Scalability Planning
- **Sharding Strategy**: Primary shard key on workspace_id for horizontal scaling
- **Read Replicas**: Analytics queries can be directed to read-only replicas
- **Archival Strategy**: Old meetings and activities archived to separate collections

### Security Considerations
- **Field-Level Encryption**: Sensitive data encrypted using MongoDB's FLE
- **Access Control**: Role-based access control at the application level
- **Audit Logging**: All administrative actions logged to activities collection

### Real-time Data Handling
- **Change Streams**: MongoDB change streams used for real-time WebSocket updates
- **Optimistic Locking**: Version fields for handling concurrent updates
- **Presence Management**: Efficient updates for high-frequency position data
