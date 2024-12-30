# Workshop #1

# Задача

Создать сервис отзывов на товары, который будет иметь два HTTP API:

1. POST products/{id}/reviews - Создать отзыв.
2. GET products/{id}/reviews - Получить все отзывы по товару

## Создание отзыва

Все параметры обязательны.

```http request
POST products/{id}/reviews
{
  "sku": 100,
  "comment": "Отлично",
  "user_id": "35fe1358-a03a-4bba-b2aa-b0b1c49264d2"
}
```

## Получение отзывов

```http request
GET products/{id}/reviews
{
  "reviews": [
     {
        "sku": 100,
        "comment": "Отлично",
        "user_id": "35fe1358-a03a-4bba-b2aa-b0b1c49264d2"
      }
   ]
}
```

# Clean Arch
![](docs/clean_arch.png)

![](docs/clean-arch-impl.png)

# Полезные ссылки

+ Effective Go - https://go.dev/doc/effective_go
+ Go Wiki: Common Mistakes https://go.dev/wiki/CommonMistakes
+ Go Wiki: Go Code Review Comments https://go.dev/wiki/CodeReviewComments
+ Структура проекта - https://github.com/golang-standards/project-layout/blob/master/README_ru.md
+ Uber style guide - https://github.com/uber-go/guide/blob/master/style.md
