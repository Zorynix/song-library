# 🎶 Онлайн библиотека песен

## 📌 Описание задания

Реализовать REST-сервис онлайн библиотеки песен с возможностью:

### 📖 Функционал:

1. **Получение списка песен**

   - Фильтрация по всем полям ( название группы, название песни и т.д. )
   - Пагинация результатов

2. **Получение текста песни**

   - Пагинация по куплетам

3. **Удаление песни**

4. **Редактирование данных песни**

5. **Добавление новой песни**\
   Формат JSON-запроса:

   ```json
   {
     "group": "Muse",
     "song": "Supermassive Black Hole"
   }
   ```

   При добавлении необходимо сделать запрос к внешнему API (Swagger спецификация описана ниже).

---

### 🌐 Внешний API для обогащения данных:

```yaml
openapi: 3.0.3
info:
  title: Music info
  version: 0.0.1
paths:
  /info:
    get:
      parameters:
        - name: group
          in: query
          required: true
          schema:
            type: string
        - name: song
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Ok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SongDetail'
        '400':
          description: Bad request
        '500':
          description: Internal server error
components:
  schemas:
    SongDetail:
      required:
        - releaseDate
        - text
        - link
      type: object
      properties:
        releaseDate:
          type: string
          example: 16.07.2006
        text:
          type: string
          example: |
            Ooh baby, don't you know I suffer?
            Ooh baby, can you hear me moan?
            You caught me under false pretenses
            How long before you let me go?

            Ooh
            You set my soul alight
            Ooh
            You set my soul alight
        link:
          type: string
          example: https://www.youtube.com/watch?v=Xsp3_a-PMTw
```

---

## 📈 Технические требования:

- **PostgreSQL** для хранения данных (структура БД создаётся миграциями при старте)
- Конфигурация вынесена в `.env` файл
- Коды покрыты логами уровней `debug` и `info`
- Сгенерирован Swagger для собственного API

---

## 🚀 Инструкция по запуску

### 1. Клонировать репозиторий

```bash
git clone git@github.com:Zorynix/song-library.git
cd song-library
```

### 2. Настроить конфигурацию

Создайте файл `.env` в корне проекта и заполните его своими параметрами, по примеру из `.env.example`. Укажите нужные параметры в `config.yaml`

### 3. Запустить сервис

```bash
docker-compose up -d
```

Миграции будут применены автоматически.

---

## 📄 Swagger-документация

После запуска сервиса Swagger UI доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

---

## 🛠️ Стек

- **Go**
- **Chi**
- **Zerolog**
- **PostgreSQL**
- **Docker / Docker Compose**
- **Swagger (OpenAPI)**
- **Prometheus**
- **Grafana**


