# calc-base-api

Это распределённый калькулятор на golang. Он роасполагается на http сервере. Что-бы им воспользоваться, нужно послать POST запрос на 'localhost:8080/api/v1/calculate' С JSON типа '{"expression": "выражение"}', и он пришлёт ответ с JSON типа {"id": "ID выражения"} и код 201. Это означает, что выражение находится в очереди на решение. Он также может выдать 2 кода ошибок:

    422 - Если выражение не соответствуют требованиям приложения (Нелегальные символы, или выражение не решаемо.)
    500 - Если в теле запроса есть ошибки (Запрос не оформлен по правилам JSON)

Чтобы просмотреть выражения, которые находятся у сервера, нужно послать GET запрос на 'localhost/api/v1/expressions', сервер пришлёт список со всеми выражениями.


## Запуск

1. Cклонировать репозиторий (Нужна программа git)
```bash
git clone https://github.com/Se623/calc-base-api
```
2. Перейти в директорию программы
```bash
cd ./calc-base-api
```
3. Установить зависимости
```bash
go install ./cmd
```
4. Запустить калькулятор
```bash
go run ./cmd
```

Сервер распологается на порту 8080.

## Примеры

### Пример 1 (Обычное выражение)
Запрос:\
Bash(Linux): `curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "2+2*2"}'`\
Cmd: `curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"2+2*2\"}"`\

Ответ: `{"id": "0"}`

### Пример 1 (Ошибка)
Запрос:\
Bash(Linux): `curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "***5***"}'`\
Cmd: `curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"***5***\"}"`\

Ответ: `Error: Invalid Input` (Выражение не покажется в списке)

## Пример просмотра выражений
Запрос:\
Bash(Linux): `curl --location 'localhost:8080/api/v1/expressions'`\
Cmd: `curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"2+2*2\"}"`\

Ответ: `{"expressions":[{"id":0,"status":"Queued","result":-1}]}` (Может быть другой ответ в зависимости от решаемых выражений)

