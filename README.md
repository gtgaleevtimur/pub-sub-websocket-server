# pub-sub-websocket-server

## Задание.
Необходимо реализовать сервер в виде консольного приложения, который
регистрирует подключения клиентских приложений (каждому клиенту при
подключении присваивается уникальный идентификатор) по веб-сокетам в группах
(hub) не более n - клиентов
При достижении лимита (n) подключений создается новый hub
При выполнении в консоли сервера команды “send --hub” (где параметр --hub -
номер hub) осуществляется broadcast рассылка произвольного сообщения все
клиентам указанного hub
При выполнении в консоли сервера команды “sendс --id” (где параметр --id -
идентификатор клиента) осуществляется отправка сообщения сообщения
конкретному клиенту.

## Реализация.
Сервер запускается по дефолту на localhost:8080 .

Для передачи сервису настроек количества соединений в одном хабе при запуске
с флагом задать int значение по дефолту стоит 10.

>go run -n 10

## Взаимодействие или возможные команды в консоли:
Для рассылки сообщения пользователям одного хаба:
>send 0 hello hub!

Для рассылки сообщения одному пользователю:
>sendc 0 hello user!

Для выходя из приложения дважды повторить Interrupt.

Примечание:
единственный хендлер для установления соединения с приложением клиента.
- GET / — для установления соединения с приложением клиента;

Для тестов использовался https://www.piesocket.com/websocket-tester .

Затраченное время: 5 часов.