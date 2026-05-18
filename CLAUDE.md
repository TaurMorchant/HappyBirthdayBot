# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

All Go code lives in `service/`. Commands must be run from that directory.

```bash
# Build
cd service && go build -o bot .

# Run locally (requires TELEGRAM_BOT_TOKEN env var and a config directory path as CLI arg)
# Configs live in HappyBirthdayBot-configs repo under profiles/<env>/configs/
$env:TELEGRAM_BOT_TOKEN="<token>"
.\bot <path-to-configs>     # e.g. ../../HappyBirthdayBot-configs/profiles/test/configs

# Tests
cd service && go test ./...

# Single test
cd service && go test ./handlers/... -run Test_getDaysWord
```

## Docker

CI builds and publishes the image automatically on every push to `main`:
```
ghcr.io/taurmorchant/happybirthdaybot:latest
```

The image expects a config directory mounted at `/config` (default `CMD`). Deploy is triggered manually via the `HappyBirthdayBot-configs` repository — see that repo's `CLAUDE.md`.

## Architecture Overview

The bot is a single-process Go application with two concurrent flows:

1. **Telegram update loop** (`main.go`) — polls for updates, routes to command handlers or reply/callback handlers
2. **Cron scheduler** (`mybot/reminder_task.go`) — checks birthdays on schedule and sends proactive reminders

### Key packages

- **`mybot/`** — bot wrapper (`Bot` struct embedding `tgbotapi.BotAPI`), command constants, cron reminder task
- **`handlers/`** — one file per command, all implementing `IHandler` interface (3 methods: `Handle`, `HandleReply`, `HandleCallback`)
- **`db/`** — SQLite data layer (`modernc.org/sqlite`, pure Go); `Init`, `ReadUsers`, `InsertUser`, `DeleteUser`, `UpdateWishlist`, `UpdateFlags`
- **`usr/`** — `User` and `Users` domain types; `daysBeforeBirthday` is computed at read-time
- **`date/`** — `Birthday` type (day+month only, no year); Russian month name parsing/formatting
- **`config/`** — loads all config files from the directory passed as CLI arg; also reads `birthdayChats.csv`
- **`resources/`** — cat images embedded via `//go:embed *`; `ImageKey` pointing to a directory triggers random image selection
- **`cache/`** — generic wrapper over `go-cache` for pending interaction state

### Multi-step interaction pattern

Multi-step commands (join, wishlist, exit) work via two in-memory caches in `handlers/handlers.go`:

- `WaitingForReplyHandlers` — maps `userID → IHandler`, set when bot sends a `ForceReply` message; cleared after `HandleReply` succeeds
- `WaitingForCallbackHandlers` — maps `messageID → CallbackElement`, set when bot sends inline keyboard; cleared after `HandleCallback`

Both caches have a 5-minute TTL.

### Reminder logic (reminder_task.go)

On each cron tick, reads all users from SQLite and applies these transitions (flags written back per-user immediately after each action):
- `DaysBeforeBirthday == 0` → send birthday congratulation, set `BirthdayGreetings=true`
- `≤ 14 days` → send 14-day reminder (with birthday chat link if configured)
- `≤ 30 days` → send 30-day reminder, pin message, send wishlist to birthday chat
- `> 30 days && BirthdayGreetings == true` → reset all three flags for the next year

### Data storage (SQLite)

File: `data/data.db` relative to working directory (mounted as a Docker volume in production).

Table `users`: `id | name | birthday | wishlist | reminder30days | reminder15days | birthday_greetings`

Birthday format: `"25 марта"` (day + Russian month name). DB opened with WAL mode and `busy_timeout=5000ms`.

## Configuration Files

Config directory is passed as the first CLI argument (e.g., `../../HappyBirthdayBot-configs/profiles/test/configs`).

| File | Purpose |
|---|---|
| `application.properties` | `mainChatId`, `adminChatId`, `reminderTriggerCron` (cron with seconds) |
| `allowedUsers.properties` | `<telegramUserId>: <name>` — who can interact with the bot |
| `allowedChats.properties` | `<chatId>: <name>` — additional allowed chats beyond birthday chats |
| `birthdayChats.csv` | `UserId, Name, ChatLink, ChatId` — per-person birthday discussion chats |

Config files live in the `HappyBirthdayBot-configs` repository under `profiles/<env>/configs/`.

## Environment Variable

`TELEGRAM_BOT_TOKEN` — must be set before running the bot. The bot panics on startup if missing.