# Итоговое задание
## Описание
В итоговом задании реализован веб-сервер на Go, который реализует функциональность простейшего планировщика задач. Это аналог TODO-листа.
В задании реализованы добаление, редактирование, удаление, вывод задач, отметка о выполнение задачи и авторизация

## Задания со звездочкой
Добавлены переменные окружения TODO_PORT, TODO_DBFILE, TODO_PASSWORD
Реализованы правила повторения задач для недель и месяцев
Реализован поиск по дате, названию и комментарию
Реализована авторизация
Создан Dockerfile READE.md

## Запуск проекта локально
### Клинирование репозитория
```
https://github.com/egor1221/final_project.git
cd final_project
```
Если настроен ssh
```
git@github.com:egor1221/final_project.git
cd final_project
```
### Добавление переменных окружения
Добавляются в файл .env или при помощи команды
```
export ПЕРЕМЕННАЯ=ЗНАЧЕНИЕ
```
### Запуск
```
cd cmd
go run .
```
После запуска перейдите по ссылке [](http://localhost:7540)
### Настройка тестов
Перед запуском убедитесь, что в файле settings.go указаны верные данные 
```
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `Ваш токен при генерации`
```

### Запуск тестов
Запуск всех тестов
```
go test ./tests
```
Запуск тестов по отдельности
```
go test -run ^TestApp$ ./tests
go test -run ^TestDB$ ./tests
go test -run ^TestNextDate$ ./tests
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestTask$ ./tests
go test -run ^TestEditTask$ ./tests
go test -run ^TestDone$ ./tests
go test -run ^TestDelTask$ ./tests
```

### запуск через Docker
```
docker build --tag app:v1 . 

docker run -d -p 7540:7540 app:v1 
```