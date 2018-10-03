.PHONY: run build-upload

run:
	GOTALK_HOST=0.0.0.0 GOTALK_PORT=4040 GOTALK_USE_SYSTEMD_SOCKET=false go run main.go

build-upload:
	GOOS=linux GOARCH=amd64 go build -v .
	rsync -v --progress ./talktask root@gotalk.waltton.com.br:/root/talktask
	rm ./talktask
