#!/bin/bash

# Имя файла .env
ENV_FILE=".env"

# Проверка, существует ли файл .env
if [ -f "$ENV_FILE" ]; then
    echo "$ENV_FILE уже существует. Перезаписать его? (y/n)"
    read answer
    if [ "$answer" != "${answer#[Yy]}" ] ;then
        echo "Перезаписываем $ENV_FILE..."
    else
        echo "Отмена перезаписи. Скрипт завершен."
        exit 0
    fi
fi

# Запрос ввода от пользователя
read -p "Введите порт для gRPC (по умолчанию 6969): " GRPC_PORT
GRPC_PORT=${GRPC_PORT:-6969}

read -p "Введите таймаут для gRPC (по умолчанию 5s): " GRPC_TIMEOUT
GRPC_TIMEOUT=${GRPC_TIMEOUT:-5s}

read -p "Введите хост для PostgreSQL (по умолчанию localhost): " POSTGRES_HOST
POSTGRES_HOST=${POSTGRES_HOST:-localhost}

read -p "Введите порт для PostgreSQL (по умолчанию 5435): " POSTGRES_PORT
POSTGRES_PORT=${POSTGRES_PORT:-5435}

read -p "Введите имя пользователя для PostgreSQL (по умолчанию example_user): " POSTGRES_USER
POSTGRES_USER=${POSTGRES_USER:-example_user}

read -p "Введите пароль для PostgreSQL (по умолчанию example_password): " POSTGRES_PASSWORD
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-example_password}

read -p "Введите имя базы данных для PostgreSQL (по умолчанию example_db): " POSTGRES_DBNAME
POSTGRES_DBNAME=${POSTGRES_DBNAME:-example_db}

read -p "Введите режим SSL для PostgreSQL (по умолчанию disable): " POSTGRES_SSLMODE
POSTGRES_SSLMODE=${POSTGRES_SSLMODE:-disable}

read -p "Введите путь к миграциям (по умолчанию ./migrations): " MIGRATIONS_PATH
MIGRATIONS_PATH=${MIGRATIONS_PATH:-./migrations}

# Заполнение файла .env
cat <<EOL > $ENV_FILE
# gRPC Configuration
GRPC_PORT=$GRPC_PORT
GRPC_TIMEOUT=$GRPC_TIMEOUT

# PostgreSQL Configuration
POSTGRES_HOST=$POSTGRES_HOST
POSTGRES_PORT=$POSTGRES_PORT
POSTGRES_USER=$POSTGRES_USER
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
POSTGRES_DBNAME=$POSTGRES_DBNAME
POSTGRES_SSLMODE=$POSTGRES_SSLMODE

# Migrations Configuration
MIGRATIONS_PATH=$MIGRATIONS_PATH
EOL

echo "$ENV_FILE успешно создан и заполнен."
