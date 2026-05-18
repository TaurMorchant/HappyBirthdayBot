# Задачи по развитию бота

## ✅ 1. Миграция хранилища данных на SQLite

**Что:** Заменить Google Sheets на локальную SQLite-базу. Удалить зависимости от Google Sheets API, google-credentials JSON, oauth2/jwt, пакета `sheets/`.

**Зачем:** Google Sheets — внешняя зависимость с аутентификацией, квотами и сложным SDK. SQLite — файл на диске, который можно открыть любым инструментом, скопировать как бэкап, и который не требует интернета.

**Детали реализации:**
- Таблица `users`: `id | name | birthday | wishlist | reminder30days | reminder15days | birthday_greetings`
- Формат даты рождения: `"25 марта"` (день + русское название месяца)
- Драйвер: `modernc.org/sqlite` (чистый Go, без CGO)
- WAL-режим, `busy_timeout=5000ms`
- Файл базы: `data/data.db` относительно рабочей директории (монтируется как Docker volume в проде)

**Статус: выполнено**

---

## ✅ 2. Деплой на VPS через Docker

**Что:** Упаковать бота в Docker-образ и настроить CI/CD для автоматической сборки и ручного деплоя на VPS.

**Зачем:** Воспроизводимое окружение, изоляция от системного Go на сервере, автоперезапуск при сбоях и перезагрузке VPS.

**Детали реализации:**
- Dockerfile: многоэтапная сборка (Go-образ для компиляции, `alpine`/`scratch` для финального образа)
- Docker Compose: `.env` рядом с `compose.yml` содержит `TELEGRAM_BOT_TOKEN`; конфиг-директория и SQLite-файл монтируются с хоста
- `restart: unless-stopped` — автоперезапуск после перезагрузки сервера
- CI workflow `HappyBirthdayBot`: сборка образа и пуш в ghcr.io при каждом пуше в `main`
- Deploy workflow `HappyBirthdayBot-configs`: ручной запуск с выбором окружения (`prod`/`test`); подготовка VPS, копирование конфигов, `docker compose pull && docker compose up -d`
- Секреты деплоя (`TELEGRAM_BOT_TOKEN`, `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_SSH_KEY`) хранятся в репо `HappyBirthdayBot-configs`; для деплоя используется отдельная SSH-пара ключей

**Статус: выполнено**

---

## ✅ 3. Реструктуризация директорий по профилям

**Что:** Заменить плоскую структуру `configs-prod/`, `configs-test/`, `data-prod/`, `data-test/` на структуру по профилям.

**Зачем:** Сейчас конфиги и данные для разных окружений лежат вперемешку на одном уровне. Группировка по профилям делает структуру очевидной и упрощает добавление новых окружений.

**Целевая структура:**
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

**Что затрагивает:**
- CLI-аргумент бота: теперь передаётся путь к `profiles/<env>/configs/`
- `HappyBirthdayBot-configs`: структура репо, пути в deploy workflow, volume-маппинги в docker-compose
- Локальный запуск: команда в README/CLAUDE.md

**Статус: выполнено**

---

## ✅ 4. Первоначальная загрузка данных из Google Forms

**Что:** Ручной GitHub Actions workflow, который применяет SQL-скрипт с начальными данными на VPS.

**Зачем:** При первом деплое база пустая. Данные участников (имена, даты рождения, вишлисты) уже собраны в Google Forms — нужен способ загрузить их в SQLite без ручного SSH на сервер.

**Детали реализации:**
- `profiles/prod/scripts/seed.sql` — 19 реальных записей из Google Sheets
- `profiles/test/scripts/seed.sql` — 3 фиктивные записи для проверки напоминалок
- `seed.sql` деплоится на VPS вместе с конфигами через workflow Deploy
- Чекбокс `seed` в workflow Deploy: опциональный шаг после деплоя, выполняет `sqlite3 data.db < seed.sql`
- `INSERT OR IGNORE` — безопасно запускать повторно

**Статус: выполнено**