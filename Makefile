# Имя файла для Bash скрипта
ENV_SCRIPT = scripts/fill_env.sh

# Цель по умолчанию
.PHONY: all
all: fill-env

# Цель для создания и заполнения .env файла
.PHONY: fill-env
fill-env:
	@echo "Запуск скрипта для заполнения .env файла..."
	./$(ENV_SCRIPT)

# Очистка
.PHONY: clean
clean:
	@rm -f .env
