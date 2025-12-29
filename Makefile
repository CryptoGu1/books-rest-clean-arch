.PHONY: run stop build

# Запуск всего проекта (с пересборкой)
run:
	docker-compose up --build

# Запуск в фоновом режиме
run-d:
	docker-compose up -d --build

# Остановка контейнеров
stop:
	docker-compose down

# Очистка (удалит контейнеры и тома базы данных - ОСТОРОЖНО)
clean:
	docker-compose down -v --remove-orphans