# SportSync

SportSync is a platform to provide sport hobbyist to find opponent.

## Tech Stack

**Language:** Go

**Database:** MongoDB

**Communication:** RestAPI, Websocket

## Installation

1. Create Environtment

   create file .env or copy from .env.example

   ```bash
   cp .env.example .env
   ```

   adjust value of each variable

2. Run be.sportsync.id

   ```bash
   go run cmd/main.go
   ```

3. Run with Hotreload

   ```bash
   air
   ```

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

| Name                      | Default              | Description                  |
| ------------------------- | -------------------- | ---------------------------- |
| API_PORT                  | 4000                 | Port API listened            |
| CONTEXT_TIMEOUT           | 5                    | context timeout              |
| DB_USER                   | user                 | mongodb username             |
| DB_PASS                   | pass                 | mongodb pass                 |
| DB_HOST                   | 127.0.0.1            | mongodb host                 |
| DB_NAME                   | sportsync_db         | mongodb database name        |
| ACCESS_TOKEN_EXPIRY_HOUR  | 1                    | JWT expiry in hour           |
| REFRESH_TOKEN_EXPIRY_HOUR | 1                    | refresh token expiry in hour |
| ACCESS_TOKEN_SECRET       | s3cret               | JWT secret                   |
| REFRESH_TOKEN_SECRET      | s3cret               | refresh token secret         |
| ALLOW_ORIGINS             | https://sportsync.id | allow origin                 |

## Roadmap

1. Authentication
   - Login
   - Forgot Password
   - Register
2. Manage Team
   - Create Team
   - Delete Team
   - Remove or Invite Player into Team
3. Competition
   - create Competition
   - join Competition
   - delete Competition
4. Leaderboard
   - Point get from Competition or 1 vs 1
   - by location Global or region
5. Quich Match
   - Find opponent 1 vs 1
6. Chatting
   - Chat on Team (all member)
   - Chat on Competition (captain | manager | pic )
