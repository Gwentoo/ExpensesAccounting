# Аналитика расходов
## Запуск
### 1. Создать файл `.env` в корне проекта со следующим содержимым:
 ```yaml
#DB
DB_HOST=postgres
DB_PORT=5432
POSTGRES_EXTERNAL_PORT=5434
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=expense
DB_SSLMODE=disable

  #Redis
REDIS_ADDR=redis:6379
REDIS_EXTERNAL_PORT=6379

  #Superset
SUPERSET_USERNAME=admin
SUPERSET_PASSWORD=admin
SUPERSET_EXTERNAL_PORT=8088

  #Auth service
AUTH_HOST=auth-service
AUTH_SERVICE_PORT=50051

  #Expenses service
EXPENSES_HOST=expenses-service
EXPENSES_SERVICE_PORT=50052

  #Gateway service
GATEWAY_EXTERNAL_PORT=8080

  #Google
GOOGLE_PASS=YOUR_GOOGLE_PASS (КАК ЕГО ПОЛУЧИТЬ В ПУНКТЕ 2)
GOOGLE_FROM=YOUR_GOOGLE_GMAIL

  #JWT
JWT_SECRET_KEY=YOUR_SECRET_KEY
```
Изменить `GOOGLE_PASS`, `GOOGLE_FROM`, `JWT_SECRET_KEY`
### 2. Получение GOOGLE_PASS
* [Пароли приложений](https://myaccount.google.com/apppasswords?continue=https://myaccount.google.com/security)
* Создать приложение и получить код вида "abcd efgh ijkl mnop"
* Вставить в поле `GOOGLE_PASS` без пробелов

### 3. Получение SUPERSET_DASHBOARD_ID
* Из корня проекта выполняем команды
```Exec
make build
make start
```
* В браузере заходим на `localhost:8088`
* Логинимся, используя `admin` и `admin`
* В Dashboards импортируем архив `dashboard_export_20260527T000604.zip`
* В Datasets выбираем базу данных и в редактировании меняем XXXXXXXX на пароль
* Возвращаемся в Dashboards, выбираем нужный, жмем на 3 точки сверху и выбираем `Embed Dashboard`
* Вставляем uuid в `/web-app/.env` в поле `REACT_APP_DASHBOARD_OVERVIEW`
* Из корня проекта выполняем команду 
```Exec
docker-compose up --build
```

### 4. Запуск веб-приложения
`Должен быть установлен Node.js`
* Переходим в директорию `/web-app`
* Выполняем команды
```Exec
npm install
npm start
```