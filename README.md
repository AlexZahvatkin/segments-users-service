# Тестовое задание
# Dynamic users segmentation service

## Getting Sarted

### Запуск через docker
1. Переименовать файл .env.example => .env () ИЛИ можно запустить без пунктов 1 и 2 (в docker-compose прописаны параметры по умолчанию)
2. При необходимости поменять значения переменных среды в .env
3. Выполнить команду
  ```sh
  make run
  ```
### Локальный запуск
В случае, если возникли какие-либо проблемы с запуском, можно запустить проект локально. Для этого нужно:
1. Переименовать файл .env.example => .env
2. При необходимости поменять значения переменных среды в .env
3. Установить PostgreSQL
4. Выполнить команду с нужными параметрами:
  ```sh
  make createandmigrate DB_USER=postgres DB_PASSWORD=password DB_HOST=localhost DB_PORT=5432 DB_DATABASE=segments DB_SSLMODE=disable
  ```
## Используемые библиотеки и технологии
Проект использует следующие библиотеки и технологии:
- PostreSQL (для хранения сущностей и отношений между ними)
- Docker (для запуска сервиса)
- Swagger (для документации)
- chi (библиотека для роутинга)
- sqlc (генерация моделей и кода взимодействия с БД)
- golang/testify (для написания тество)
- golang/migrate (для миграций (необходимы в тестах))

## Примеры запросов и ответов 
Документация по API доступна через Swagger: http://localhost:8080/swagger/index.html#/ (порт может отличаться в зависимости от настроек .env)

### Создание нового пользователя
```
POST /v1/users
```
#### Описание: 
Создает нового пользователя с заданым именем 
#### Тело запроса:
```
{
  "name": "Alexander"
}
```
#### Пример ответа:
```
{
  "id": 13,
  "name": "Alexander",
  "created_at": "2023-08-31T20:22:30.900337Z",
  "updated_at": "2023-08-31T20:22:30.900337Z"
}
```

### Удаление пользователя по заданному ID
```
DELETE /v1/users/{userId}
```
#### Описание: 
Удаляет пользователя с ID, преданным в URL
#### Пример ответа:
```
Status: 200 OK
```

### Создание сегмента 
```
POST /v1/segments
```
#### Описание: 
Создает сегмент с заданным именем и описанием. Имя сегмента обязательный параметр. 
Также в запросе можно передать процент пользователей, которые должны попасть в данный сегмент. 
Если процент был передан, после создания сегмента данный сегмент будет добалвен даннному проценту пользователей (пользователи выбираются случайно)
#### Тело запроса:
```
{
  "name": "segment_name",
  "description": "short description",
  "percent": 30
}
```
#### Пример ответа:
```
{
  "segment": {
    "name": "SEGMENT_NAME",
    "description": "short description",
    "created_at": "2023-08-31T20:27:29.357976Z",
    "updated_at": "2023-08-31T20:27:29.357976Z"
  },
  "added_users_ids": [
    1,
    2,
    7,
    8
  ]
}
```

### Удаление сегмента
```
POST /v1/segments?name={name}
```
#### Описание:
Удаляет сегмент с данным именем (slug-ом). 
При создании сегмента все имена приводятся в единый формат: пробелы заменяются на _, строчные буквы меняются на заглавные.
Но при удалении данного форматирования не происходит (сделано для того, чтобы не удалить сегмент случайно), поэтому важно передать сегмент в верном формате.
#### Пример ответа:
```
Status: 200 OK
```

### Добавление и удаление сегментов для пользователя
```
POST v1/segments/assign/{userId}
```
#### Описание:
Добавляет сегменты, переданные в параметре to_add, для данного пользоватея. Удаляет сегменты, переданные в параметре to_delete, для данного пользователя.
Приоритет отдается удалению. Поэтому, если сегмент был передан в обоих параметрах, он все равно будет удален.
#### Тело запроса:
```
{
  "to_add": [
    "AVITO_DISCOUNT_50"
  ],
  "to_delete": [
    "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"
  ]
}
```
#### Пример ответа:
```
[
  {
    "user_id": 1,
    "segment_name": "AVITO_DISCOUNT_50",
    "created_at": "2023-08-31T20:56:57.784489Z",
    "updated_at": "2023-08-31T20:56:57.784489Z",
    "expire_at": "0001-01-01T00:00:00Z"
  }
]
```

### Задание пользователю сегмента с определенным TTL
```
POST v1/segments/ttl/{userId}
```
#### Описание: 
Добавляет пользователю сегмент с определенным TTL(заданным в часах). По истечение заданного времени, данный сегмент считается неактивным для данного пользователя. 
#### Тело запроса: 
```
{
  "segment_name": "AVITO_PERFORMANCE_VAS",
  "ttl": 10
}
```
#### Пример ответа:
```
{
  "user_id": 1,
  "segment_name": "AVITO_PERFORMANCE_VAS",
  "created_at": "2023-08-31T21:03:04.319455Z",
  "updated_at": "2023-08-31T21:03:04.319455Z",
  "expire_at": "2023-09-01T07:03:04.319455Z"
}
```

### Получение активных сегментов для данного пользователя:
```
GET v1/segments/{userId}
```
#### Описание:
Получает список активных сегментов для данного пользователя. Сегменты, у которых вышел TTL, для данного пользователя - возвращены не будут.
#### Пример ответа:
```
[
  "TEST_NEW_ADD_10",
  "SEGMENT_NAME",
  "AVITO_DISCOUNT_50",
  "AVITO_PERFORMANCE_VAS"
]
```

### Получение истории добавления/удаления пользователя в сегмент
```
GET v1/segments/history/{userId}?from={from}&to={to}
```
#### Описание:
Возвращает CSV, в котором перечислена информация о том, когда для данного пользователя были удалены/добавлены сегменты в заданном промежутке времени
#### Пример ответа:
```
1,AVITO_VOICE_MESSAGES,deleted,2023-08-31 20:43:42
1,AVITO_PERFORMANCE_VAS,deleted,2023-08-31 20:56:57
1,AVITO_DISCOUNT_30,deleted,2023-08-31 20:56:57
1,AVITO_DISCOUNT_50,inserted,2023-08-31 20:56:57
1,AVITO_PERFORMANCE_VAS,inserted,2023-08-31 21:03:04
```

## Проблемы, с которыми столкнулся, и их решения
1. Для хранения ID в БД можно было использовать UUID. Но был выбран формат bigserial, так как это облегчает использование API
2. Формат хранения даннх не был описан в задании. Для сущности пользователей в качестве первичного ключа был создан суррогатный ключ. Для хранения сегментов в качестве ключа было выбрано их имя, так как имя сегмента однозначно задает его цель, а создание сегмента с уже существующим именем может привести к неопределенности.
3. API для создания и удаления пользователей по заданию не требовалось. Но я решил создать их, так как это облегчает работу с API
4. В задании сказано, что важно, чтобы сегменты не перетирались. Для этого при добавлении в БД все сегменты приводятся к одному формату: пробелы заменяются на _, буквы переходят в верхний регистр. Такой формат был выбран, так как примеры в задании используют именно его.
5. Cитуация, когда мы добавляем сегмент пользователю, у которого в данном сегменте стоит TTL: мне показалось логично в данной ситуации выполнить upsert. Пользователю добавляется данный сегмент, если у него его нет. В ситуации, если у пользоватея есть сегмент c TTL - TTL становится пустым.
6. Для хранения истории добавления/удаления сегментов пользователей - создал вспомогательную таблицу истории, куда с помощью триггеров записываются данные.
7. Логика при добалении пользователю сегмента, в котором он уже состоит: в данном случае никаких ошибок и уведомлений не происходит.
8. При добавлении пользователю сегмента с заданным TTL: TTL передается в часах. TTL не обеспеичивал бы достаточной гибкости, а в минутах был бы неудобен для использования. Тем не менее, при необходимости формат TTL можно легко поменять.
9. Для обеспечения функциональности TTL в БД создано поле expire_at. Оно показывает, когда данный сегмент для данного пользователя можно считать недействительным. В базу данных можно добавить EVENT, который будет проходится раз по таблице (например, раз в месяц), и удалять сегменты для пользователя, у которых уже вышел TTL.
10. В целом, в задании не совсем ясно указано: должен ли в TTL передаваться в формате конкретного времени, когда данный сегмент становится невалидным для пользователя, или же в формате временного периода. Я выбрал второй вариант.
11. При добавлении сегментов пользователю происходит проверка наличия переданных сегментов и пользователя в БД. Это может негативно сказываться на производительности, однако дает возможность дать более точный ответ клиенту, почему его запрос вернулся с ошибкой. Однако, и это не дает полной гарантии: ведь сегмент или пользователь могли быть удалены после того, как мы получили запрос, но перед тем, как мы выполнили проверку. Но данная ситуация довольно маловероятна. 
