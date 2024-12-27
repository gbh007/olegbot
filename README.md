# OlegBot

**Внимание:** вторая версия не является совместимой с 1-й.

Бот для Telegram который может выдавать рандомные цитаты из добавленных.

## Переменные окружения

| Переменная            | Описание                           | Пример                                                            |
| --------------------- | ---------------------------------- | ----------------------------------------------------------------- |
| `REPO`                | Строка подключения для PostgreSQL  | `postgres://oleg:olegpass@localhost:5432/olegbot?sslmode=disable` |
| `ADDR`                | Веб адрес для метрик и CMS         | `:8080`                                                           |
| `CMS_LOGIN`           | Логин для CMS                      | `user`                                                            |
| `CMS_PASSWORD`        | Пароль для CMS                     | `password`                                                        |
| `CMS_STATIC_DIR_PATH` | Путь до статических файлов для CMS | `/app/static`                                                     |
| `DEBUG`               | Режим отладки                      | `true`                                                            |

## Метрики

Бот отдает стандартные метрики для Prometheus, дашборд для графаны `grafana-dashboard.json`

## CMS

У бота есть система управления контентом, запускается если указан веб адрес.  
Если логин и пароль не пустые, то применяется базовая аутентификация по ним, кроме ендпоинта метрик.
