# Домашнее задание: 
Работу продолжаем в репозитории домашнего задания 2.

1.  __Сохранение данных в бд__
Подключить работу с базой данных, добавить конфиг подключения, инициализацию коннекта.

Работу с мапкой/файлом оставляем неизменными, доменные структуры можно переиспользовать

Написать CRUD операции для работы с бд
Должны быть реализованы методы записи и чтения данных простой системы  хранения ПВЗ

###### _Подсказка:_
_Вам могут пригодится следующие методы_
- _GetByID_
- _List_
- _Update_
- _Create_
- _Delete_

_Так же можете реализовать те методы, которые вы делали с файлом._

2. __Разработать HTTP сервер.__

- Необходимо реализовать HTTP сервер, который будет работать с методами базы данных, реализованными в 1 пункте. 

- Методы должны позволять манипулировать данными(create,read,update,delete) для системы хранения пвз

- Методы должны быть выполнены в restful стиле. Необходимо корректно обрабатывать все коды ошибок 

- Входящие параметры должны быть представлены либо в формате json либо в query параметрах(зависит от метода)

- Сервис должен использовать порт 9000

3. __В ридми приложить curl запросы, на каждую ручку. Запросы должны быть валидными и возвращать нужный код ответа__
4. __Необходимо реализовать middleware, который будет логгировать поля POST,PUT,DELETE запросов__


###### _Подсказка:_
_Посмотрите на результат выполнения дз 2. Сервис должен делать похожий flow, но используя бд как хранилище и http как интерфейс взаимодействия с пользователем_
## Дополнительно:
1. Поддержать https. Можно использовать самоподписный сертификат от Let's Encrypt
2. Реализовать middleware с basic auth. Юзер/пароль можно задать как в конфиге, так и хранить в базе(создать круд юзеров)


## Ограничения дз:
- Нельзя использовать orm или sql билдеры
- Для реализации http сервера можно использовать как net/http так и gin/fasthttp и прочее
- Коды ошибок должны соответствовать поведению сервиса. Хендлеры, которые отдают только 500 в случае ошибки - не принимаются
- В хендлерах должна быть базовая валидация данных, соответствующая бизнес-логике


# Информация для запуска:

Перед началом работы:  
1) Локально запустить postgres и выбрать любую базу данных где не занята таблица pvz  
2) Задать конфиг бд в /config/bdcfg.json
3) Чтобы подготовить базу данных, нужно в папке /internal/sqlserver/db/migrations выполнить команду  
goose postgres "user=USER dbname=DBNAME password=PASSWORD sslmode=disable" up  
    
Код из прошлого дз запускается под флагом: 2  
Для третьего дз нужно запустить программу с флагом: 3  
Нарпример   
go build -o pvz  
./pvz 3  
  
Curl'ы:  
  
Наполнить базу данных для тестов можно так  
curl localhost:9000/pvz -X POST -d '{"Name":"velikolepnii","Adress":"ulitsa Pushkina","Contacts":"+345345345"}' -i -u login:password  
curl localhost:9000/pvz -X POST -d '{"Name":"zachetnii","Adress":"prospekt Lenina","Contacts":"+5345245235"}' -i -u login:password  
curl localhost:9000/pvz -X POST -d '{"Name":"luchshii","Adress":"406540 stroenie 65","Contacts":"+53425324"}' -i -u login:password  
  
curl localhost:9000/pvz -X POST -d '{"Name":"luchshii3","Contacts":"+53425334"}' -i -u login:password  
  
  
Остальные запросы  
curl localhost:9000/pvz/1 -X GET -i -u login:password  
curl localhost:9000/pvz/2 -X GET -i -u login:password  
curl localhost:9000/pvz/12323 -X GET -i -u login:password  
curl localhost:9000/pvz/-2 -X GET -i -u login:password  
  
curl localhost:9000/pvz -X PUT -d '{"ID":2,"Name":"antonova","Adress":"pereulok rodnoi","Contacts":"+353453454"}' -i -u login:password  
curl localhost:9000/pvz -X PUT -d '{"ID":15,"Name":"antonova","Adress":"pereulok rodnoi","Contacts":"+353453454"}' -i -u login:password  
  
curl localhost:9000/pvz/1 -X DELETE -i -u login:password  
curl localhost:9000/pvz/1 -X DELETE -i -u login:password  
curl localhost:9000/pvz/100 -X DELETE -i -u login:password  
  
curl localhost:9000/pvz/3 -X PATCH -d '{"Name":"neeeeneee","Adress":"pereulok rodnoi","Contacts":"+353434553454"}' -i -u login:password  
curl localhost:9000/pvz/3 -X GET -i -u login:password  
curl localhost:9000/pvz/100 -X PATCH -d '{"Name":"nuurrrr","Adress":"pereulok rodnoi","Contacts":"+353434553454"}' -i -u login:password  
  
чтобы проверить работу https можно положить в папку/config/tls свои файлы server.crt и server.key или находясь в этой папке прописать  
openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt  
  
curl -k https://localhost:9001/pvz/3 -X GET -i -u login:password  
  
После всего можно выполнить  
goose postgres "user=USER dbname=DBNAME password=PASSWORD sslmode=disable" down  
чтобы вернуть базу данных к исходному состоянию  