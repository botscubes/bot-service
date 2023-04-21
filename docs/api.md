# Bot API

## Методы

## New bot

Создать нового бота

```plaintext
POST /new
```

```json
{
    "user_id": "int64",
    "title": "string"
}
```

Параметры | Тип | Обязательный | Описание 
--------- | ---- | -------- | -----------
`user_id` | int64 | yes | Id пользователя
`title` | string | yes | Название бота

#### Ответ

```json
{
    "ok": "bool",

    //  includes only one of the fields `data`, `error`
    "data": {
        "id": "int64"
    },
    "error": {
        "code": "int",
        "message": "string"
    }
}

```

#### Список возможных ошибок

Код ошибки | Текст ошибки | Описание ошибки
----- | ----- | -----
1400 | Invalid request | Техническая ошибка
1401 | Required parameters are missing | Отсутсвует обязательный параметр
103 | Title is too long | Название бота слишком длинное
