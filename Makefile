compile:
	env GOOS=linux GOARCH=amd64 go build -o ./rel/bin/auto-node-repair -i cmd/auto-node-repair/main.go

release: 
	env GOOS=linux GOARCH=amd64 go build -o ./rel/bin/auto-node-repair -i cmd/auto-node-repair/main.go
	docker build -t kvmprashanth/auto-node-repair:v4 .
	docker push kvmprashanth/auto-node-repair:v4
