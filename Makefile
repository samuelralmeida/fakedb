build:
	docker build -t dumpfakedb .

run:
	docker run -d --name dumpfakedb --env-file .env dumpfakedb

bash:
	docker exec -it dumpfakedb bash

remove:
	docker container rm dumpfakedb -f

fakedb-all:
	docker exec -it dumpfakedb ./fakedb

fakedb-all-dev:
	docker exec -it dumpfakedb ./fakedb --target=dev

fakedb-trem-staging:
	docker exec -it dumpfakedb ./fakedb --source=trem --target=staging