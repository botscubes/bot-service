# API управления компонентами

- [Главная](../README.md)

## Methods

- **COMPONENTS:**
- [Add component](#add-component)
- [Delete component](#delete-component)
- [Delete set of components](#delete-set-of-components)
- [Update component](#update-component)
- [Get bot structure](#get-bot-structure)
- **NEXT STEPS:**
- [Set next step for component](#set-next-step-for-component)
- [Delete next step for component](#delete-next-step-for-component)
- [Set next step for command](#set-next-step-for-command)
- [Delete next step for command](#delete-next-step-for-command)
- **COMMANDS:**
- [Add command](#add-command)
- [Delete command](#delete-command)
- [Update command](#update-command)


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

Параметры тела запроса

```json
{
    "data": {
        "type": "string",
        "content": [
            "content"
        ]
    },
    "commands": [
        {
            "type": "string",
            "data": "string"
        },
    ],
    "position": {
        "x": "integer",
        "y": "integer"
    }
}
```

Поле              | Тип                       | Описание
------------------|---------------------------|-----------------------------------------------------------
`data`            | object                    | Данные компонента
`data.type`       | string                    | Тип компонента
`data.content`    | [content][type_content][] | Список c данными, специфичными для каждого типа компонента
`commands`        | [command][type_command][] | Список команд
`commands[].type` | string                    | Тип команды
`commands[].data` | string                    | Текст команды
`position`        | object                    | Координаты компонента на поле редактора
`position.x`      | integer                   | Координата X
`position.y`      | integer                   | Координата Y

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


## Delete set of components

[Наверх][toup]

Удаление набора компонентов из структуры бота

```plaintext
POST /api/bots/{botId}/components/del
```

Параметры пути

Поле    | Описание
--------|---------
`botId` | id бота

Параметры тела запроса

```json
{
    "data": [ "componentID" ]
}
```

Поле          | Тип           | Описание
--------------|---------------|------------------------
`data`        | []componentID | Список с id компонентов
`componentID` | integer       | id компонента

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

<details>
    <summary>Пример</summary>
   
`Запрос`

```plaintext
POST /api/bots/64/components/del
```

Тело запроса

```json
{
    "data": [21, 24]
}
```

`Ответ` 

```json
{
    "ok": true
}
```
</details>


- - -


## Update component

[Наверх][toup]

Обновление компонента в структуре бота

```plaintext
PATCH /api/bots/{botId}/components/{compId}
```

Параметры пути

Поле     | Описание
---------|--------------
`botId`  | id бота
`compId` | id компонента

Параметры тела запроса

Запрос должен включать только необходимые для обновления поля.  
Т.е основные поля тела запроса являются необзательными (см. пример).

```json
{
    "data": {
        "type": "string",
        "content": [
            "content"
        ]
    },
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
`data.content` | [content][type_content][] | Список c данными, специфичными для каждого типа компонента
`position`     | object                    | Координаты компонента на поле редактора
`position.x`   | integer                   | Координата X
`position.y`   | integer                   | Координата Y


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

<details>
    <summary>Пример</summary>
   
`Запрос`

```plaintext
PACTH /api/bots/64/components/5
```

Обновление текста и позиции

Тело запроса

```json
{
    "data": {
        "type": "text",
        "content": [
            {
                "text": "Updated Hello Telegram"
            }
        ]
    },
    "position": {
        "x": 141,
        "y": 112
    }
}
```

Обновление только позиции

Тело запроса

```json
{
    "position": {
        "x": 111,
        "y": 222
    }
}
```


</details>


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
            "isMain": true,
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
            "isMain": false,
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

Параметры тела запроса

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
POST /api/bots/{botId}/components/{compId}/commands/{commandId}/next
```

Параметры пути

Поле        | Описание
------------|--------------
`botId`     | id бота
`compId`    | id компонента
`commandId` | id команды

Параметры тела запроса

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
DELETE /api/bots/{botId}/components/{compId}/commands/{commandId}/next
```

Параметры пути

Поле        | Описание
------------|--------------
`botId`     | id бота
`compId`    | id компонента
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

Добавление [команды][type_command] в компоненте

```plaintext
POST /api/bots/{botId}/components/{compId}/commands
```

Параметры пути

Поле     | Описание
---------|--------------
`botId`  | id бота
`compId` | id компонента

Параметры тела запроса

```json
{
    "type": "string",
    "data": "string"
}
```

Поле   | Тип    | Описание
-------|--------|---------------------------------------------
`type` | string | Тип команды
`data` | string | Данные, специфичные для каждого типа команды

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
DELETE /api/bots/{botId}/components/{compId}/commands/{commandId}
```

Параметры пути

Поле        | Описание
------------|--------------
`botId`     | id бота
`compId`    | id компонента
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


## Update command

[Наверх][toup]

Обновление [команды][type_command] в компоненте

```plaintext
PATCH /api/bots/{botId}/components/{compId}/commands/{commandId}
```

Параметры пути

Поле        | Описание
------------|--------------
`botId`     | id бота
`compId`    | id компонента
`commandId` | id команды

Параметры тела запроса

```json
{
    "type": "string",
    "data": "string"
}
```

Поле   | Тип    | Описание
-------|--------|---------------------------------------------
`type` | string | Тип команды
`data` | string | Данные, специфичные для каждого типа команды

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

<details>
    <summary>Пример</summary>  

`Запрос`

```plaintext
POST /api/bots/67/components/5/commands/5
```

Тело запроса

```json
{
    "type": "text",
    "data": "updated"
}
```

`Ответ` 

```json
{
    "ok": true,
}
```
</details>


[//]: # (LINKS)
[type_component]: ../objects.md#component
[type_content]: ../objects.md#content
[type_command]: ../objects.md#command
[toup]: #api-управления-компонентами