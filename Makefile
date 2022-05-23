build:
	docker build -t dumpfakedb .

run:
	docker run -d --name dumpfakedb --env-file .env -p 3306:3306 dumpfakedb

bash:
	docker exec -it dumpfakedb bash

remove:
	docker container rm dumpfakedb -f

fakedb:
	docker exec -it dumpfakedb ./fakedb