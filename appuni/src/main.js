import Vue from 'vue'
import App from './App'
import store from './store'
import index from "./utils";
import {api} from "../config/api"
var qcloud = require("wafer2-client-sdk/index.js");

Vue.config.productionTip = false

const tui = {
  toast: function(text, duration, success) {
    uni.showToast({
      title: text,
      icon: success ? 'success' : 'none',
      duration: duration || 2000
    })
  },
  constNum: function() {
    const res = uni.getSystemInfoSync();
    return res.platform.toLocaleLowerCase() == "android" ? 300 : 0;
  },
  px: function(num) {
    return uni.upx2px(num) + 'px';
  },
  interfaceUrl: function() {
    //接口地址
    return "https://www.thorui.cn";
  },
  request: function(url, postData, method, type, hideLoading) {
    //接口请求
    if (!hideLoading) {
      uni.showLoading({
        mask: true,
        title: '请稍候...'
      })
    }
    return new Promise((resolve, reject) => {
      uni.request({
        url: this.interfaceUrl() + url,
        data: postData,
        header: {
          'content-type': type ? 'application/x-www-form-urlencoded' : 'application/json',
          'authorization': this.getToken(),
          'security': 1
        },
        method: method, //'GET','POST'
        dataType: 'json',
        success: (res) => {
          !hideLoading && uni.hideLoading()
          resolve(res.data)
        },
        fail: (res) => {
          if (!hideLoading) {
            this.toast("网络不给力，请稍后再试~")
          }
          reject(res)
        }
      })
    })
  },
  uploadFile: function(src) {
    const that = this
    uni.showLoading({
      title: '请稍候...'
    })
    return new Promise((resolve, reject) => {
      const uploadTask = uni.uploadFile({
        url: 'https://abc.cc',
        filePath: src,
        name: 'file',
        header: {
          'content-type': 'multipart/form-data'
        },
        formData: {},
        success: function(res) {
          uni.hideLoading()
          let d = JSON.parse(res.data)
          if (d.code === 1) {
            let fileObj = JSON.parse(d.data)[0];
            //文件上传成功后把图片路径数据提交到服务器，数据提交成功后，再进行下张图片的上传
            resolve(fileObj)
          } else {
            that.toast(res.message);
          }
        },
        fail: function(res) {
          reject(res)
          uni.hideLoading();
          that.toast(res.message);
        }
      })
    })
  },
  setToken: function(token) {
    uni.setStorageSync("token", token)
  },
  getToken() {
    return uni.getStorageSync("token")
  },
  getUserInfo() {
    return uni.getStorageSync("userInfo")
  },
  isLogin: function() {
    const token =  uni.getStorageSync("token")
    const userInfo = uni.getStorageSync('userInfo');
    if (!token || !userInfo) {
      return false
    }
    return true
  },
  webURL:function(){
    return "https://www.thorui.cn/wx"
  },
  checkLogin: function () {
    if (this.isLogin()) {
      return true
    }
    console.log("fffffffff")
    qcloud.setLoginUrl(api.auth.login.url);
    var that  = this
    uni.getSetting({
      success (res){
        if (res.authSetting['scope.userInfo']) {
          // 已经授权，
          // uni.login({
          //   success: result => {
          //     uni.getUserInfo({
          //       code: result.code,
          //       withCredentials: true,
          //       success: function (res) {
          //         const userInfo = res.userInfo;
          //         request('/Api/User/user_info', {
          //           code: result.code,
          //           encryptedData: res.encryptedData,
          //           iv: res.iv
          //         }).then(res => {
          //           uni.setStorageSync('login', res.data);
          //           _t.setData({
          //             canIUse: true
          //           })
          //           _t.userMy()
          //         })
          //       }
          //     })
          //   }
          // })
          qcloud.login({
            success: res => {
              uni.hideLoading();
              uni.setStorageSync("userInfo", res);
              var openId = res.openId;
              uni.setStorageSync("openId", openId);
              that.setToken(qcloud.Session.get().skey)
              that.toast("登录成功",1500,"success")
            },
            fail: err => {
              console.log(err);
              uni.hideLoading();
            }
          });
        } else {
          uni.showToast({
            title: '您之前拒绝过授权，需要你点击右上角重新授权',
            icon: 'none',
            duration: 2000
          })
        }
      },
      fail(err) {
        console.log(err)
        reject({
          msg: '网络错误'
        })
      }
    })
    // uni.showModal({
    //   title: '登录',
    //   content: '您尚未登录，请先登录',
    //   showCancel: false,
    //   confirmColor: "#5677FC",
    //   confirmText: '确定',
    //   success: function success(e) {
    //     uni.redirectTo({
    //       url: '/pages/login/main?url='+url });
    //
    //   } });
    return false
  },
  /*获取当前路由*/
  getCurPage: function (){
    let pages = getCurrentPages();
    return pages[pages.length - 1].route
  },
}
// 将微信request方法赋值到vue上
Vue.prototype.wxhttp=index
Vue.prototype.tui = tui
Vue.prototype.$eventHub = Vue.prototype.$eventHub || new Vue()
Vue.prototype.$store = store
App.mpType = 'app'

const app = new Vue({
  ...App
})
app.$mount()
