# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
# make命令如果没有指定目标，默认会执行第一个目标，即help。所以help要放在第一个。
# 给规则添加伪目标(.PHONY: help)，可以防止目标和项目中的真实文件名冲突。
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


# Create the new confirm target.
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# BUILD
# ==================================================================================== #

# 在目标名称上加命名空间，推荐 xxx/xxx 的方式，可以更好的组织大型项目
current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

# 通过指示 Go 的链接器去(-ldflags="-s")除二进制文件中的 DWARF 调试信息和符号表，可以将二进制文件的大小大约减少 25%。
## build/cmd: build the application
.PHONY: build/cmd
build/cmd:
	@echo 'Building cmd...'
	go build -o=./inspector/wps-mcp-server.exe -ldflags=${linker_flags} ./cmd
