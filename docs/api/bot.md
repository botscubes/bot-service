# API управления ботом

- [Главная](../README.md)

## Methods

- [New bot](#new-bot)
- [Set token](#set-token)
- [Delete token](#delete-token)
- [Start](#start)
- [Stop](#stop)


- - -


## New bot

[Наверх][toup]

Создание нового бота

```plaintext
POST /api/bots
```

Параметры тела запроса

```json
{
    "title": "string"
}
```

Поле    | Тип    | Описание
--------|--------|--------------
`title` | string | Название бота

#### Ответ

Включает только одно из полей: `data`, `error`  
После создания бот включает в свою структуру стартовый компонент

```json
{
    "ok": "bool",
    "data": {
        "botId": "integer",
        "component": "component"
    },
    "error": {
        "code": "integer",
        "message": "string"
    }
}
```

_data_

Поле        | Тип                         | Описание
------------|-----------------------------|----------
`botId`     | integer                     | id бота
`component` | [component][type_component] | Компонент

<details>
    <summary>Пример</summary>
   
`Запрос`

```plaintext
POST /api/bots
```

Тело запроса

```json
{
    "title": "qwerty"
}
```

`Ответ` 

```json
{
    "ok": true,
    "data": {
        "botId": 66,
        "component": {
            "id": 1,
            "data": {
                "type": "start",
                "content": []
            },
            "keyboard": {
                "buttons": []
            },
            "commands": [],
            "nextStepId": null,
            "isMain": true,
            "position": {
                "x": 50,
                "y": 50
            }
        }
    }
}
```
</details>


- - -


## Set token

[Наверх][toup]

Установка токена бота

```plaintext
POST /api/bots/{botId}/token
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

Параметры тела запроса

```json
{
    "token": "string"
}
```

Параметры | Тип    | Описание
----------|--------|---------
`token`   | string | Токен

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

[Наверх][toup]

Удаление токена бота

```plaintext
DELETE /api/bots/{botId}/token
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

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

[Наверх][toup]

Запуск бота

```plaintext
PATCH /api/bots/{botId}/start
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

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

[Наверх][toup]

Остановка бота

```plaintext
PATCH /api/bots/{botId}/stop
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

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

[//]: # (LINKS)
[type_component]: ../objects.md#component
[toup]: #api-управления-ботом 