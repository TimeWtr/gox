# 安装依赖
.PHONY: tidy
tidy:
	@go mod tidy

# go文件格式化
.PHONY: fmt
fmt:
	@sh ./.scripts/fmt.sh

# 项目检查命令
.PHONY: check
check:
	@$(MAKE) --no-print-directory tidy
	@$(MAKE) --no-print-directory fmt
