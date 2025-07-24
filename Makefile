buildDocker:
	docker build -t subscription-service .

upContainer:
	docker compose up -d

upServer: buildServer
	./build/main --config ./config/local.yml

buildServer:
	mkdir -p build
	go build -o ./build/ -v cmd/*.go

clean:
	rm -rf build