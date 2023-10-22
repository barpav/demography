up:
	sudo docker-compose up -d
down:
	sudo docker-compose down

up-debug:
	sudo docker-compose -f debug.compose.yaml up -d
down-debug:
	sudo docker-compose -f debug.compose.yaml down