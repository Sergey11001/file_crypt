# Базовый образ для сборки React-приложения
FROM node:20 as build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем весь проект в контейнер
COPY . .

# Собираем приложение
RUN npm run build

# Базовый образ для сервера Nginx
FROM nginx:alpine

# Копируем собранное приложение из предыдущего шага в папку для статического контента
COPY --from=build /app/build /usr/share/nginx/html

# Добавляем кастомный конфиг для Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Экспонируем порт для фронтенда
EXPOSE 80

# Запуск Nginx
CMD ["nginx", "-g", "daemon off;"]
