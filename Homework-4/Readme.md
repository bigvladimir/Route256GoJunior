## Домашнее задание №4 «Выдача заказов в разной упаковке»
### Цель:

Модифицировать ваш сервис, добавить возможность в ПВЗ выдавать заказы в любой из трех различных упаковках

### Задание:

- **Модифицируйте Go-приложение**, добавьте в метод "Принять заказ от курьера" возможность передавать параметр упаковки
- Всего есть три вида упаковки: пакет, коробка, пленка
- Реализуйте функционал так, чтобы в будущем можно было просто добавить еще один вид упаковки
- При выборе пакета необходимо проверять, что вес заказа меньше 10 кг, если нет, то возвращаем информативную ошибку
- При выборе пакета стоимость заказа увеличивается на 5 рублей
- При выборе коробки необходимо проверить, что вес заказа меньше 30 кг, если нет, то возвращаем информативную ошибку
- При выборе коробки стоимость заказа увеличивается на 20 рублей
- При выборе пленки дополнительных проверок не требуется
- При выборе пленки стоимость заказа увеличивается на 1 рубль

### Дополнительное задание (за него можно получить 10 баллов):

- Опишите архитектуру своего решения любым известным стандартом (например, UML)
- В MR вложите файл с описанием архитектуры
- При выборе стандарта необходимо описать какой был выбран стандарт и дать ссылку на его документацию
- Запрещается использовать генерилки диаграмм, а также инструменты генерации связей между таблицами в БД в качестве описания 
  
## Информация для запуска:
  
Нужно настроить конфиг бд в config  
Чтобы подготовить базу данных, нужно в папке /internal/pkg/db/migrations выполнить команду  
goose postgres "user=USER dbname=DBNAME password=PASSWORD sslmode=disable" up  
  
Примеры команд:  
help  
1 takeorder 1 1 2025-02-02 40 1 bag  
1 takeorder 1 1 2025-02-02 5 1 bag  
