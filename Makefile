export tag=v1.0
root:
	export ROOT=github.com/myapp/

build:
	echo "building myapp binary"
	mkdir -p bin/amd64
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/amd64/myapp ./cmd/server/main.go

release: build
	echo "building httpserver container"
	docker ps -a -q --filter ancestor=myapp:${tag} | xargs -r docker rm -f
	docker images -a -q --filter "reference=myapp:${tag}" | xargs -r docker rmi -f
	docker build -t myapp:${tag} .

push: release
	echo "pushing myapp"
	docker push myapp:${tag}

run: release
	echo "run myapp"
	docker run -v /home/hyk/HYK/myapp/db/:/yk/db/ -v /home/hyk/HYK/myapp/uploads/:/yk/uploads/  -p 8080:8080 --name myapp -d myapp:${tag}

