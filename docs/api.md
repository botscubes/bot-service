# Bot API

## Methods

- [New bot](#new-bot)
- [Set token](#set-token)
- [Delete token](#delete-token)
- [Start](#start)
- [Stop](#stop)


## New bot

[Наверх](#methods)

Создание нового бота

```plaintext
POST /api/bot/new
```

```json
{
    "title": "string"
}
```

Параметры | Тип | Описание 
--------- | ---- | -----------
`title` | string | Название бота

#### Ответ

Включает только одно из полей: `data`, `error`

```json
{
    "ok": "bool",
    "data": {
        "id": "integer"
    },
    "error": {
        "code": "integer",
        "message": "string"
    }
}

```
_response.data_

Параметры | Тип | Описание 
--------- | ---- | -----------
`id` | integer  | id бота

- - -

## Set token

[Наверх](#methods)

Установка токена бота

```plaintext
POST /api/bot/setToken
```

```json
{
    "bot_id": "integer",
    "token": "string"
}
```

Параметры | Тип | Описание 
--------- | ---- | -----------
`bot_id` | integer | id бота
`token` | string | Токен

#### Ответ

В случае успеха включает только поле `ok`

```json
{
    "ok": "bool",
    "error": {
        "code": "integer",
        "message": "string"
    }
}

```

- - -

## Delete token

[Наверх](#methods)

Удаление токена бота

```plaintext
POST /api/bot/deleteToken
```

```json
{
    "bot_id": "integer"
}
```

Параметры | Тип  | Описание 
--------- | ---- | -----------
`bot_id` | integer | id бота

#### Ответ

В случае успеха включает только поле `ok`

```json
{
    "ok": "bool",
    "error": {
        "code": "integer",
        "message": "string"
    }
}

```

- - -

## Start

[Наверх](#methods)

Запуск бота

```plaintext
POST /api/bot/start
```

```json
{
    "bot_id": "integer"
}
```

Параметры | Тип  | Описание 
--------- | ---- | -----------
`bot_id` | integer | id бота

#### Ответ

В случае успеха включает только поле `ok`

```json
{
    "ok": "bool",
    "error": {
        "code": "integer",
        "message": "string"
    }
}
```

- - -

## Stop

[Наверх](#methods)

Остановка бота

```plaintext
POST /api/bot/stop
```

```json
{
    "bot_id": "integer"
}
```

Параметры | Тип  | Описание 
--------- | ---- | -----------
`bot_id` | integer | id бота

#### Ответ

В случае успеха включает только поле `ok`

```json
{
    "ok": "bool",
    "error": {
        "code": "integer",
        "message": "string"
    }
}
```
