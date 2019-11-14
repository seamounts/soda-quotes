build:
	docker build -t hsxue/soda-quotes:latest .
	docker push hsxue/soda-quotes:latest
	
PHONY: build