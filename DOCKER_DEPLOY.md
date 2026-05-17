# Docker Deploy

## Как это работает

При каждом пуше в `main` GitHub Actions собирает Docker-образ и публикует его в GitHub Container Registry:
```
ghcr.io/taurmorchant/happybirthdaybot:latest
```

На VPS образ запускается через Docker Compose с автоматическим рестартом.

---

## Первоначальная настройка VPS (один раз)

### Шаг 1 — Скопировать конфиги вручную (единственный ручной шаг)

Конфиги содержат чувствительные данные и не хранятся в репозитории — копируются один раз с локальной машины:

```powershell
scp C:\go_modules\happy_birthday_bot\configs-prod\* root@VPS_IP:/opt/happy-birthday-bot/configs/
```

### Шаг 2 — Запустить pipeline

Сделать любой пуш в `main`. Pipeline автоматически:
- Установит Docker и Docker Compose plugin (если не стоят)
- Создаст директории и `.env` с токеном
- Скопирует `docker-compose.yml`
- Запустит бота

---

## Обновление бота

### Только код (docker-compose.yml не менялся)

```bash
cd /opt/happy-birthday-bot
docker compose pull && docker compose up -d
```

### Если изменился docker-compose.yml

Скопировать обновлённый файл с локальной машины (Windows):
```powershell
scp C:\go_modules\happy_birthday_bot\docker-compose.yml root@VPS_IP:/opt/happy-birthday-bot/
```

Затем на сервере пересоздать контейнер:
```bash
cd /opt/happy-birthday-bot
docker compose up -d --force-recreate
```

---

## Полезные команды

```bash
# Посмотреть логи в реальном времени
docker compose logs -f bot

# Статус контейнера
docker compose ps

# Перезапустить
docker compose restart bot

# Остановить
docker compose down
```

---

## Структура файлов на VPS

```
/opt/happy-birthday-bot/
├── configs/
│   ├── application.properties
│   ├── allowedUsers.properties
│   ├── allowedChats.properties
│   ├── birthdayChats.csv
│   └── happybirthdaybot-454814-2dec5157295e.json
├── logs/                  # персистентные лог-файлы (создаётся автоматически)
├── .env                   # TELEGRAM_BOT_TOKEN (не в git!)
└── docker-compose.yml
```
