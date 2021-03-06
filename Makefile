GOAPTH=$PWD/.dependencies

default: run

test:
	go test ./...

docker:
	docker build --tag registry.wawan.pro/yabloko/alkobot .
	docker push registry.wawan.pro/yabloko/alkobot

run: test
	go run main.go

update:
	docker pull registry.wawan.pro/yabloko/alkobot
	docker run --env-file=.env --net=host -d --restart=always --name=alkobot registry.wawan.pro/yabloko/alkobot

stop:
	docker stop alkobot && docker rm alkobot