# Описание объектов 

- [Главная](./README.md)
- [Типы компонентов][types_of_components]

## Objects

- [Component](#component)
- [Content](#content)
- [Command](#command)
- [Types of components](#types-of-components)



## Component

[Наверх][toup]

Содержит информацию о компоненте  
Поле `keyboard` пока НЕ ИСПОЛЬЗУЕТСЯ. 

```json
{
    "id": "integer",
    "data": {
        "type": "string",
        "content": [
            "content"
        ]
    },
    "keyboard": "keyboard",
    "commands": [
        "command"
    ],
    "nextStepId": "integer",
    "isMain": "bool",
    "position": {
        "x": "integer",
        "y": "integer"
    }
}
```

Поле           | Тип                   | Описание
---------------|-----------------------|-----------------------------------------------------------
`id`           | integer               | id компонента
`data`         | object                | Данные компонента
`data.type`    | string                | [Тип][types_of_components] компонента
`data.content` | [content](#content)[] | Список c данными, специфичными для каждого типа компонента
`keyboard`     | keyboard              | Структура клавиатуры (НЕ ИСПОЛЬЗУЕТСЯ)
`commands`     | [command](#command)[] | Список команд
`nextStepId`   | integer               | id следующего шага (компонента)
`isMain`       | bool                  | Определяет, является ли компонент начальным
`position`     | object                | Координаты компонента на поле редактора
`position.x`   | integer               | Координата X
`position.y`   | integer               | Координата Y

<details>
    <summary>Примеры</summary>

```json
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
}
```

```json
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
            "nextStepId": 1
        },
        {
            "id": 3,
            "type": "text",
            "data": "abc",
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
```
</details>

## Content

[Наверх][toup]

Соджержит данные, специфичные для каждого типа компонента  

Структуры данных смотреть в

--> [Типы компонентов][types_of_components] <--


- - - 


## Command

[Наверх][toup]

Содержит информацию о команде

```json
{
    "id": "integer",
    "type": "string",
    "data": "string",
    "componentId": "integer",
    "nextStepId": "integer"
}
```

| Поле          | Тип     | Описание                        |
|---------------|---------|---------------------------------|
| `id`          | integer | id команды                      |
| `type`        | string  | Тип команды                     |
| `data`        | string  | Текст команды             |
| `componentId` | integer | id компонента ?!?!              |
| `nextStepId`  | integer | id следующего шага (компонента) |

Тип компонента на данный момент есть только `text`

<details>
    <summary>Примеры</summary>

```json
{
    "id": 1,
    "type": "text",
    "data": "First button",
    "componentId": 2,
    "nextStepId": null
}
```

```json
{
    "id": 2,
    "type": "text",
    "data": "Second button",
    "componentId": 2,
    "nextStepId": 1
}
```
</details>


- - -


## Types of components

[Наверх][toup]

Список типов компонентов и описание их структуры

| Тип     | Описание                      |
|---------|-------------------------------|
| `start` | Начальный компонент           |
| `text`  | Отправка текстового сообщения |

**Описание структуры компонентов**

#### start

Начальный компонент

`No data`

#### text

Отправка текстового сообщения

```json
{
    "text": "string"
}
```


[//]: # (LINKS)
[types_of_components]: #types-of-components
[toup]: #описание-объектов 