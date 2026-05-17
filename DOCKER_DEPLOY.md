# Docker Deploy

## Как это работает

При каждом пуше в `main` GitHub Actions собирает Docker-образ и публикует его в GitHub Container Registry:
```
ghcr.io/taurmorchant/happybirthdaybot:latest
```

На VPS образ запускается через Docker Compose с автоматическим рестартом.

---

## Первоначальная настройка VPS (один раз)

### Шаг 1 — На сервере: установить Docker и создать директорию

```bash
# Установить Docker (если не стоит)
curl -fsSL https://get.docker.com | sh

# Установить Docker Compose plugin (если Ubuntu — из apt может не поставиться, ставим вручную)
mkdir -p /usr/local/lib/docker/cli-plugins
curl -SL "https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64" \
  -o /usr/local/lib/docker/cli-plugins/docker-compose
chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

# Создать рабочую директорию
mkdir -p /opt/happy-birthday-bot/configs
```

### Шаг 2 — На локальной машине (Windows): скопировать файлы на сервер

```powershell
# Конфиги
scp C:\go_modules\happy_birthday_bot\configs-prod\* root@VPS_IP:/opt/happy-birthday-bot/configs/

# docker-compose.yml
scp C:\go_modules\happy_birthday_bot\docker-compose.yml root@VPS_IP:/opt/happy-birthday-bot/
```

### Шаг 3 — На сервере: создать .env и запустить

```bash
cd /opt/happy-birthday-bot

# Создать файл с токеном бота
echo "TELEGRAM_BOT_TOKEN=ВАШ_ТОКЕН" > .env
chmod 600 .env

# Запустить
docker compose pull
docker compose up -d
```

---

## Обновление бота

```bash
cd /opt/happy-birthday-bot
docker compose pull && docker compose up -d
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
