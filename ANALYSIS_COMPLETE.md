# Backend vs Frontend Models Analysis

## Overview
This document provides a comprehensive comparison between backend domain models (Go) and frontend entity types (TypeScript).

---

## ✅ Models Already Aligned

### 1. **CrosshairSettings**
- **Backend**: `backend/internal/domain/crosshair.go`
- **Frontend**: `frontend/src/entities/crosshair/types/types.ts`
- **Status**: ✅ Fully aligned

### 2. **DeadlockMMR** 
- **Backend**: `backend/internal/domain/deadlock_mmr.go`
- **Frontend**: `frontend/src/entities/deadlock/types/types.ts`
- **Status**: ✅ Fully aligned

### 3. **PerformanceDynamics & Trend**
- **Backend**: `backend/internal/domain/player_profile.go`
- **Frontend**: `frontend/src/entities/player/types/types.ts` & `deadlock/types/types.ts`
- **Status**: ✅ Fully aligned

### 4. **PersonalRecords**
- **Backend**: `backend/internal/domain/personal_records.go`
- **Frontend**: `frontend/src/entities/deadlock/types/types.ts`
- **Status**: ✅ Fully aligned

### 5. **MateStat**
- **Backend**: `backend/internal/domain/mate_stat.go`
- **Frontend**: `frontend/src/entities/deadlock/types/types.ts`
- **Status**: ✅ Fully aligned

### 6. **FeaturedHero**
- **Backend**: `backend/internal/domain/featured_hero.go`
- **Frontend**: `frontend/src/entities/deadlock/types/types.ts`
- **Status**: ✅ Fully aligned

### 7. **HeroMMRHistory**
- **Backend**: `backend/internal/domain/deadlock_mmr.go`
- **Frontend**: `frontend/src/entities/deadlock/types/types.ts`
- **Status**: ⚠️ Needs alignment - frontend uses nested object structure

---

## ⚠️ Models Needing Alignment

### 1. **Crosshair**
**Backend** (`backend/internal/domain/crosshair.go`):
```go
type Crosshair struct {
    ID          uuid.UUID       `json:"id"`
    AuthorID    uuid.UUID       `json:"author_id"`
    Author      *User           `json:"author,omitempty"`
    Title       string          `json:"title"`
    Description string          `json:"description"`
    Settings    json.RawMessage `json:"settings"`
    LikesCount  int             `json:"likes_count"`
    IsPublic    bool            `json:"is_public"`
    ViewCount   int             `json:"view_count"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
}
```

**Frontend** (`frontend/src/entities/crosshair/types/types.ts`):
```typescript
// Has two different interfaces:
interface CrosshairListItem {
  id: string
  title: string
  description: string
  settings: CrosshairSettings
  likes_count: number
  created_at: string
  author_id: string
  author_name?: string
  author_avatar?: string
  is_liked: boolean  // ❌ Not in backend
  is_public: boolean
  view_count: number
}

interface PublishedCrosshair {  // ❌ Should be removed
  id: string
  settings: CrosshairSettings
  likes: number
  author_id: string
  createdAt: string  // ❌ Wrong casing
}
```

**Issues**:
- Frontend has `author_name` and `author_avatar` (should use nested `author` object)
- Frontend has `is_liked` field (client-side only)
- `PublishedCrosshair` interface is redundant
- Missing `updated_at` field
- Inconsistent date field naming

---

### 2. **CrosshairLike**
**Backend** (`backend/internal/domain/crosshair.go`):
```go
type CrosshairLike struct {
    ID          uuid.UUID `json:"id"`
    UserID      uuid.UUID `json:"user_id"`
    CrosshairID uuid.UUID `json:"crosshair_id"`
    CreatedAt   time.Time `json:"created_at"`
}
```

**Frontend**: ❌ **Missing entirely**

---

### 3. **Match**
**Backend** (`backend/internal/domain/match.go`):
```go
type Match struct {
    ID                   string    `json:"match_id"`  // ⚠️ "match_id" not "id"
    HeroID               int       `json:"hero_id"`
    PlayerKills          int       `json:"player_kills"`
    PlayerDeaths         int       `json:"player_deaths"`
    PlayerAssists        int       `json:"player_assists"`
    NetWorth             int       `json:"net_worth"`
    MatchDurationS       int       `json:"match_duration_s"`
    MatchResult          int       `json:"match_result"`
    PlayerTeam           int       `json:"player_team"`
    StartTime            int64     `json:"start_time"`
    HeroName             string    `json:"hero_name"`
    HeroAvatar           string    `json:"hero_avatar,omitempty"`
    PlayerRankAfterMatch int       `json:"player_rank_after_match"`
    RankName             string    `json:"rank_name"`
    SubRank              int       `json:"sub_rank"`
    RankImage            string    `json:"rank_image"`
    PlayerRankChange     int       `json:"player_rank_change"`
    Kills                int       `json:"kills,omitempty"`
    Deaths               int       `json:"deaths,omitempty"`
    Assists              int       `json:"assists,omitempty"`
    DurationMinutes      int       `json:"duration_minutes,omitempty"`
    MatchTime            time.Time `json:"match_time,omitempty"`
    Result               string    `json:"result"`
}
```

**Frontend** (`frontend/src/entities/player/types/types.ts`):
```typescript
interface Match {
  id: string  // ❌ Should be "match_id"
  hero_name: string
  hero_avatar: string
  result: 'Win' | 'Loss'
  player_kills: number
  player_deaths: number
  player_assists: number
  match_duration_s: number
  player_rank_change: number
  player_rank_after_match: number
  rank_name: string
  sub_rank: number | undefined
  rank_image: string | undefined
  match_time: string
  souls: number  // ❌ Not in backend Match model
  player_score: number  // ❌ Not in backend Match model
}
```

**Issues**:
- Frontend uses `id` instead of `match_id`
- Frontend missing many fields: `hero_id`, `net_worth`, `match_result`, `player_team`, `start_time`, etc.
- Frontend has `souls` and `player_score` which aren't in backend Match
- Frontend has `result` as enum, backend has both `result` string and `match_result` int

---

### 4. **HeroStat**
**Backend** (`backend/internal/domain/hero_stat.go`):
```go
type HeroStat struct {
    HeroID     int     `json:"hero_id"`
    HeroName   string  `json:"hero_name"`
    Matches    int     `json:"matches_played"`  // ⚠️ "matches_played"
    WinRate    float64 `json:"win_rate"`
    KDA        float64 `json:"kda"`
    HeroAvatar string  `json:"hero_avatar,omitempty"`
}
```

**Frontend** (`frontend/src/entities/player/types/types.ts`):
```typescript
interface HeroStat {
  hero_name: string
  matches: number  // ❌ Should be "matches_played"
  win_rate: number
  kda: number
  hero_avatar?: string
}
```

**Issues**:
- Missing `hero_id` field
- `matches` should be `matches_played`

---

### 5. **User**
**Backend** (`backend/internal/domain/user.go`):
```go
type User struct {
    ID         uuid.UUID `json:"id"`
    SteamID    string    `json:"steam_id"`
    Nickname   string    `json:"nickname"`
    AvatarURL  string    `json:"avatar_url"`
    ProfileURL string    `json:"profile_url"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

**Backend DTO** (`backend/internal/dto/user_search_result.go`):
```go
type UserSearchResult struct {
    ID         string     `json:"id"`
    SteamID    string     `json:"steam_id"`
    Nickname   string     `json:"nickname"`
    AvatarURL  string     `json:"avatar_url"`
    ProfileURL string     `json:"profile_url"`
    CreatedAt  *time.Time `json:"created_at,omitempty"`
    UpdatedAt  *time.Time `json:"updated_at,omitempty"`
    
    AccountID   int    `json:"account_id,omitempty"`
    CountryCode string `json:"countrycode,omitempty"`
    LastUpdated int64  `json:"last_updated,omitempty"`
    Realname    string `json:"realname,omitempty"`
    
    IsDeadlockPlayer    bool `json:"is_deadlock_player"`
    DeadlockStatusKnown bool `json:"deadlock_status_known"`
}
```

**Frontend** (`frontend/src/entities/user/types/types.ts`):
```typescript
interface User {
  id: string
  steam_id: string
  nickname: string
  avatar_url: string
  profile_url: string
  created_at?: Date  // ❌ Should be string (ISO format)
  updated_at?: Date  // ❌ Should be string (ISO format)
  
  account_id?: number
  countrycode?: string
  last_updated?: number
  realname?: string
  
  is_deadlock_player: boolean
  deadlock_status_known: boolean
}
```

**Issues**:
- Frontend uses `Date` type but backend sends ISO string
- Frontend matches DTO structure more than domain model

---

### 6. **PlayerProfile**
**Backend** (`backend/internal/domain/player_profile.go`):
```go
type PlayerProfile struct {
    SteamID               string              `json:"steam_id"`
    Nickname              string              `json:"nickname"`
    AvatarURL             string              `json:"avatar_url"`
    ProfileURL            string              `json:"profile_url"`
    CreatedAt             time.Time           `json:"created_at"`
    LastMatchTime         time.Time           `json:"last_match_time"`
    PlayerRank            int                 `json:"player_rank"`
    RankName              string              `json:"rank_name"`
    SubRank               int                 `json:"sub_rank"`
    RankImage             string              `json:"rank_image"`
    WinRate               float64             `json:"win_rate"`
    KDRatio               float64             `json:"kd_ratio"`
    AvgMatchesPerDay      float64             `json:"avg_matches_per_day"`
    FavoriteHero          string              `json:"favorite_hero"`
    LastUpdatedAt         time.Time           `json:"last_updated_at"`
    TotalMatches          int                 `json:"total_matches"`
    TotalKills            int                 `json:"total_kills"`
    TotalDeaths           int                 `json:"total_deaths"`
    TotalAssists          int                 `json:"total_assists"`
    MaxKillsInMatch       int                 `json:"max_kills_in_match"`
    AvgDamagePerMatch     float64             `json:"avg_damage_per_match"`
    AvgObjectivesPerMatch float64             `json:"avg_objectives_per_match"`
    AvgSoulsPerMin        float64             `json:"avg_souls_per_min"`
    RecentMatches         []Match             `json:"recent_matches"`
    HeroStats             []HeroStat          `json:"hero_stats"`
    PerformanceDynamics   PerformanceDynamics `json:"performance_dynamics"`
}
```

**Frontend** (`frontend/src/entities/player/types/types.ts`):
```typescript
interface PlayerProfile {
  steam_id: string
  nickname: string
  avatar_url: string
  last_match_time: string
  last_updated_at: string
  player_rank: number
  rank_name: string
  sub_rank: number | undefined
  rank_image: string
  win_rate: number
  kd_ratio: number
  total_matches: number
  total_kills: number
  total_deaths: number
  total_assists: number
  max_kills_in_match: number
  avg_souls_per_min: number
  recent_matches: Match[]
  hero_stats: HeroStat[]
  performance_dynamics: PerformanceDynamics
}
```

**Issues**:
- Missing many fields from backend:
  - `profile_url`
  - `created_at`
  - `avg_matches_per_day`
  - `favorite_hero`
  - `avg_damage_per_match`
  - `avg_objectives_per_match`

---

## ❌ Missing Frontend Models

### 1. **Build**
**Backend** (`backend/internal/domain/builds.go`):
```go
type Build struct {
    ID          uuid.UUID `json:"id"`
    AuthorID    uuid.UUID `json:"author_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    GameVersion string    `json:"game_version"`
    IsPublic    bool      `json:"is_public"`
    ViewCount   int       `json:"view_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Frontend**: ❌ **Missing entirely** - needs to be added

---

### 2. **Comment**
**Backend** (`backend/internal/domain/comments.go`):
```go
type Comment struct {
    ID          uuid.UUID `json:"id"`
    AuthorID    uuid.UUID `json:"author_id"`
    ParentID    uuid.UUID `json:"parent_id"`
    ContentType string    `json:"content_type"`
    ContentID   uuid.UUID `json:"content_id"`
    Body        string    `json:"body"`
    CreatedAt   time.Time `json:"created_at"`
}
```

**Frontend**: ❌ **Missing entirely** - needs to be added

---

### 3. **ContentTag**
**Backend** (`backend/internal/domain/content_tags.go`):
```go
type ContentTag struct {
    TagID       int       `json:"tag_id"`
    ContentType string    `json:"content_type"`
    ContentID   uuid.UUID `json:"content_id"`
}
```

**Frontend**: ❌ **Missing entirely** - needs to be added

---

### 4. **Tag**
**Backend** (`backend/internal/domain/tags.go`):
```go
type Tag struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
```

**Frontend**: ❌ **Missing entirely** - needs to be added

---

### 5. **Vote**
**Backend** (`backend/internal/domain/votes.go`):
```go
type Vote struct {
    UserID      uuid.UUID `json:"user_id"`
    ContentType string    `json:"content_type"`
    ContentID   uuid.UUID `json:"content_id"`
    VoteValue   int       `json:"vote_value"`
    CreatedAt   time.Time `json:"created_at"`
}
```

**Frontend**: ❌ **Missing entirely** - needs to be added

---

### 6. **PlayerStats**
**Backend** (`backend/internal/domain/player_stats.go`):
```go
type PlayerStats struct {
    UserID           uuid.UUID `json:"user_id"`
    KDRatio          float64   `json:"kd_ratio"`
    WinRate          float64   `json:"win_rate"`
    AvgMatchesPerDay float64   `json:"avg_matches_per_day"`
    FavoriteHero     string    `json:"favorite_hero"`
    LastUpdatedAt    time.Time `json:"last_updated_at"`
}
```

**Frontend**: ❌ **Missing entirely** - needs to be added

---

### 7. **MateStatAPI**
**Backend** (`backend/internal/domain/mate_stat_api.go`):
```go
type MateStatAPI struct {
    MateID        int `json:"mate_id"`
    Wins          int `json:"wins"`
    MatchesPlayed int `json:"matches_played"`
}
```

**Frontend**: ❌ **Missing entirely** - needs to be added

---

### 8. **SteamProfileSearch**
**Backend** (`backend/internal/domain/steam_profile_search.go`):
```go
type SteamProfileSearch struct {
    AccountID   int    `json:"account_id"`
    Avatar      string `json:"avatar"`
    CountryCode string `json:"countrycode"`
    LastUpdated int64  `json:"last_updated"`
    Personaname string `json:"personaname"`
    Profileurl  string `json:"profileurl"`
    Realname    string `json:"realname"`
}
```

**Frontend**: ❌ **Missing entirely** - needs to be added

---

## 📝 Summary of Changes Needed

### High Priority
1. ✅ Align `Crosshair` model - remove redundant interfaces, add missing fields
2. ✅ Align `Match` model - use correct field names from backend
3. ✅ Align `HeroStat` model - use `matches_played`, add `hero_id`
4. ✅ Align `User` model - fix date types
5. ✅ Add missing `CrosshairLike` model

### Medium Priority
6. ✅ Add `Build` model
7. ✅ Add `Comment` model
8. ✅ Add `Vote` model
9. ✅ Add `Tag` model
10. ✅ Add `ContentTag` model

### Low Priority
11. ✅ Add `PlayerStats` model
12. ✅ Add `MateStatAPI` model
13. ✅ Add `SteamProfileSearch` model
14. ✅ Complete `PlayerProfile` model with all fields

---

## Implementation Plan

1. **Phase 1**: Update existing misaligned models
   - Fix Crosshair models
   - Fix Match model
   - Fix HeroStat model
   - Fix User model

2. **Phase 2**: Add missing critical models
   - Add CrosshairLike
   - Add Build
   - Add Comment
   - Add Vote
   - Add Tag
   - Add ContentTag

3. **Phase 3**: Add remaining models
   - Add PlayerStats
   - Add MateStatAPI
   - Add SteamProfileSearch

4. **Phase 4**: Update components using these models
   - Search through codebase for usages
   - Update component props and state
   - Update API calls

---

*Generated: ${new Date().toISOString()}*
