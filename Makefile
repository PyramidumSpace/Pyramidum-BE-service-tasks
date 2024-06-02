.PHONY: envfile

envfile:
	@echo "Creating .env file from .env.example"
	cp .env.example .env
