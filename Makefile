build:
	@go build -o bin/ConsultationsReminder.exe

run: build
	@./bin/ConsultationsReminder.exe

bau: # Build and update
	@go build -o bin/ConsultationsReminder
	pscp -i "%USERPROFILE%/Documents/prin.ppk" -r ./bin root@77.232.42.104:/root/
	
