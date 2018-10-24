.PHONY: run build-upload

run:
	GOTALK_HOST=0.0.0.0 GOTALK_PORT=4040 GOTALK_USE_SYSTEMD_SOCKET=false GOTALK_ACD_POOL_SIZE=8 GOTALK_ACD_QUEUEL_SIZE=32 go run main.go

build-upload:
	GOOS=linux GOARCH=amd64 go build -v .
	rsync -v --progress ./talktask root@gotalk-demo.waltton.com.br:/root/talktask
	rm ./talktask
