PROTO_DIR := proto
OUT_DIR := ./proto

# Поиск всех .proto файлов в директории proto
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

# Команда для генерации Go-кода из proto файлов
gen:
	# Генерация Go-кода и gRPC-кода
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) --proto_path=$(PROTO_DIR) $(PROTO_FILES)

# Цель по умолчанию
.PHONY: generate