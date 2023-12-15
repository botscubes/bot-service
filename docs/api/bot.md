# API управления ботом

- [Главная](../README.md)

## Methods

- [New bot](#new-bot)
- [Get bots](#get-bots)
- [Set token](#set-token)
- [Delete token](#delete-token)
- [Start](#start)
- [Stop](#stop)
- [Wipe bot](#wipe-bot)


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

После создания бот включает в свою структуру стартовый компонент

```json
{
    "botId": "integer",
    "component": "component"
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
```
</details>


- - -


## Get bots

[Наверх][toup]

Получение списка ботов пользователя

```plaintext
GET /api/bots
```

#### Ответ

После создания бот включает в свою структуру стартовый компонент

```json
[
    "bot"
]
```

Поле   | Тип             | Описание
-------|-----------------|-------------
`data` | bot[]           | Список ботов
`bot`  | [bot][type_bot] | Бот

<details>
    <summary>Пример</summary>
   
`Запрос`

```plaintext
GET /api/bots
```

`Ответ` 

```json
[
    {
        "id": 79,
        "title": "qwerty",
        "status": 1
    },
    {
        "id": 80,
        "title": "qwerty",
        "status": 0
    },
    {
        "id": 124,
        "title": "--",
        "status": 1
    }
]

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

В случае успеха http статус 204 без тела ответа.



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

В случае успеха http статус 204 без тела ответа.



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

В случае успеха http статус 204 без тела ответа.



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

В случае успеха http статус 204 без тела ответа.


- - -


## Wipe bot

[Наверх][toup]

Очистка бота:  
- Удаление всех компонентов в структре, кроме начального  
- Удаление всех команд  
- Сброс токена  

```plaintext
PATCH /api/bots/{botId}/wipe
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

#### Ответ

В случае успеха http статус 204 без тела ответа.


[//]: # (LINKS)
[type_component]: ../objects.md#component
[type_bot]: ../objects.md#bot
[toup]: #api-управления-ботом 
