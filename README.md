# Army API

## Описание

REST API для работы с родами войск. 
Система позволяет создавать, публиковать, просматривать и удалять военные ветки, а также работать с лайками и пользователями.

## API методы

| Метод | URL | Назначение |
|---|---|---|
| GET | /api/military_branches | Получение опубликованных военных веток |
| GET | /api/military_branches?speed=4 | Получение веток с фильтрацией по скорости |
| POST | /api/military_branches | Создание новой военной ветки |
| GET | /api/military_branches/draft | Получение черновика текущего пользователя |
| GET | /api/military_branches/feed | Получение ленты опубликованных веток |
| GET | /api/military_branches/:id | Получение военной ветки |
| GET | /api/military_branches/:id?next=true | Получение следующей опубликованной ветки |
| POST | /api/military_branches/:id/like | Добавление или удаление лайка |
| PUT | /api/military_branches/:id | Публикация военной ветки |
| DELETE | /api/military_branches/:id | Удаление военной ветки |
| POST | /api/users | Регистрация пользователя |

## Таблицы базы данных

### military_branches

Основная таблица военных веток.

Основные поля:
- id — идентификатор;
- name — название;
- description — описание;
- status — статус;
- image_url — изображение;
- video_url — видео;
- speed_plain — скорость;
- food — потребность в продовольствии;
- creator_id — автор;
- created_at — дата создания;
- formed_at — дата публикации.

### users

Таблица пользователей.

Основные поля:
- id — идентификатор пользователя;
- login — логин;
- password — пароль.

### military_branch_likes

Таблица лайков.

Основные поля:
- id — идентификатор;
- user_id — пользователь;
- military_branch_id — военная ветка.

## Текущий пользователь

В системе используется фиксированный текущий пользователь:

CurrentUserID = 1
