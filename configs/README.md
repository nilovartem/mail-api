# configs

### Конфигурационный файл

 mailbox - папка, где хранится почта, пример в test/data, название папки может быть любым.

 Ее структура: mailbox/username:

 Например:

 mailbox/ivanov@mail.ru

 mailbox/petrov@mail.ru

    "mailbox":"путь к папке с почтой",
    "ttl": время жизни ссылки, нужно указывать с единицей измерения ("1s", "15m"),
    "log_level": уровень логирования, например: "debug",
    "bind_address": адрес в формате "host:port",
    "users":{
        "имя пользователя":"пароль",
    }
