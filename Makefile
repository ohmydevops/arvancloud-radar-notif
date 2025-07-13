OUTPUT_DIR = build
SOURCE = ./cmd/radar-notif/

all: linux windows mac

linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/radar-linux $(SOURCE)
	cp icon.png $(OUTPUT_DIR)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/radar-windows.exe $(SOURCE)
	cp icon.png $(OUTPUT_DIR)

mac:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/radar-mac $(SOURCE)
	cp icon.png $(OUTPUT_DIR)

clean:
	rm -rf $(OUTPUT_DIR)
