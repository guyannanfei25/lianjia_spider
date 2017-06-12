all:
	export GOPATH=$$GOPATH:`pwd`; go build -o bin/send_mail src/main/send_mail.go

clean:
	rm -fv bin/*
