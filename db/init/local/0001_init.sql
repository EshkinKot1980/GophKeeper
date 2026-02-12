CREATE USER keeperuser
    PASSWORD 'keeperpass';

CREATE DATABASE gophkeeper
    OWNER 'keeperuser'
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';
