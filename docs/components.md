# API управления компонентами

- [Главная](./README.md)

## Methods

- [Add component](#add-component)
- [Delete component](#delete-component)
- [Get bot structure](#get-bot-structure)
- [Set next step for component](#set-next-step-for-component)
- [Delete next step for component](#delete-next-step-for-component)
- [Set next step for command](#set-next-step-for-command)
- [Delete next step for command](#delete-next-step-for-command)
- [Add command](#add-command)
- [Delete command](#delete-command)


- - -


## Add component

[Наверх][toup]

Добавление нового компонента в структуру бота

```plaintext
POST /api/bots/{botId}/components
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

Параметры тела

```json
{
    "data": {
        "type": "string",
        "content": [
            "content"
        ]
    },
    "commands": [
        "command"
    ],
    "position": {
        "x": "integer",
        "y": "integer"
    }
}
```

Поле           | Тип                       | Описание
---------------|---------------------------|-----------------------------------------------------------
`data`         | object                    | Данные компонента
`data.type`    | string                    | Тип компонента
`data.content` | [content][type_content][] | Список c данными, спецефичными для каждого типа компонента
`commands`     | [command][type_command][] | Список команд
`position`     | object                    | Координаты компонента на поле редактора
`position.x`   | integer                   | Координата X
`position.y`   | integer                   | Координата Y

#### Ответ

Включает только одно из полей: `data`, `error`  

```json
{
    "ok": "bool",
    "data": {
        "id": "integer",
    },
    "error": {
        "code": "integer",
        "message": "string"
    }
}
```

_data_

Поле | Тип     | Описание
-----|---------|--------------
`id` | integer | id компонента

<details>
    <summary>Пример</summary>
   
`Запрос`

```plaintext
POST /api/bots/64/components
```

Тело запроса

```json
{
    "data": {
        "type": "text",
        "content": [
            {
                "text": "Hello Telegram"
            }
        ]
    },
    "commands": [
        {
            "type": "text",
            "data": "First button"
        },
        {
            "type": "text",
            "data": "Second button"
        }
    ],
    "position": {
        "x": 141,
        "y": 112
    }
}
```

`Ответ` 

```json
{
    "ok": true,
    "data": {
        "id": 7
    }
}
```
</details>


- - -


## Delete component

[Наверх][toup]

Удаление компонента из структуры бота

```plaintext
DELETE /api/bots/{botId}/components/{compId}
```

Параметры пути

Поле     | Описание
---------|--------------
`botId`  | id бота
`compId` | id компонента

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


## Get bot structure

[Наверх][toup]

Получение структуры бота  
(Список всех компонентов входящих в структуру)

```plaintext
GET /api/bots/{botId}/components
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

#### Ответ

Включает только одно из полей: `data`, `error` 

```json
{
    "ok": "bool",
    "data": [
        "component"
    ],
    "error": {
        "code": "integer",
        "message": "string"
    }
}
```

Поле        | Тип                         | Описание
------------|-----------------------------|-------------------
`data`      | component[]                 | Список компонентов
`component` | [component][type_component] | Компонент

<details>
    <summary>Пример</summary>
   
`Запрос`

```plaintext
GET /api/bots/67/components
```

`Ответ` 

```json
{
    "ok": true,
    "data": [
        {
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
            "isStart": true,
            "position": {
                "x": 50,
                "y": 50
            }
        },
        {
            "id": 2,
            "data": {
                "type": "text",
                "content": [
                    {
                        "text": "Hello Telegram"
                    }
                ]
            },
            "keyboard": {
                "buttons": []
            },
            "commands": [
                {
                    "id": 1,
                    "type": "text",
                    "data": "First button",
                    "componentId": 2,
                    "nextStepId": null
                },
                {
                    "id": 2,
                    "type": "text",
                    "data": "Second button",
                    "componentId": 2,
                    "nextStepId": null
                }
            ],
            "nextStepId": null,
            "isStart": false,
            "position": {
                "x": 141,
                "y": 112
            }
        }
    ]
}
```
</details>


- - -


## Set next step for component

[Наверх][toup]

Установка номера следующего шага для компонента

```plaintext
POST /api/bots/{botId}/components/{compId}/next
```

Параметры пути

Поле     | Описание
---------|------------------------
`botId`  | id бота
`compId` | id исходного компонента

Параметры тела

```json
{
    "nextStepId": "integer"
}
```

Поле         | Тип     | Описание
-------------|---------|--------------------------------
`nextStepId` | integer | id следующего шага (компонента)

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


## Delete next step for component

[Наверх][toup]

Удаление номера следующего шага для компонента

```plaintext
DELETE /api/bots/{botId}/components/{compId}/next
```

Параметры пути

Поле     | Описание
---------|------------------------
`botId`  | id бота
`compId` | id исходного компонента

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


## Set next step for command

[Наверх][toup]

Установка номера следующего шага для команды в компоненте

```plaintext
POST /api/bots/{botId}/commands/{commandId}/next
```

Параметры пути

Поле        | Описание
------------|-----------
`botId`     | id бота
`commandId` | id команды

Параметры тела

```json
{
    "nextStepId": "integer"
}
```

Поле         | Тип     | Описание
-------------|---------|--------------------------------
`nextStepId` | integer | id следующего шага (компонента)

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


## Delete next step for command

[Наверх][toup]

Удаление номера следующего шага для команды в компоненте

```plaintext
DELETE /api/bots/{botId}/commands/{commandId}/next
```

Параметры пути

Поле        | Описание
------------|-----------
`botId`     | id бота
`commandId` | id команды

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


## Add command

[Наверх][toup]

Добавление команды в компоненте

```plaintext
POST /api/bots/{botId}/components/{compId}/commands
```

Параметры пути

Поле     | Описание
---------|--------------
`botId`  | id бота
`compId` | id компонента

Параметры тела

```json
{
    "type": "string",
    "data": "string"
}
```

Поле   | Тип    | Описание
-------|--------|---------------------------------------------
`type` | string | Тип команды
`data` | string | Данные, спецефичные для каждого типа команды

#### Ответ

Включает только одно из полей: `data`, `error`  

```json
{
    "ok": "bool",
    "data": {
        "id": "integer",
    },
    "error": {
        "code": "integer",
        "message": "string"
    }
}
```

_data_

Поле | Тип     | Описание
-----|---------|-----------
`id` | integer | id команды

<details>
    <summary>Пример</summary>  

`Запрос`

```plaintext
POST /api/bots/67/components/3/commands
```

Тело запроса

```json
{
    "type": "text",
    "data": "abc"
}
```

`Ответ` 

```json
{
    "ok": true,
    "data": {
        "id": 4
    }
}
```
</details>


- - -


## Delete command

[Наверх][toup]

Удаление команды в компоненте

```plaintext
DELETE /api/bots/{botId}/commands/{commandId}
```

Параметры пути

Поле        | Описание
------------|-----------
`botId`     | id бота
`commandId` | id команды

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
[type_component]: ./objects.md#component
[type_content]: ./objects.md#content
[type_command]: ./objects.md#command
[toup]: #api-управления-компонентами