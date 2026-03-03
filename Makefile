ANDROID_SDK_ROOT := C:/Users/timof/AppData/Local/Android/Sdk
EMULATOR := $(ANDROID_SDK_ROOT)/emulator/emulator
ADB := $(ANDROID_SDK_ROOT)/platform-tools/adb


AVD_NAME := Pixel_9_PRO_API_35

generate:
	protoc -I ./backend/proto/gen \
    		--proto_path=./backend/proto \
    		--go_out ./backend/proto/gen --go_opt paths=source_relative \
    		--go-grpc_out ./backend/proto/gen --go-grpc_opt paths=source_relative \
    		--grpc-gateway_out ./backend/proto/gen --grpc-gateway_opt paths=source_relative \
    		--openapiv2_out ./backend/proto/gen --openapiv2_opt logtostderr=true \
    		./backend/proto/auth/auth.proto \
    		./backend/proto/expenses/expenses.proto \

lint:
	cd backend && golangci-lint run --fix

build:
	docker-compose build

start:
	docker-compose up

emulator-start:
	@echo "Starting Android emulator: $(AVD_NAME)"
	"$(EMULATOR)" -avd $(AVD_NAME) -netdelay none -netspeed full &

mobile-start:
	cd mobile-app && npm start