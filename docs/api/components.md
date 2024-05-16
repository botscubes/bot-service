# API управления компонентами

- [Главная](../README.md)

## Methods

- **Components:**
    - [Get components](#get-components)
    - [Add component](#add-component)
    - [Delete component](#delete-component)
    - [Update component data](#update-component-data)
    - [Update component position](#update-component-position)
- **Connections:**
    - [Add connection](#add-connection)
    - [Delete connection](#delete-connection)



- - -

## Get components

[Наверх][toup]

Получение компонентов бота

```plaintext
GET /api/bots/{botId}/groups/{groupId}/components
```

Параметры пути

- botId: integer - id бота
- groupId: integer - id группы компонентов

#### Ответ

В случае успеха статус 200 с телом ответа: 

```plaintext
[
    ..."components"
]
```

Структура компонента:

```plaintext
{
    "id": "integer",
    "type": "string",
    "data": "object",
    "path": "string",
    "outputs": {
        "nextComponentId": "integer",
        ...
    },
    "connectionPoints": {
        "<pointId: string>": {
            "sourceComponentId": "integer",
            "sourcePointName": "string",
            "relativePointPosition": {
                "x": "integer",
                "y": "integer"
            }
        }
        ...
    },
    "position": {
        "x": "integer",
        "y": "integer"
    }
}
```

- - -


## Add component

[Наверх][toup]

Добавление нового компонента в структуру бота

```plaintext
POST /api/bots/{botId}/groups/{groupId}/components
```

Параметры пути

- botId: integer - id бота
- groupId: integer - id группы компонентов

Параметры тела запроса

```json
{
    "type": "string",
    "position": {
        "x": "integer",
        "y": "integer"
    }
}
```
#### Ответ

В случае успеха статус 201 с телом ответа: 

```json
{
    "id": "integer"
}
```

где id: integer - id вновь созданного бота.


<details>
    <summary>Пример</summary>
   
`Запрос`

```plaintext
POST /api/bots/64/groups/1/components
```

Тело запроса

```json
{
    "type": "message",
    "position": {
        "x": 141,
        "y": 112
    }
}
```

`Ответ` 

```json
{
    "id": 7
}
```
</details>


- - -


## Delete component

[Наверх][toup]

Удаление компонента из структуры бота

```plaintext
DELETE /api/bots/{botId}/groups/{groupId}/components/{compId}
```

Параметры пути

- botId: integer - id бота
- groupId: integer - id группы компонентов
- compId: integer - id компонента


#### Ответ

В случае успеха статус 204 без тела ответа.


- - -


## Update component data

[Наверх][toup]

Обновление данных компонента

```plaintext
PATCH /api/bots/{botId}/groups/{groupId}/components/{compId}/data
```

Параметры пути

- botId: integer - id бота
- groupId: integer - id группы компонентов
- compId: integer - id компонента


Параметры тела запроса

```plaintext
{
    "<property name>": "any",
    ...
}
```
#### Ответ

В случае успеха статус 204 без тела ответа.


- - -


## Update component position

[Наверх][toup]

Обновление позиции компонента


```plaintext
PATCH /api/bots/{botId}/groups/{groupId}/components/{compId}/position
```

Параметры пути

- botId: integer - id бота
- groupId: integer - id группы компонентов
- compId: integer - id компонента

Параметры тела запроса

```json
{
    "x": "integer",
    "y": "integer"
}
```

#### Ответ

В случае успеха статус 204 без тела ответа.


- - - 


## Add connection

[Наверх][toup]

Добавление соединения между компонентами.


```plaintext
POST /api/bots/{botId}/groups/{groupId}/connections
```

Параметры пути

- botId: integer - id бота
- groupId: integer - id группы компонентов

Параметры тела запроса

```json
{
    "sourceComponentId": "integer",
    "sourcePointName": "string",
    "targetComponentId": "integer",
    "relativePointPosition": {
        "x": "integer",
        "y": "integer"
    }
}
```

где
- sourceComponentId - id компонента, от которого будет переход
- sourcePointName - имя точки компонента, от которого будет переход
- targetComponentId - id компонента, к которому будет переход
- relativePointPosition - расположение точки относительно компонента, от которого будет переход

#### Ответ

В случае успеха статус 201 без тела ответа.


- - - 


## Delete connection

[Наверх][toup]

Добавление соединения между компонентами.


```plaintext
DELETE /api/bots/{botId}/groups/{groupId}/connections
```

Параметры пути

- botId: integer - id бота
- groupId: integer - id группы компонентов

Параметры тела запроса

```json
{
    "sourceComponentId": "integer",
    "sourcePointName": "string"
}
```

где
- sourceComponentId - id компонента, от которого будет переход
- sourcePointName - имя точки компонента, от которого будет переход

#### Ответ

В случае успеха статус 204 без тела ответа.








[//]: # (LINKS)
[type_component]: ../objects.md#component
[type_content]: ../objects.md#content
[type_command]: ../objects.md#command
[toup]: #api-управления-компонентами
