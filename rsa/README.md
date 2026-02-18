## Сертификат

Директория содержит сертификат `server.crt` и приватный ключ `server.key` для обеспечения `HTTPS/TLS` соединения.

### Генерация самоподписанного сертификата (Ubuntu)

Для локального и dev окружения можно использовать самоподписанный сертификат

``` bash
openssl req -newkey rsa:4096 \
  -x509 \
  -sha256 \
  -days 365 \
  -nodes \
  -out rsa/tls.crt \
  -keyout rsa/tls.key \
  -subj "/C=RU/ST=MO/O=MyOrg/OU=IT/CN=www.example.com"

```

## JWT

Директория так же содержит публичный и приватный ключи (`jwt-pub.pem`,  `jwt-priv.pem`) для создания токена.

### Генерация ключей (Ubuntu)

``` bash 
openssl genrsa -out rsa/jwt-priv.pem 4096
openssl rsa -in rsa/jwt-priv.pem -pubout -out rsa/jwt-pub.pem
```

## Примечание

Указанные команды выполняются в корне поекта.
Пути к ключам и сертификату можно изменить через файл конфигурации
`tls_cert` , `tls_key`, `jwt_priv`, `jwt_pub` либо через переменные среды `TLS_CERT`, `TLS_KEY`, `JWT_PRIV`, `JWT_PUB`.