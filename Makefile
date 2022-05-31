run: cmd/run_server/run_server
	./cmd/run_server/run_server

cmd/run_server/run_server: **/*.go
	cd cmd/run_server && go build

README.pdf: README.md
	pandoc -V geometry:margin=1in -o README.pdf README.md
