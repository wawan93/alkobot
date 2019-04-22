default: run

test:
	go test ./...

docker:
	docker rmi -f registry.wawan.pro/wawan/alkobot
	docker build --tag registry.wawan.pro/wawan/alkobot .
	docker push registry.wawan.pro/wawan/alkobot

run: test
	go run main.go

update:
	docker pull registry.wawan.pro/wawan/alkobot
	docker run --env-file=.env --net=host -d --restart=always --name=alkobot registry.wawan.pro/wawan/alkobot

stop:
	docker stop alkobot && docker rm alkobot