# Название образа и контейнера
IMAGE_NAME=my_forum_image
CONTAINER_NAME=my_forum_container

# Правило по умолчанию
.PHONY: all
all: build run

# Правило для сборки Docker-образа
.PHONY: build
build:
	docker build -t $(IMAGE_NAME) .

# Правило для удаления старого контейнера
.PHONY: clean
clean:
	@if [ "$(shell docker ps -aq -f name=$(CONTAINER_NAME))" ]; then \
		echo "Removing old container..."; \
		docker rm -f $(CONTAINER_NAME); \
	fi

# Правило для запуска нового контейнера
.PHONY: run
run: clean
	docker run -d --name $(CONTAINER_NAME) -p 8080:8080 $(IMAGE_NAME)

# Правило для остановки контейнера
.PHONY: stop
stop:
	@if [ "$(shell docker ps -aq -f name=$(CONTAINER_NAME))" ]; then \
		echo "Stopping container..."; \
		docker stop $(CONTAINER_NAME); \
		docker rm $(CONTAINER_NAME); \
	fi

# Правило для повторного запуска (стоп + запуск)
.PHONY: restart
restart: stop run

# Правило для проверки состояния контейнера
.PHONY: status
status:
	docker ps -a -f name=$(CONTAINER_NAME)
