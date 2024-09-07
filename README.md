# Файлы для итогового задания

Настройки для тестов:
var Port = 7540                       //TODO_PORT из файла .env
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6IiJ9.sgOfht2gj7-D_-HdrvnboC_ehEq-UGdzQQpl3devG_s`


Адрес в браузере: http://localhost:7540/

Пароль: 1324657980         //TODO_PASSWORD из файла .env

Проект Todo List. 
Может получать, изменять, вносить, подтверждать задачи. И также можно осуществлять поиск через поисковую строку.
Все задачи хранятся в базе данных.
Возможна авторизация по паролю.

В папке hendler хранятся все используемые хэндлеры. 
На данном этапе также там лежат вспомогательные функции.

В bd лежит логика работы с базой данных.

В entity лежат структуры, которые используются в процессе реализации.

В папке middlewares лежит мидварь для проверки токена, используется только для проверки аутентификации.

В папке web лежит весь фронт.

Создаю и запускаю контейнер: 
 docker build  -t maliven1/todo-list:v1.0.0 .   
 docker run -p 7540:7540 -v $PWD/schedule.db:/schedule.db  maliven1/todo-list:v1.0.0     // для винды
