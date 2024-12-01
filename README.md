EMtest

Запросы:
GET /info 
Получить данные о песни по названию песни и имполнителю в параметрах запроса

GET /songs 
Получить песни с фильтрацией по всем полям в параметрах запроса и пагинацией.
Если параметры page, limit переданы не были, пагинация выключается

GET /songs/{id}
Получить песню по id

GET /songs/{id}/text
Получить текст песни по id с пагинацией через параметры запроса.
значение limit по умолчанию - 4 (Возвращается по 4 строки текста).
Если параметр page не был передан, возвращается весь текст

POST /songs
Создание песни через тело запроса

POST /songs/{id}/info
Создание информации о _СУЩЕСТВУЮЩЕЙ_ песне через тело запроса
Для одной песни разрешено иметь одну запись с ниформацией о ней

DELETE /songs/{id}
Удаление песни

DELETE /songs/{id}/info
Удаление записи с ниформацией о песне

PATCH /songs/{id}
Обновление песни

PATCH /songs/{id}/info
Обновление записи с ниформацией о песне

База данных:
Используется БД postgres.
Схема имеет в себе две таблицы, песни и их информация.