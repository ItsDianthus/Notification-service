# Шаблон Go-проекта для домашних заданий

Шаблон для домашних заданий [Академии Бэкенда 2024](https://edu.tinkoff.ru/all-activities/courses/870efa9d-7067-4713-97ae-7db256b73eab).

Цель данного репозитория – познакомить вас с процессом разработки приложений на Go с использованием наиболее распространенных практик, инструментов и библиотек.

## Структура проекта

Проект содержит в себе следующие компоненты:

- **`cmd/`** – директория с исполняемыми файлами, здесь находятся два сервиса:
  - `bot/` – точка входа в сервис Telegram-бота.
  - `scrapper/` – точка входа в сервис для отслеживания изменений на сайтах.

  Каждая папка содержит `main.go` и представляет отдельное исполняемое приложение.

- **`internal/`** – директория с внутренними пакетами, которые не предназначены для использования вне проекта.
  - **`bot/`** – реализация логики бота:
    - `application/` – юзкейсы и сценарии взаимодействия пользователя с ботом.
    - `config/` – конфигурация и параметры запуска бота.
    - `domain/` – доменные сущности и интерфейсы, описывающие бизнес-логику бота.
    - `infrastructure/` – инфраструктурный слой (работа с Telegram API, хранением данных и т.д.).

  - **`scrapper/`** – реализация логики скраппера:
    - `application/` – юзкейсы по управлению ссылками, подписками и проверкой обновлений.
    - `config/` – конфигурация для запуска сервиса скраппера.
    - `domain/` – модели и интерфейсы для работы со ссылками и подписками, в том числе enum с возможными источниками ссылок для обработки
    - `infrastructure/` – взаимодействие с сетью, БД, логированием и другими внешними сервисами.

  - **`api/openapi/`** – сгенерированные клиенты по OpenAPI-спецификации:
    - `bot_api/` – клиент для взаимодействия скраппера с ботом.
    - `scrapper_api/` – клиент для взаимодействия бота со скраппером.

- **`api/`** – директория с исходными спецификациями API:
  - `openapi/` – YAML/JSON-файлы для генерации HTTP-клиентов и серверов.
  - `proto/` – файлы .proto.

- **`pkg/`** – общие утилиты, которые могут быть переиспользованы в других проектах:
  - `slogger/` – обёртка над логированием (`slog`), реализующая единый подход к логированию по всему проекту.

- **`config/`** – глобальные конфигурационные файлы для всего проекта (YAML).

- **`bin/`** – (опционально) директория для сборки бинарников.

- **`build/`** – директория для Docker-файлов (и др.)
