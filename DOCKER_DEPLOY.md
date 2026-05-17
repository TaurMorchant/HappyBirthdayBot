# Docker Deploy

## Репозитории

| Репозиторий | Назначение |
|---|---|
| `HappyBirthdayBot` | Код бота, Dockerfile. При пуше в `main` — сборка и публикация образа в ghcr.io |
| `HappyBirthdayBot-configs` | Конфиги, `docker-compose.yml`, deploy workflow. Запускается вручную |

## Как это работает

**Обновление кода** — пуш в `main` репозитория `HappyBirthdayBot`:
1. Собирает Docker-образ
2. Публикует в `ghcr.io/taurmorchant/happybirthdaybot:latest`

**Деплой на сервер** — ручной запуск workflow в `HappyBirthdayBot-configs`:
1. Подготавливает VPS (Docker, директории, `.env`)
2. Копирует конфиги и `docker-compose.yml` на сервер
3. Запускает бота с выбранным окружением (`prod` или `test`)

---

## Первоначальная настройка (один раз)

### 1. Добавить секреты в `HappyBirthdayBot-configs`

**https://github.com/TaurMorchant/HappyBirthdayBot-configs/settings/secrets/actions**

| Секрет | Описание |
|---|---|
| `TELEGRAM_BOT_TOKEN` | Токен бота |
| `DEPLOY_HOST` | IP сервера |
| `DEPLOY_USER` | Пользователь SSH (`root`) |
| `DEPLOY_SSH_KEY` | Приватный SSH-ключ для доступа к серверу |

Как сгенерировать SSH-ключ для деплоя:
```powershell
ssh-keygen -t ed25519 -C "github-deploy" -f deploy_key
```
Публичный ключ добавить на сервер:
```powershell
type deploy_key.pub | ssh root@VPS_IP "cat >> ~/.ssh/authorized_keys"
```
Приватный ключ (`deploy_key`) — в секрет `DEPLOY_SSH_KEY`.

### 2. Запустить деплой

**https://github.com/TaurMorchant/HappyBirthdayBot-configs/actions** → **Deploy** → **Run workflow** → выбрать окружение (`prod` / `test`).

---

## Обновление бота

| Что изменилось | Действие |
|---|---|
| Код бота | Пуш в `main` → образ пересобирается автоматически. Затем запустить деплой |
| Конфиги или `docker-compose.yml` | Пуш в `HappyBirthdayBot-configs`, затем запустить деплой |

---

## Полезные команды на сервере

```bash
# Логи в реальном времени
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
│   ├── configs-prod/
│   │   ├── application.properties
│   │   ├── allowedUsers.properties
│   │   ├── allowedChats.properties
│   │   └── birthdayChats.csv
│   └── configs-test/
├── data-prod/
│   └── data.db                # SQLite база prod
├── data-test/
│   └── data.db                # SQLite база test
├── docker-compose.yml         # копируется из HappyBirthdayBot-configs
├── logs/                      # персистентные лог-файлы
└── .env                       # TELEGRAM_BOT_TOKEN + ENVIRONMENT (не в git!)
```
