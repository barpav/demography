up:
	sudo docker-compose up -d
down:
	sudo docker-compose down

up-debug:
	sudo docker-compose -f debug.compose.yaml up -d
down-debug:
	sudo docker-compose -f debug.compose.yaml down

# make person N=Lev P=Nikovaevich S=Tolstoy
person:
	curl -v -X POST	-H "Content-Type: application/vnd.newPersonData.v1+json" \
	-d '{"surname": "$(S)", "name": "$(N)", "patronymic": "$(P)"}' \
	localhost:8080/v1/people