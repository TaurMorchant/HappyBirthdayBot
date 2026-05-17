# Redesign Plan

## Удаляем

### Команды `/join` и `/exit`
Пользователи больше не регистрируются сами. Список участников задаётся статически через конфиг и не меняется в рантайме.

### Google Sheets как хранилище данных
Убираем зависимость от Google Sheets API, google-credentials JSON, oauth2/jwt конфигурацию. Пакет `sheets/` удаляется целиком.

---

## Меняем

### Список участников — из конфига, не из Sheets

Вместо `/join`/`/exit` список пользователей задаётся в конфиг-файле. Формат — аналогично текущему `allowedUsers.properties` (Telegram user ID → имя), но расширенный: добавляется дата рождения.

Пример:
```
287959887 = Влад, 15 мая
1811998414 = Миша, 3 февраля
```

Или отдельный CSV — решить при реализации.

### Хранилище данных — SQLite ✅

Данные о пользователях (имя, дата рождения, вишлист, флаги напоминаний) хранятся в SQLite-файле на диске.

- Флаги `reminder30days`, `reminder15days`, `birthdayGreetings` по-прежнему сбрасываются в конце цикла (>30 дней до ДР)
- База доступна для ручного просмотра/редактирования через DB Browser for SQLite или sqlite3 CLI
- Бэкап — копирование файла
- Драйвер: `modernc.org/sqlite` (чистый Go, без CGO), WAL-режим, `busy_timeout=5000ms`

### Конфигурация упрощается

Убираем: google credentials JSON, `spreadsheetList`, `spreadsheetID`.  
Остаётся: `mainChatId`, `adminChatId`, `reminderTriggerCron`, путь к SQLite-файлу.  
`birthdayChats.csv` — остаётся без изменений.  
`allowedUsers.properties` — расширяется датой рождения (заменяет ручной ввод через `/join`).

---

## Оставляем без изменений

### Команды бота
- `/start` — приветственное сообщение
- `/list` — список всех участников с датами
- `/reminders` — ближайшие 3 именинника
- `/wishlist` — просмотр и редактирование своего вишлиста (через force-reply и inline keyboard)

### Логика напоминаний
Cron-задача без изменений:
- **≤ 30 дней** — сообщение в основной чат, пин, отправка вишлиста в birthday chat (если есть)
- **≤ 14 дней** — напоминание в основной чат со ссылкой на birthday chat
- **День X** — поздравление в основной чат
- **> 30 дней** — сброс флагов для следующего года

### Birthday chats
Отдельные чаты обсуждения подарка для каждого именинника — остаются. Конфигурируются через `birthdayChats.csv`.

### Мультишаговые диалоги
Механика `WaitingForReplyHandlers` / `WaitingForCallbackHandlers` для `/wishlist` — остаётся.

---

## Новое — деплой на VPS через Docker

### Dockerfile
Многоэтапная сборка: сборка бинарника в официальном Go-образе, финальный образ на `alpine` или `scratch`.

### Docker Compose
На VPS запускается через `docker compose`. Файл `.env` рядом с `compose.yml` содержит `TELEGRAM_BOT_TOKEN`.  
Конфиг-директория и SQLite-файл монтируются как volume с хоста.

### CI/CD

**Сборка** (`HappyBirthdayBot`) — автоматически при пуше в `main`: сборка образа и пуш в ghcr.io.

**Деплой** (`HappyBirthdayBot-configs`) — ручной запуск workflow с выбором окружения (`prod`/`test`):
- Подготовка VPS (Docker, директории, `.env`)
- Копирование конфигов и `docker-compose.yml` на сервер
- `docker compose pull && docker compose up -d`

Секреты деплоя хранятся в `HappyBirthdayBot-configs`: `TELEGRAM_BOT_TOKEN`, `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_SSH_KEY`.  
Для деплоя используется отдельная SSH-пара ключей (не личный ключ разработчика).

### Автоперезапуск
`restart: unless-stopped` в compose-файле — бот поднимается автоматически после перезагрузки сервера.

---

## Открытые вопросы

- **Формат расширенного `allowedUsers`** — актуально только если уберём `/join`; пока команда оставлена

---

## В плане

### Реструктуризация папок по профилям

Текущая структура (`configs-prod/`, `configs-test/`, `data-prod/`, `data-test/`) заменяется на профили:

```
profiles/
├── prod/
│   ├── configs/
│   │   ├── application.properties
│   │   ├── allowedUsers.properties
│   │   ├── allowedChats.properties
│   │   └── birthdayChats.csv
│   └── data/
│       └── data.db
└── test/
    ├── configs/
    └── data/
```

Затрагивает: `HappyBirthdayBot-configs` (структура репо, deploy workflow, docker-compose volumes), путь к конфигам в CLI-аргументе бота.

### Workflow для первоначальной загрузки данных из Google Forms

Ручной workflow в `HappyBirthdayBot-configs`, который применяет SQL-скрипт с данными на сервере:

- Данные (имена, даты рождения, вишлисты) выгружаются из Google Forms в SQL INSERT-скрипт вручную
- Workflow копирует скрипт на VPS и выполняет `sqlite3 data.db < seed.sql`
- Используется один раз для первоначального наполнения базы; после этого данные ведутся через бота