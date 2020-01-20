build: cmd/server/main.go
	CGO_ENABLED=0 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w -s" cmd/server/main.go
run:
	go run cmd/server/main.go
deploy:
	sudo docker build -t registry.heroku.com/choose-your-own-adventure-go/web -f build/Dockerfile .
	sudo docker push registry.heroku.com/choose-your-own-adventure-go/web 
	heroku container:release web -a choose-your-own-adventure-go