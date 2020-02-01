import querystring  from 'querystring'

function formatNumber(n) {
  const str = n.toString()
  return str[1] ? str : `0${str}`
}

export function formatTime(date) {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()

  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()

  const t1 = [year, month, day].map(formatNumber).join('/')
  const t2 = [hour, minute, second].map(formatNumber).join(':')

  return `${t1} ${t2}`
}


//-------------------------------------------------------------------------请求的封装

const isMock = false

var qcloud = require("wafer2-client-sdk/index.js");

//请求封装
function innerRequest(url, method, data, header = {}) {
  uni.showLoading({
    title: '加载中' //数据请求前loading
  })
  const session = qcloud.Session.get();
  if (session) {
    header = {
      'content-type': 'application/json', // 默认值
      "Access-Token":session.skey,
    }
  }else {
    header = {
      'content-type': 'application/json', // 默认值
    }
  }
  return new Promise((resolve, reject) => {
    if (method == "DOWNLOAD") {
      uni.getSetting({
        success: (response) => {
          console.log(response)
          if (!response.authSetting['scope.writePhotosAlbum']) {
            uni.authorize({
              scope: 'scope.writePhotosAlbum',
              success: () => {
                console.log('yes')
              },
              fail(err) {
                console.log("err",err)
              }
            })
          }
        }
      })
      const downloadTask = uni.downloadFile({
        url: url, //仅为示例，并非真实的接口地址
        method: "GET",
        header: header,
        success (res) {
          // 代码需要授权
          uni.saveImageToPhotosAlbum({
            filePath: res.tempFilePath,
            success(res) {
              uni.showToast({
                title: '保存成功',
                icon: 'success',
                duration: 2000
              })
            },
            fail: function (err) {
              console.log(err);
              if (err.errMsg === "saveImageToPhotosAlbum:fail auth deny") {
                console.log("当初用户拒绝，再次发起授权")
                uni.openSetting({
                  success(settingdata) {
                    console.log(settingdata)
                    if (settingdata.authSetting['scope.writePhotosAlbum']) {
                      console.log('获取权限成功，给出再次点击图片保存到相册的提示。')
                    } else {
                      console.log('获取权限失败，给出不给权限就无法正常使用的提示')
                    }
                  },fail: function (err) {
                    console.log(err);
                  }
                })
              }
            },
            complete(res) {
              console.log(res);
            }
          })
        }
      })

          // console.log("Res",res)
          // uni.playVoice({
          //   filePath: res.tempFilePath
          // })
      //   }
      // })
      //
      // downloadTask.onProgressUpdate((res) => {
      //   console.log('下载进度', res.progress)
      //   console.log('已经下载的数据长度', res.totalBytesWritten)
      //   console.log('预期需要下载的数据总长度', res.totalBytesExpectedToWrite)
      // })
      //
      // downloadTask.abort() // 取消下载任务
    }else {
      uni.request({
        url: url,
        method: method,
        data: data,
        header: header,
        success: function (res) {
          uni.hideLoading();
          // 如果服务端返回401状态码，跳转到登录界面
          if (res.data.code === 401) {
            qcloud.Session.clear();
            uni.removeStorageSync("token");
            uni.showToast({
              title: "登录异常，清除缓存，请再次重新登录",
              icon:  'none',
              duration: 2000
            })
            return
            //鉴权不通过，清除对应session
            // uni.navigateTo({
            //   url: "/pages/login/main"
            // });
          }

          resolve(res.data)
        },
        fail: function (error) {
          uni.hideLoading();
          reject(false)
        },
        complete: function () {
          uni.hideLoading();
        }
      })
    }

  })
}

export default {
  // api:请求path; data:请求payload; queris:path后queries参数（即"?"后参数）
  route(api, data, queries) {
    if (isMock) {
      return api.mock.data
    }else {
      var url = api.url
      // 如果url里包含: ，说明是特殊的url，类似于/goods/:id
      if (url.indexOf(":") != -1 && data != undefined) {
        // 遍历queries，拿到需要放到url的"?"后的参数
        let queryStr = querystring.stringify(queries)
        if (queryStr !== "") {
          url += "?"+queryStr
        }
        // 遍历数据，拿到data里的key和value
        Object.keys(data).forEach(function(key) {
          if (url.indexOf("/:"+key) != -1) {
                url = url.replace("/:"+key,"/"+data[key])
                delete(data[key])
              }
            }
        )
        // if (api.method == "GET") {
        //   data = {}
        // }


      }

      return innerRequest(url, api.method, data)
    }
  }
}

//-------------------------------------------------------------------------请求的封装


//----------------------------------------------用户是否登录 未登录跳转到登录页面 -------------------------

//
// export function toLogin() {
//   const userInfo = uni.getStorageSync('userInfo');
//   if (!userInfo) {
//     uni.navigateTo({
//       url: "/pages/login/main"
//     });
//   } else {
//     return true;
//   }
// }
//
// export function login() {
//   const userInfo = uni.getStorageSync('userInfo');
//   if (userInfo) {
//     return userInfo;
//   }
// }
//
// //----------------------------------------------用户是否登录 未登录跳转到登录页面 -------------------------
//
//
// export function getStorageOpenid() {
//   const openId = uni.getStorageSync("openId");
//   if (openId) {
//     return openId;
//   } else {
//     return ''
//   }
// }
//
//
//
//
// export function getOpenid() {
//   // uni.login({
//   //   success: res => {
//   //     if (res.code) {
//   //       //发起网络请求
//   //       uni.request({
//   //         url: 'https://api.weixin.qq.com/sns/jscode2session',
//   //         data: {
//   //           "appid": "wx601ce71bde7b9add",
//   //           "secret": "abed5421d88eb8236e933c6e42d5c14e",
//   //           "js_code": res.code,
//   //           "grant_type": "authorization_code"
//   //         },
//   //         success: function (data) {
//   //           var openid = data.data.openid;
//   //           uni.setStorageSync("openid", openid);
//   //         }
//   //       })
//   //     } else {
//   //       console.log('登录失败！' + res.errMsg)
//   //     }
//
//   //   },
//   //   fail: () => {},
//   //   complete: () => {}
//   // });
// }
