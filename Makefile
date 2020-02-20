MUSES_SYSTEM:=github.com/i2eco/muses/pkg/system
APPNAME:=egoshop
APPPATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
APPOUT:=$(APPPATH)/appgo/$(APPNAME)

# 执行wechat
wechat:
	@cd $(APPPATH)/appuni && npm run dev:mp-weixin


ant:
	@cd $(APPPATH)/adminant && npm start

# 执行go指令
go:
	@cd $(APPPATH)/appgo && $(APPPATH)/tool/build.sh $(APPNAME) $(APPOUT) $(MUSES_SYSTEM)
	@cd $(APPPATH)/appgo && $(APPOUT) start --conf=conf/conf.toml


install:
	@cd $(APPPATH)/appgo && $(APPPATH)/tool/build.sh $(APPNAME) $(APPOUT) $(MUSES_SYSTEM)
	@cd $(APPPATH)/appgo && $(APPOUT) install --conf=conf/conf.toml

install.create:
	@cd $(APPPATH)/appgo && $(APPPATH)/tool/build.sh $(APPNAME) $(APPOUT) $(MUSES_SYSTEM)
	@cd $(APPPATH)/appgo && $(APPOUT) install --conf=conf/conf.toml --mode=create

install.clear:
	@cd $(APPPATH)/appgo && $(APPPATH)/tool/build.sh $(APPNAME) $(APPOUT) $(MUSES_SYSTEM)
	@cd $(APPPATH)/appgo && $(APPOUT) install --conf=conf/conf.toml --clear=true


