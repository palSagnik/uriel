# MongoDB Database Schema Design for Interactive Game Map & Player Tracking System

This document outlines the comprehensive MongoDB database schema designed to support an immersive game world with multiple distinct maps, real-time player tracking, and points of interest. It covers the structure of `maps`, `players`, and `points_of_interest` collections, including fields, data types, and crucial indexing strategies for performance and scalability.

-----

## I. `maps` Collection

This collection defines each distinct map or instance within your game world.

**Purpose:** To categorize and manage separate playable areas, each potentially with its own characteristics, geographical boundaries, default spawn locations, and associated visual configurations.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f7a1a"), // MongoDB's default primary key (auto-generated)
    "map_id": "overworld-continent-alpha",      // A unique, application-generated identifier for the map (e.g., a slug or UUID)
    "name": "The Emerald Continent",            // Display name of the map
    "description": "A vast land of lush forests, towering mountains, and ancient ruins.", // Optional: Brief description
    "type": "overworld",                        // Categorization of the map (e.g., "overworld", "dungeon", "city", "instance")
    "boundaries": {                             // Optional: Defines the geographical extent of the map
        "type": "Polygon",
        "coordinates": [
            [
                [-120.0, 30.0],
                [-120.0, 40.0],
                [-110.0, 40.0],
                [-110.0, 30.0],
                [-120.0, 30.0]
            ]
        ]
    },
    "default_spawn_location": {                 // Optional: Where players might spawn when first entering this map
        "type": "Point",
        "coordinates": [-118.2437, 34.0522]     // [longitude, latitude]
    },
    "mapbox_style_id": "mapbox/streets-v11",    // Optional: Reference to a specific Mapbox style for this map
    "custom_tile_source_url": "https://tiles.yourgame.com/emerald_continent/{z}/{x}/{y}.png", // Optional: For custom raster/vector tile servers
    "created_at": ISODate("2025-01-01T00:00:00Z"), // Timestamp of map creation
    "updated_at": ISODate("2025-07-16T10:00:00Z")  // Timestamp of last map definition update
}
```

**Schema Fields and Data Types for `maps`:**

  * `_id`: `ObjectId` (MongoDB's default primary key)
  * `map_id`: `String` (Unique identifier for the map, e.g., a slug or UUID generated by the application)
  * `name`: `String` (Display name of the map)
  * `description`: `String` (Optional)
  * `type`: `String` (Categorization: `overworld`, `dungeon`, `city`, `instance`, etc.)
  * `boundaries`: `Object` (Optional: GeoJSON Polygon defining the traversable area of the map)
      * `type`: `String` (Must be "Polygon")
      * `coordinates`: `Array<Array<Array<Number>>>` (Array of linear rings)
  * `default_spawn_location`: `Object` (Optional: GeoJSON Point)
      * `type`: `String` (Must be "Point")
      * `coordinates`: `Array<Number>` (Contains `[longitude, latitude]`)
  * `mapbox_style_id`: `String` (Optional)
  * `custom_tile_source_url`: `String` (Optional)
  * `created_at`: `Date`
  * `updated_at`: `Date`

**Indexing Strategy for `maps`:**

  * `{ "map_id": 1 }`: Unique index for efficient lookup by map identifier.
  * `{ "name": 1 }`: Index for searching/listing maps by name.
  * `{ "type": 1 }`: Index for filtering maps by type.
  * `{ "boundaries": "2dsphere" }`: Optional geospatial index if you frequently query maps by geographical intersection.

-----

## II. `players` Collection

This collection stores information about each player, including their authentication credentials and their real-time location within a specific map.

**Purpose:** To manage player profiles and facilitate real-time location updates and queries, contextualized by their current map.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f6a7b"), // MongoDB's default primary key
    "player_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef", // A UUID for external system reference
    "username": "GamerTagPro",                          // Player's unique display name
    "email": "gamertagpro@example.com",                 // Player's email (optional)
    "password_hash": "$2a$10$abcdefghijklmnopqrstuvwxyza.1234567890abcdefgHIJKL...", // Hashed & salted password (e.g., bcrypt)
    "role": "player",                                   // e.g., "player", "admin", "game_master"
    "current_map_id": "overworld-continent-alpha",      // **NEW FIELD:** References `map_id` from the `maps` collection
    "location": {
        "type": "Point",
        "coordinates": [-118.2437, 34.0522]             // [longitude, latitude] as per GeoJSON standard, understood to be *within* `current_map_id`
    },
    "last_location_update_at": ISODate("2025-07-16T14:30:00Z"), // Last reported location timestamp
    "created_at": ISODate("2025-01-15T08:00:00Z"),         // Player account creation timestamp
    "updated_at": ISODate("2025-07-16T14:30:00Z"),         // Last update to player profile
    "is_online": true,                                  // Boolean: true if player is currently connected via WebSocket
    "session_id": "some-websocket-session-id"           // Optional: To map player to active WebSocket connection
    // Other player specific data can be added here, e.g., "inventory": [...], "stats": {...}
}
```

**Schema Fields and Data Types for `players`:**

  * `_id`: `ObjectId`
  * `player_id`: `String` (UUID generated by application for external use)
  * `username`: `String` (Unique display name)
  * `email`: `String` (Optional)
  * `password_hash`: `String`
  * `role`: `String`
  * **`current_map_id`: `String` (References `map_id` from the `maps` collection)**
  * `location`: `Object` (GeoJSON Point structure)
      * `type`: `String` (Must be "Point")
      * `coordinates`: `Array<Number>` (Contains `[longitude, latitude]`)
  * `last_location_update_at`: `Date`
  * `created_at`: `Date`
  * `updated_at`: `Date`
  * `is_online`: `Boolean`
  * `session_id`: `String` (Optional)

**Indexing Strategy for `players`:**

  * `{ "username": 1 }`: Unique index to ensure usernames are distinct.
  * `{ "email": 1 }`: Unique index (if `email` is enforced as unique) for faster lookup.
  * `{ "player_id": 1 }`: Index for fast lookups by application-generated UUID.
  * **`{ "current_map_id": 1, "location": "2dsphere" }`**: **Crucial compound geospatial index.** This enables efficient queries like "find all players *on a specific map* within a certain radius." It filters by map first, then applies the geospatial query.
  * `{ "current_map_id": 1, "is_online": 1 }`: Compound index for quickly filtering online players on a specific map.
  * `{ "last_location_update_at": -1 }`: Descending index for fetching recent updates or sorting players by last activity.

-----

## III. `points_of_interest` Collection

This collection stores all static and dynamic points of interest that appear on the game map, explicitly linked to a specific map.

**Purpose:** To manage and query various interactive or informational locations within the game world, contextualized by their associated map.

**Example Document Structure:**

```json
{
    "_id": ObjectId("60c72b2f9f1b2c3d4e5f6c8d"), // MongoDB's default primary key
    "poi_id": "f5e4d3c2-b1a0-9876-5432-1fedcba98765", // A UUID for external system reference
    "name": "Emerald Oasis",                         // Display name of the POI
    "type": "resource_node",                         // Category of the POI (e.g., "shop", "quest_giver")
    "map_id": "overworld-continent-alpha",           // **NEW FIELD:** References `map_id` from the `maps` collection
    "location": {
        "type": "Point",
        "coordinates": [-118.2550, 34.0450]          // [longitude, latitude], understood to be *within* `map_id`
    },
    "description": "A lush oasis where rare emeralds can be mined.",
    "metadata": {                                    // Flexible field for type-specific data (embedded document)
        "resource_type": "emerald",
        "respawn_time_seconds": 3600,
        "difficulty": "medium"
    },
    "is_active": true,                               // Boolean: Is this POI currently visible/active?
    "created_by": "admin_user_id_123",               // Optional: ID of the admin/GM who created it
    "created_at": ISODate("2025-02-01T12:00:00Z"),
    "updated_at": ISODate("2025-07-10T09:15:00Z")
}
```

**Schema Fields and Data Types for `points_of_interest`:**

  * `_id`: `ObjectId`
  * `poi_id`: `String` (UUID generated by application for external use)
  * `name`: `String`
  * `type`: `String`
  * **`map_id`: `String` (References `map_id` from the `maps` collection)**
  * `location`: `Object` (GeoJSON Point structure)
      * `type`: `String` (Must be "Point")
      * `coordinates`: `Array<Number>` (Contains `[longitude, latitude]`)
  * `description`: `String`
  * `metadata`: `Object` (Embedded document for arbitrary, type-specific data)
  * `is_active`: `Boolean`
  * `created_by`: `String` (Optional: References `player_id` of the creator)
  * `created_at`: `Date`
  * `updated_at`: `Date`

**Indexing Strategy for `points_of_interest`:**

  * `{ "poi_id": 1 }`: Index for fast lookups by application-generated UUID.
  * **`{ "map_id": 1, "location": "2dsphere" }`**: **Crucial compound geospatial index.** This enables efficient queries like "find all POIs *on a specific map* within a certain area."
  * `{ "map_id": 1, "type": 1 }`: Compound index for filtering POIs by type on a given map.
  * `{ "name": 1 }`: Index for fast lookups by name.
  * `{ "is_active": 1 }`: Index for filtering active/inactive POIs.

-----

## General Considerations for the Schema

  * **GeoJSON Standard:** All geographical coordinates (`location` fields) are stored in `[longitude, latitude]` array format, as specified by the GeoJSON standard. This is critical for MongoDB's geospatial query operators to function correctly.
  * **UUIDs for External IDs:** Using application-generated UUIDs (`player_id`, `poi_id`, `map_id`) alongside MongoDB's internal `_id` provides stable, globally unique identifiers. These are more convenient for inter-service communication, external APIs, and client-side references than exposing raw MongoDB ObjectIds.
  * **Denormalization:** The `metadata` field in `points_of_interest` is a prime example of denormalization. Instead of creating separate collections for each POI type, type-specific data is embedded directly within the POI document. This pattern is common in NoSQL databases to optimize read performance by reducing the need for joins.
  * **Compound Indexes:** The strategic use of compound indexes (e.g., `map_id` combined with `location`) is vital. By placing the `map_id` first in these indexes, queries can efficiently filter data down to the specific map before applying more complex geospatial calculations, significantly improving performance in a multi-map scenario.
  * **Time-to-Live (TTL) Indexes (Optional):** While not explicitly in the primary schema above, for auxiliary data like historical player movement trails (if stored separately) or temporary game events, TTL indexes can be used to automatically expire and remove old documents from a collection after a specified period, preventing unbounded data growth.
  * **Scalability Planning:** This schema design is inherently sharding-friendly. The `players` and `points_of_interest` collections, with their `map_id` and geospatial indexes, are well-prepared for a sharding strategy based on `map_id` (if queries primarily target a single map) or a custom shard key that considers the geospatial distribution (if queries frequently span map boundaries).
