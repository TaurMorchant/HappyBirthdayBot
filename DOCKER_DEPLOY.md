# Docker Deploy

## Как это работает

При каждом пуше в `main` GitHub Actions собирает Docker-образ и публикует его в GitHub Container Registry:
```
ghcr.io/taurmorchant/happybirthdaybot:latest
```

На VPS образ запускается через Docker Compose с автоматическим рестартом.

---

## Первоначальная настройка VPS (один раз)

```bash
# 1. Установить Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
# перелогиниться

# 2. Создать рабочую директорию
mkdir -p /opt/happy-birthday-bot/configs
cd /opt/happy-birthday-bot

# 3. Скопировать конфиги с локальной машины
scp -r configs-prod/* user@VPS_IP:/opt/happy-birthday-bot/configs/

# 4. Создать .env с токеном бота
echo "TELEGRAM_BOT_TOKEN=ВАШ_ТОКЕН" > .env
chmod 600 .env

# 5. Скопировать docker-compose.yml
scp docker-compose.yml user@VPS_IP:/opt/happy-birthday-bot/

# 6. Запустить
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
├── .env                   # TELEGRAM_BOT_TOKEN (не в git!)
└── docker-compose.yml
```
