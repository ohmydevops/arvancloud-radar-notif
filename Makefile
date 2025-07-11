OUTPUT_DIR = build
SOURCE = main.go

all: linux windows mac

linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/radar-linux $(SOURCE)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/radar-windows.exe $(SOURCE)

mac:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/radar-mac $(SOURCE)

clean:
	rm -rf $(OUTPUT_DIR)
