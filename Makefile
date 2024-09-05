build:
	@go build -o bin/ConsultationsReminder.exe

run: build
	@./bin/ConsultationsReminder.exe
