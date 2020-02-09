MUSES_SYSTEM:=github.com/goecology/muses/pkg/system
APPNAME:=egoshop
APPPATH:=$(GOPATH)/src/github.com/goecology/egoshop
APPOUT:=$(APPPATH)/build/$(APPNAME)

# 执行wechat
wechat:
	@cd $(APPPATH)/appuni && npm run dev:mp-weixin


ant:
	@cd $(APPPATH)/adminant && npm start

# 执行go指令
go:
	@cd $(APPPATH)/appgo && $(APPPATH)/tool/build.sh $(APPNAME) $(APPOUT) $(MUSES_SYSTEM)
	@cp $(APPPATH)/appgo/conf $(APPPATH)/build/ -R
	@cd $(APPPATH)/build && $(APPOUT) start --conf=conf/conf.toml


install:
	@cd $(APPPATH)/appgo && $(APPPATH)/tool/build.sh $(APPNAME) $(APPOUT) $(MUSES_SYSTEM)
	@cp $(APPPATH)/appgo/conf $(APPPATH)/build/ -R
	@cd $(APPPATH)/build && $(APPOUT) install --conf=conf/conf.toml

all:
	@rm -r $(APPPATH)/build
	@cd $(APPPATH)/appgo && $(APPPATH)/tool/build.sh $(APPNAME) $(APPOUT) $(MUSES_SYSTEM)
	@cp $(APPPATH)/appgo/conf $(APPPATH)/build/ -R
	@cd $(APPPATH)/adminant && npm run build
	@cp $(APPPATH)/adminant/dist $(APPPATH)/build/ -R
	@mv $(APPPATH)/build/dist $(APPPATH)/build/adminant
	@cd $(APPPATH)/appuni && npm run build:mp-weixin
	@cp $(APPPATH)/appuni/dist/build/mp-weixin $(APPPATH)/build/ -R
	@tar -zxvf build.tar.gz build

