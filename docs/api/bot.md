# API управления ботом

- [Главная](../README.md)

## Methods

- [New bot](#new-bot)
- [Get bots](#get-bots)
- [Delete bot](#delete-bot)
- [Set token](#set-token)
- [Delete token](#delete-token)
- [Start](#start)
- [Stop](#stop)
- [Get status](#get-status)

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
`component` | component                   | Компонент

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
        "type": "start",
        "path": "",
        "outputs": {
            "nextComponentId": null
        },
        "position": {
            "x": 50,
            "y": 50
        }
    }
}
```
</details>


- - -
## Delete bot

[Наверх][toup]

Удаление бота

```plaintext
DELETE /api/bots/{botId}
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

#### Ответ

В случае успеха http статус 204 без тела ответа.


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
`bot`  | bot             | Бот


Структура бота:

```
{
    "id": "integer",
    "title": "string",
    "status": "integer"
}
```


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


## Get status

[Наверх][toup]

Получение статуса бота

```plaintext
GET /api/bots/{botId}/status
```

В случае успеха http статус 200 с телом ответа, содержащим номер статуса бота:
- 0 - остановлен;
- 1 - запущен.


[//]: # (LINKS)
[type_component]: ../objects.md#component
[type_bot]: ../objects.md#bot
[toup]: #api-управления-ботом 
