.PHONY: all build run clean test

# 设置Go编译器参数
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 项目名称和主程序入口
BINARY_NAME=hertz-server
MAIN_FILE=main.go
BINARY_PATH=bin/$(BINARY_NAME)

# 默认目标
all: build

# 编译项目
build:
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_FILE)

# 运行项目
run:
	$(GORUN) $(MAIN_FILE)

# 清理编译文件
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# 运行测试
test:
	$(GOTEST) -v ./...

# 更新依赖
deps:
	$(GOMOD) tidy

# 生成Hertz代码
hz-gen:
	hz update -I idl -idl idl/short/*.proto

# 安装hz工具
install-hz:
	$(GOGET) github.com/cloudwego/hertz/cmd/hz@latest