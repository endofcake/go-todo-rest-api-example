  CREATE DATABASE todoapp;
  CREATE ROLE docker WITH LOGIN PASSWORD 'qwerty';
  GRANT ALL PRIVILEGES ON DATABASE todoapp TO docker;