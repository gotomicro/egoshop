ROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
MUSES_SYSTEM:=github.com/goecology/muses/pkg/system
PORT?=3306
ADMINNAME:=egoshopadmin
APINAME:=egoshopapi
APPPATH:=$(GOPATH)/src/github.com/goecology/egoshop
ADMINOUT:=$(APPPATH)/appgo/apps/shopadmin/bin/shopadmin
APIOUT:=$(APPPATH)/appgo/apps/shopapi/bin/shopapi

APPS = shopapi shopadmin

# 同步自己开发环境mysql表结构
local.createdb:
	@cd $(APPPATH)/tool/createdb && go build && $(APPPATH)/tool/createdb/createdb -m "root:root@tcp(localhost:3306)"

local.mockdb:
	@cd $(APPPATH)/tool/mockdb && go build && $(APPPATH)/tool/mockdb/mockdb -m "root:root@tcp(localhost:$(PORT))" --endpoints="" --key="" --secret="" --bucket=""


# 执行wechat
wechat:
	@cd $(APPPATH)/appuni && npm run dev:mp-weixin

# 执行go指令
go.api:
	@cd $(APPPATH)/appgo/apps/shopapi && $(APPPATH)/tool/build.sh $(APINAME) $(APIOUT) $(MUSES_SYSTEM) && $(APIOUT) start --conf=conf/conf.toml

go.admin:
	@cd $(APPPATH)/appgo/apps/shopadmin && $(APPPATH)/tool/build.sh $(ADMINNAME) $(ADMINOUT) $(MUSES_SYSTEM) && $(ADMINOUT) start --conf=conf/conf.toml
