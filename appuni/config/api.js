let host = ""

import ConfigDev from './config-dev.js';
import ConfigProd from './config-prod.js';


// uEnvDev
if (process.env.NODE_ENV === 'development') {
  host = ConfigDev.host
}
// uEnvProd
if (process.env.NODE_ENV === 'production') {
  host = ConfigProd.host
}



const api = {
  home: {
    index: {
      url: host+`/api/v1/home/index`,
      method: 'GET',
    }
  },
  cate: {
    index: {
      url: host+`/api/v1/cate/index`,
      method: 'GET',
    }
  },
  //商城
  com: {
    index: {
      url:host+`/api/v1/com/index`,
      method: 'GET',
    },
    list: {
      url: host+`/api/v1/com/list`,
      method: 'GET',
    },
    info: {
      url: host+`/api/v1/com/info/:id`,
      method: 'GET',
    },
  },


  feed: {
    list: {
      url: host+`/api/v1/feed/list`,
      method: 'GET',
    },
    info: {
      url: host+`/api/v1/feed/info/:feedId`,
      method: 'GET',
    },
  },
  comment: {
    list: {
      url: host+`/api/v1/comment/list`,
      method: 'GET',
    },
    listTop3: {
      url: host+`/api/v1/comment/listTop3`,
      method: 'GET',
    },
    create: {
      url: host+`/api/v1/comment/create`,
      method: 'POST',
    },
  },

  goods: {
    // 查询免费或非免费（即所有）图片的单个goods。非免费（即所有）图片的goods需要鉴权
    info: {
      url: host+`/api/v1/goods/:gid`,
      method: 'GET',
    },
    // 查询免费或非免费（即所有）图片的所有goods。非免费（即所有）图片的goods需要鉴权
    list: {
      url: host+`/api/v1/goods`,
      method: 'GET',
    },
    // 获取商品评论
    getComment: {
      url: host+`/api/v1/goods/:gid/comments`,
      method: 'GET',
    },
    // 提交商品评论
    postComment: {
      url: host+`/api/v1/goods/:gid/comments`,
      method: 'POST',
    },
    // 修改商品评论
    putComment: {
      url: host+`/api/v1/goods/:gid/comments/:cmtid`,
      method: 'PUT',
    },
    // 删除商品评论
    deleteComment: {
      url: host+`/api/v1/goods/:gid/comments/:cmtid`,
      method: 'DELETE',
    },
    // 获取下载链接
    genDownload: {
      url: host+`/api/v1/goods/:gid/download_url`,
      method: 'GET',
    },
    // 上传资源
  },
  userGoods: {
    // 查询用户“关注”、“分享”、“上传”、“购买”的所有商品
    list: {
      url: host+`/api/v1/users/goods/list`,
      method: 'GET',
    },
    // 查询用户单个商品
    stats: {
      url: host+`/api/v1/users/goods/stats/:gid/:typeId`,
      method: 'GET',
    },
    // 关注商品
    star: {
      url: host+`/api/v1/users/star`,
      method: 'POST',
    },
    // 取消关注商品
    unstar: {
      url: host+`/api/v1/users/unstar`,
      method: 'POST',
    },
    // 关注商品
    collect: {
      url: host+`/api/v1/users/collect`,
      method: 'POST',
    },
    // 取消关注商品
    uncollect: {
      url: host+`/api/v1/users/uncollect`,
      method: 'POST',
    },
    // 分享商品
    share: {
      url: host+`/api/v1/users/share`,
      method: 'POST',
    },
    // 购买商品
    pay: {
      url: host+`/api/v1/users/goods/:gid/pay`,
      method: 'POST',
    },
    repay: {
      url: host+`/api/v1/users/goods/:gid/repay`,
      method: 'POST',
    },
    // 上传商品
    upload: {
      url: host+`/api/v1/users/goods`,
      method: 'POST',
    }
  },
  res: {
    post: {
      url: host+`/api/v1/res`,
      method: 'POST',
    }
  },
  topic: {
    index: {
      url: `/topic/index`,
      method: 'GET',
    }
  },

  users: {
    stats: {
      url: host+`/api/v1/users/stats`,
      method: 'GET',
    },
    signin: {
      url: host+`/api/v1/users/signin`,
      method: 'POST',
    }
  },
  auth: {
    login: {
      url: host+`/api/v1/auth/login`,
      method: 'POST',
    }
  },

  // 订单
  order : {
    create: {
      url: host+`/api/v1/order/create`,
      method: 'POST',
    },
    delete: {
      url: host+`/api/v1/order/delete`,
      method: 'POST',
    },
    list: {
      url: host+`/api/v1/order/list`,
      method: 'GET',
    },
    info: {
      url: host+`/api/v1/order/info`,
      method: 'GET',
    },
    pay: {
      url: host+`/api/v1/order/pay`,
      method: 'POST',
    }
  },

  // 购物车
  cart: {
    create: {
      url: host+`/api/v1/cart/create`,
      method: 'POST',
    },
    list: {
      url: host+`/api/v1/cart/list`,
      method: 'POST',
    },
    totalNum: {
      url: host+`/api/v1/cart/totalNum`,
      method: 'GET',
    },
  },
  // 购买商品
  pay: {
    calculate: {
      url: host+`/api/v1/pay/calculate`,
      method: 'POST',
    },
  },
  address: {
    list: {
      url: host+`/api/v1/address/list`,
      method: 'GET',
    },
    typeList: {
      url: host+`/api/v1/address/typeList`,
      method: 'GET',
    },
    default: {
      url: host+`/api/v1/address/default`,
      method: 'GET',
    },
    info: {
      url: host+`/api/v1/address/info`,
      method: 'GET',
    },
    create: {
      url: host+`/api/v1/address/create`,
      method: 'POST',
    },
    update: {
      url: host+`/api/v1/address/update`,
      method: 'POST',
    },
    delete: {
      url: host+`/api/v1/address/delete`,
      method: 'POST',
    },
  }
}

export {api}
