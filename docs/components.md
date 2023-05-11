# Список компонентов

- [Главная](./README.md)

## Components

Список oкомпонентов

| Компонент     | Описание                      |
|---------|-------------------------------|
| [`start`](#start) | Начальный компонент           |
| [`text`](#text)  | Отправка текстового сообщения |


- - -


## Start

[Наверх][toup]

Начальный компонент.

Автоматически добавляется в структуру бота при создании.

<details>
    <summary>Пример</summary>
   
Тело компонента

```json
{
    "data": {
        "type": "start",
        "content": []
    },
    "commands": [],
    "position": {}
}
```
</details>


- - - 


## Text

[Наверх][toup]

Отправка текстового сообщения

```json
{
    "text": "string"
}
```

<details>
    <summary>Пример</summary>
   
Тело компонента

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
    "commands": [],
    "position": {}
}
```
</details>


[//]: # (LINKS)
[toup]: #список-компонентов