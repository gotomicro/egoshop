ROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
MUSES_SYSTEM:=github.com/goecology/muses/pkg/system
PORT?=3306
GITHUB:=$(GOPATH)/src/github.com/goecology/egoshop
ADMINOUT:=$(GITHUB)/appgo/apps/shopadmin/bin/shopadmin
APIOUT:=$(GITHUB)/appgo/apps/shopapi/bin/shopapi

APPS = shopapi shopadmin

# 同步自己开发环境mysql表结构
local.createdb:
	cd $(GOPATH)/src/github.com/goecology/generator && make build
	@cd $(GITHUB)/tool/createdb && go build && $(GITHUB)/tool/createdb/createdb -m "root:root@tcp(localhost:3306)"

local.mockdb:
	@cd $(GITHUB)/tool/mockdb && go build && $(GITHUB)/tool/mockdb/mockdb -m "root:root@tcp(localhost:$(PORT))" --endpoints="" --key="" --secret="" --bucket=""


# 执行wechat
wechat:
	@cd $(GITHUB)/appuni && npm run dev:mp-weixin

# 执行go指令
go.api:
	@cd $(GITHUB)/appgo/apps/shopapi && $(GITHUB)/tool/build.sh $(APIOUT) $(MUSES_SYSTEM) $(MUSES_SYSTEM) && $(APIOUT) start --local=false --config=conf/conf.toml

go.admin:
	@cd $(GITHUB)/appgo/apps/shopadmin && $(GITHUB)/tool/build.sh $(ADMINOUT) $(MUSES_SYSTEM) $(MUSES_SYSTEM) && $(ADMINOUT) start --local=false --config=conf/conf.toml
