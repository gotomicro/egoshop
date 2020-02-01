<template>
	<view class="container">
		<!--header-->
		<view class="tui-header-box" :style="{height:height+'px',background:'rgba(255,255,255,'+opcity+')'}">
			<view class="tui-header" :style="{paddingTop:top+'px', opacity:opcity}">
				宠物详情
			</view>
			<view class="tui-header-icon" :style="{marginTop:top+'px'}">
				<view class="tui-icon tui-icon-arrowleft tui-icon-back" :style="{color:opcity>=1?'#000':'#fff',background:'rgba(0, 0, 0,'+iconOpcity+')'}"
				 @tap="back"></view>
				<view class="tui-icon tui-icon-more-fill  tui-icon-ml" :style="{color:opcity>=1?'#000':'#fff',background:'rgba(0, 0, 0,'+iconOpcity+')',fontSize:'20px'}"
				 @tap.stop="openMenu"></view>
				<tui-badge type="red" size="small">5</tui-badge>
			</view>
		</view>
		<!--header-->

		<!--banner-->
		<view class="tui-banner-swiper">
			<swiper :autoplay="true" :interval="5000" :duration="150" :circular="true" :style="{height:scrollH + 'px'}" @change="bannerChange">
				<block v-for="(item,index) in info.gallery" :key="index">
					<swiper-item :data-index="index" @tap.stop="previewImage">
						<image :src="item" class="tui-slide-image" :style="{height:scrollH+'px'}" />
					</swiper-item>
				</block>
			</swiper>
			<tui-tag type="translucent" shape="circleLeft" size="small">{{bannerIndex+1}}/{{info.gallery.length}}</tui-tag>
		</view>

		<!--banner-->

		<view class="tui-pro-detail">
			<view class="tui-product-title tui-border-radius">
				<view class="tui-pro-pricebox tui-padding">
					<view class="tui-pro-price">
						<view>￥<text class="tui-price">{{info.price}}</text></view>
						<tui-tag size="small" :plain="true" type="high-green" shape="circle">新品</tui-tag>
					</view>
					<view class="tui-collection tui-size" @tap="collecting">
						<view class="tui-icon tui-icon-collection" :class="['tui-icon-'+(collected?'like-fill':'like')]" :style="{color:collected?'#ff201f':'#333',fontSize:'20px'}"></view>
						<view class="tui-scale" :class="[collected?'tui-icon-red':'']">收藏</view>
					</view>
				</view>
				<view class="tui-original-price tui-gray" v-if="info.originPrice !== 0">
					价格
					<text class="tui-line-through">￥{{info.originPrice}}</text>
				</view>
				<view class="tui-pro-titbox">
					<view class="tui-pro-title">{{info.name}}</view>
					<button open-type="share" class="tui-share-btn tui-share-position">
						<tui-tag type="gray" tui-tag-class="tui-tag-share tui-size" shape="circleLeft" size="small">
							<view class="tui-icon tui-icon-partake" style="color:#999;font-size:15px"></view>
							<!-- <tui-icon name="partake" color="#999" size="15"></tui-icon> -->
							<text class="tui-share-text tui-gray">分享</text>
						</tui-tag>
					</button>
				</view>
				<view class="tui-padding">
					<view class="tui-sub-title tui-size tui-gray">此商品将于2019-06-28,10点结束闪购特卖，时尚美饰联合专场</view>
					<view class="tui-sale-info tui-size tui-gray">
						<view>阅读：{{info.cntView}}</view>
						<view>点赞：{{info.cntStar}}</view>
                        <view>收藏：{{info.cntCollect}}</view>
                        <view>评论：{{info.cntComment}}</view>
					</view>
				</view>
			</view>


			<view class="tui-basic-info tui-mtop tui-radius-all">
				<view class="tui-list-cell">
					<view class="tui-bold tui-cell-title">标签</view>
					<view class="tui-tag-coupon-box">
						<tui-tag v-for="(item,index) in info.labelsDesc" :key="index" size="small" type="red" shape="circle" tui-tag-class="tui-tag-coupon">{{item}}</tui-tag>
					</view>
					<tui-icon name="more-fill" :size="20" class="tui-right tui-top40" color="#666"></tui-icon>
				</view>
				<view class="tui-list-cell" @tap="showPopup">
					<view class="tui-bold tui-cell-title">送至</view>
					<view class="tui-addr-box">
						<view class="tui-addr-item">北京朝阳区三环到四环之间</view>
						<view class="tui-addr-item">今日23:59前完成下单，预计6月28日23:30前发货，7月1日24:00前送达</view>
					</view>
					<tui-icon name="more-fill" :size="20" class="tui-right" color="#666"></tui-icon>
				</view>
				<view class="tui-list-cell tui-last">
					<view class="tui-bold tui-cell-title">运费</view>
					<view class="tui-selected-box">在线支付免运费</view>
				</view>
				<view class="tui-guarantee">
					<view class="tui-guarantee-item">
						<tui-icon name="circle-selected" :size="14" color="#999"></tui-icon>
						<text class="tui-pl">可配送海外</text>
					</view>
					<view class="tui-guarantee-item">
						<tui-icon name="circle-selected" :size="14" color="#999"></tui-icon>
						<text class="tui-pl">店铺发货&售后</text>
					</view>
					<view class="tui-guarantee-item">
						<tui-icon name="circle-selected" :size="14" color="#999"></tui-icon>
						<text class="tui-pl">7天无理由退货</text>
					</view>
					<view class="tui-guarantee-item">
						<tui-icon name="circle-selected" :size="14" color="#999"></tui-icon>
						<text class="tui-pl">闪电退款</text>
					</view>
					<view class="tui-guarantee-item">
						<tui-icon name="circle-selected" :size="14" color="#999"></tui-icon>
						<text class="tui-pl">极速审核</text>
					</view>
				</view>
			</view>

			<view class="tui-cmt-box tui-mtop tui-radius-all">
				<view class="tui-list-cell tui-last tui-between">
					<view class="tui-bold tui-cell-title">评价</view>
					<view @tap="common">
						<text class="tui-cmt-all">查看全部</text>
						<view class="tui-icon tui-icon-more-fill" style="color:#ff201f; font-size: 20px;"></view>
						<!-- <tui-icon name="more-fill" size="20" color="#ff201f"></tui-icon> -->
					</view>
				</view>

				<view class="tui-cmt-content tui-padding">
					<view class="tui-cmt-user">
						<image src="../../../static/images/news/avatar_2.jpg" class="tui-acatar"></image>
						<view>z***9</view>
					</view>
					<view class="tui-cmt">物流很快，很适合我的风格❤</view>
					<view class="tui-attr">颜色：叠层钛钢流苏耳环（A74）</view>
				</view>

				<view class="tui-cmt-btn">
					<tui-tag type="black" :plain="true" tui-tag-class="tui-tag-cmt" @tap="common">查看全部评价</tui-tag>
				</view>
			</view>

			<view class="tui-nomore-box">
				<tui-nomore text="宝贝详情" :visible="true" bgcolor="#f7f7f7"></tui-nomore>
			</view>
			<view class="tui-product-img tui-radius-all">
				<image :src="'https://www.thorui.cn/img/detail/'+(index+1)+'.jpg'" v-for="(img,index) in 20" :key="index" mode="widthFix"></image>
			</view>
			<tui-nomore text="已经到最底了" :visible="true" bgcolor="#f7f7f7"></tui-nomore>
			<view class="tui-safearea-bottom"></view>
		</view>

		<!--底部操作栏-->
		<view class="tui-operation">
			<view class="tui-operation-left tui-col-5">
				<button open-type="contact" class="tui-operation-item operation-button" hover-class="opcity" :hover-stay-time="150">
					<tui-icon name="kefu" :size="22" color='#333'></tui-icon>
					<view class="tui-operation-text tui-scale-small">客服</view>
				</button>
				<button open-type="share" class="tui-operation-item operation-button" hover-class="opcity" :hover-stay-time="150">
					<tui-icon name="share" :size="22" color='#333'></tui-icon>
					<view class="tui-operation-text tui-scale-small">分享</view>
				</button>
				<button class='tui-operation-item operation-button' open-type="getUserInfo" @getuserinfo="onCollect" >
					<tui-icon :name="collectFlag===1?'like-fill':'like'" :size="22" color='#333'></tui-icon>
					<view class="tui-operation-text tui-scale-small">收藏</view>
				</button>
				<view class="tui-operation-item" hover-class="opcity" :hover-stay-time="150" @tap="toCart">
					<tui-icon name="cart" :size="22" color='#333'></tui-icon>
					<view class="tui-operation-text tui-scale-small">购物车</view>
					<tui-badge type="danger" size="small">{{inCartNumber}}</tui-badge>
				</view>
			</view>
			<view class="tui-operation-right tui-right-flex tui-col-7 tui-btnbox-4">
				<view class="tui-flex-1">
					<tui-button type="danger" shape="circle" size="mini" @click="showPopup">加入购物车</tui-button>
				</view>
				<view class="tui-flex-1">
					<tui-button type="warning" shape="circle" size="mini" @click="showPopup">立即购买</tui-button>
				</view>
			</view>
		</view>


		<!--底部操作栏--->

		<!--顶部下拉菜单-->
		<tui-top-dropdown tui-top-dropdown="tui-top-dropdown" bgcolor="rgba(76, 76, 76, 0.95)" :show="menuShow" :height="0"
		 @close="closeMenu">
			<view class="tui-menu-box tui-padding tui-ptop">
				<view class="tui-menu-header" :style="{paddingTop:top+'px'}">
					功能直达
				</view>
				<view class="tui-menu-itembox">
					<block v-for="(item,index) in topMenu" :key="index">
						<view class="tui-menu-item" hover-class="tui-opcity" :hover-stay-time="150" @tap="navigateToUrl(item.url)">
							<view class="tui-badge-box">
								<tui-icon :name="item.icon" color="#fff" :size="item.size"></tui-icon>
								<tui-badge type="red" tui-badge-class="tui-menu-badge" size="small" v-if="item.badge">{{item.badge}}</tui-badge>
							</view>
							<view class="tui-menu-text">{{item.text}}</view>
						</view>
					</block>
				</view>
				<view class="tui-icon tui-icon-up" style="color: #fff; font-size: 26px;" @tap.stop="closeMenu"></view>
				<!-- <tui-icon name="up" color="#fff" size="26" class="tui-icon-up" @tap.stop="closeMenu"></tui-icon> -->
			</view>

		</tui-top-dropdown>
		<!---顶部下拉菜单-->

		<!--底部选择层-->
		<tui-bottom-popup :show="popupShow" @close="hidePopup">
			<view class="tui-popup-box">
				<view class="tui-product-box tui-padding">
					<image :src="info.cover" class="tui-popup-img"></image>
					<view class="tui-popup-price">
						<view class="tui-amount tui-bold">￥{{price}}</view>
						<view class="tui-number">编号:{{info.id}}</view>
					</view>
				</view>
				<scroll-view scroll-y class="tui-popup-scroll">
					<view class="tui-scrollview-box"  v-for="(spec,index) in specList" :key="index">
						<view class="tui-bold tui-attr-title">{{spec.name}}</view>
						<view class="tui-attr-box">
							<view class="tui-attr-item" :class="[value.checked===true?'tui-attr-active':'']" v-for="(value, valueIndex) in  spec.valueList" @click="onSpecClick(spec.id,value.id)" :key="valueIndex">
								{{value.name}}
							</view>
						</view>

						<view class="tui-number-box tui-bold tui-attr-title">
							<view class="tui-attr-title">数量</view>
							<tui-numberbox :max="999" :min="1" :value="buyNum" @change="changeNum"></tui-numberbox>
						</view>

						<view class="tui-bold tui-attr-title">
							保障服务
						</view>
						<view class="tui-attr-box">
							<view class="tui-attr-item">
								半年掉钻保 ￥4.0
							</view>
						</view>

						<view class="tui-bold tui-attr-title">
							只换不修
						</view>
						<view class="tui-attr-box">
							<view class="tui-attr-item">
								三月意外换￥2.0
							</view>
							<view class="tui-attr-item">
								半年意外换￥2.0
							</view>
						</view>
					</view>
				</scroll-view>
				<view class="tui-operation tui-operation-right tui-right-flex tui-popup-btn">
					<view class="tui-flex-1">
						<tui-button type="red" shape="circle" size="mini" open-type="getUserInfo" @getuserinfo="addCart">加入购物车</tui-button>
					</view>
					<view class="tui-flex-1">
						<tui-button type="warning" shape="circle" size="mini" open-type="getUserInfo" @getuserinfo="submit">立即购买</tui-button>
					</view>
				</view>
				<view class="tui-icon tui-icon-close-fill tui-icon-close" style="color: #999;font-size:20px" @tap="hidePopup"></view>
				<!-- <tui-icon name="close-fill" color="#999" class="tui-icon-close" size="20" @tap="hidePopup"></tui-icon> -->
			</view>
		</tui-bottom-popup>
		<!--底部选择层-->

	</view>
</template>

<script>
	import tuiIcon from "@/components/icon/icon"
	import tuiTag from "@/components/tag/tag"
	import tuiBadge from "@/components/badge/badge"
	import tuiNomore from "@/components/nomore/nomore"
	import tuiButton from "@/components/button/button"
	import tuiTopDropdown from "@/components/top-dropdown/top-dropdown"
	import tuiBottomPopup from "@/components/bottom-popup/bottom-popup"
	import tuiNumberbox from "@/components/numberbox/numberbox"
	import {relativeTime} from '../../utils/common';
	import {api} from "../../../config/api"

	export default {
		components: {
			tuiIcon,
			tuiTag,
			tuiBadge,
			tuiNomore,
			tuiButton,
			tuiTopDropdown,
			tuiBottomPopup,
			tuiNumberbox
		},
		data() {
			return {
				id: 0,
				typeId: 1,
				info: {
					gallery:[],
                    labelsDesc:[],
					skuList: [],
				},
				inCartNumber:0, // 购物车数量
				collectFlag: 0, // 是否收藏
				specList: [],
				height: 64, //header高度
				top: 0, //标题图标距离顶部距离
				scrollH: 0, //滚动总高度
				opcity: 0,
				iconOpcity: 0.5,
				bannerIndex: 0,
				priceSeparator: "-",

				specIdValueIdsChecked: [],
				price: 0,
				skuListIndex: -1,
				buyNum: 1,


				checkedSkuInfo: {},
				checkedSkuListIndex: -1,
				checkedSpecIdValueIdsChecked: 0,
				topMenu: [{
					icon: "message",
					text: "消息",
					size: 26,
					badge: 3
				}, {
					icon: "home",
					text: "首页",
					size: 23,
					badge: 0,
					url: "/pages/mall/main"
				}, {
					icon: "people",
					text: "我的",
					size: 26,
					badge: 0,
					url: "/pages/my/main"
				}, {
					icon: "cart",
					text: "购物车",
					size: 23,
					badge: 2,
					url: "/pages/shopcart/main"
				}, {
					icon: "kefu",
					text: "客服小蜜",
					size: 26,
					badge: 0
				}, {
					icon: "feedback",
					text: "我要反馈",
					size: 23,
					badge: 0
				}, {
					icon: "share",
					text: "分享",
					size: 26,
					badge: 0
				}],

				goodsMode: "cart", // 两种类型,放入购物车,直接购买

				menuShow: false,
				popupShow: false,
				value: 1,
				collected: false
			}
		},
		onLoad: function(options) {
			let obj = {};
			// #ifdef MP-WEIXIN
			obj = uni.getMenuButtonBoundingClientRect();
			// #endif
			// #ifdef MP-BAIDU
			obj = swan.getMenuButtonBoundingClientRect();
			// #endif
			// #ifdef MP-ALIPAY
			my.hideAddToDesktopMenu();
			// #endif

			setTimeout(() => {
				uni.getSystemInfo({
					success: (res) => {
						this.width = obj.left || res.windowWidth;
						this.height = obj.top ? (obj.top + obj.height + 8) : (res.statusBarHeight + 44);
						this.top = obj.top ? (obj.top + (obj.height - 32) / 2) : (res.statusBarHeight + 6);
						this.scrollH = res.windowWidth
					}
				})
			}, 50)
		},
		onShow() {
			if (this.$root.$mp.query.id === undefined) {
				this.id = 1
			}else {
				this.id = this.$root.$mp.query.id;
			}
			this.getInfo();
			this.getUserGoodsData()
			this.getComments()
			// 初始化购物车大小
			this.cartTotal()
		},
		//商品转发
		onShareAppMessage() {
			return {
				title: this.info.title,
				path: "/pages/cominfo/main?id=" + this.id,
				imageUrl: this.info.cover //拿第一张商品的图片
			};
		},

		methods: {
			relativeTime:relativeTime,
			async getInfo() {
				const resp = await this.wxhttp.route(api.com.info,{id:parseInt(this.id)});
				const info = resp.data.info
				info.updatedAt = relativeTime(info.updatedAt)
				this.info = info;
				this.specList = resp.data.info.specList
				this.price = this.generatePice()
			},

			// 获取评论
			async getComments() {
				const resp = await this.wxhttp.route(api.comment.list,{gid:parseInt(this.id),typeId:this.typeId})
				if (resp.code === 0) {
					this.comments = resp.data.list.map(function (item, index) {
						item.createdAt = relativeTime(item.createdAt)
						return item;
					})
				}
			},
			// 购物车数量
			async cartTotal() {
				const resp = await this.wxhttp.route(api.cart.totalNum,{});
				if (resp.code !== 0) {
					return false
				}
				this.inCartNumber = resp.data
			},

			async getUserGoodsData() {
				if (this.tui.isLogin()) {
					const resp = await this.wxhttp.route(api.userGoods.stats,{gid:parseInt(this.id),typeId:this.typeId});
					if (resp.code === 0) {
						this.collectFlag = resp.data.isCollect
					}else {
						this.collectFlag = 0
					}
				}

			},

			// 收藏
			async onCollect() {
				if (!this.tui.checkLogin()) {
					return
				}

				// 如果之前是收藏，那么现在取消收藏
				if (this.collectFlag === 1) {
					const resp = await this.wxhttp.route(api.userGoods.uncollect,{gid:parseInt(this.id),typeId:this.typeId})
					if (resp.code !== 0) {
						this.tui.toast(resp.msg)
						return
					}
					this.collectFlag = 0;

				}else {
					const resp = await this.wxhttp.route(api.userGoods.collect  ,{gid:parseInt(this.id),typeId:this.typeId})
					if (resp.code !== 0) {
						this.tui.toast(resp.msg)
						return
					}
					this.collectFlag = 1;
				}
			},

			toCart: function() {
				uni.switchTab({
					url:"/pages/shopcart/main",
				})
			},

			bannerChange: function(e) {
				this.bannerIndex = e.detail.current
			},
			previewImage: function(e) {
				let index = e.currentTarget.dataset.index;
				uni.previewImage({
					current: this.feed.gallery[index],
					urls: this.feed.gallery
				})
			},
			back: function() {
				uni.navigateBack()
			},
			openMenu: function() {
				this.menuShow = true
			},
			closeMenu: function() {
				this.menuShow = false
			},
			showPopup: function() {
				this.popupShow = true
			},
			hidePopup: function() {
				this.popupShow = false
			},
			change: function(e) {
				this.value = e.value
			},
			collecting: function() {
				this.collected = !this.collected
			},
			common: function() {
				this.tui.toast("功能开发中~")
			},


			// 加入购物车
			async addCart() {
				// 如果没有登录
				if (!this.tui.checkLogin()) {
					return
				}

				// 如果没有选择规格,提示下
				if (!this.checkedComSkuInfo) {
					this.tui.toast("请选择商品规格")
					return false;
				}

				if (this.buyNum <= 0) {
					this.tui.toast("请选择商品数量")
					return false;
				}

				const resp = await this.wxhttp.route(api.cart.create,{
					comSkuId:parseInt(this.checkedComSkuInfo.id),
					num:parseInt(this.buyNum),
					typeId:parseInt(this.typeId),
				});
				// 添加失败
				if (resp.code !== 0) {
					this.tui.toast(resp.msg)
					return false
				}

				//添加成功后
				this.inCartNumber = resp.data.cnt;
				this.tui.toast("添加购物车成功","success")
				this.popupShow = false
			},
			// 勾选规格
			onSpecClick(specId,specValueId) {
				const comInfo = this.info

				let specIdValueIdsChecked = this.specIdValueIdsChecked

				// 取消选中
				if (specIdValueIdsChecked[specId] === specValueId) {
					delete specIdValueIdsChecked[specId]
				} else {
					// 选中
					specIdValueIdsChecked[specId] = specValueId
				}

				this.specIdValueIdsChecked = specIdValueIdsChecked

				console.log("onSpecClick specIdValueIdsChecked",specIdValueIdsChecked)

				this.initSpecList()
				// 判断是否有匹配的sku
				let specValueIdsChecked = [];
				for (let i = 0; i < specIdValueIdsChecked.length; i++) {
					if (typeof specIdValueIdsChecked[i] !== 'undefined') {
						specValueIdsChecked.push(specIdValueIdsChecked[i])
					}
				}

				this.specValueIdsChecked = specValueIdsChecked

				// 数据存在问题
				if (this.specList.length !== specValueIdsChecked.length) {
					this.tui.toast("规格数据出错")
					return
				}

				const matchResult = this.matchSku()
				console.log("onSkuMatchSuccess",matchResult)
				this.checkedComSkuInfo = matchResult.comSkuInfo
				this.checkedSkuListIndex = matchResult.skuListIndex,
				this.checkedSpecIdValueIdsChecked = this.specIdValueIdsChecked
				this.skuListIndex = matchResult.skuListIndex
				this.price = this.generatePice()
			},
			// 初始化规格
			initSpecList: function () {
				const comInfo = this.info
				let specIdValueIdsChecked = this.specIdValueIdsChecked
				comInfo.specList = comInfo.specList.map(function (item) {
					return {
						id: item.id,
						name: item.name,
						valueList: item.valueList.map(function (sub) {
							return {
								id: sub.id,
								name: sub.name,
								checked: sub.id === specIdValueIdsChecked[item.id]
							}
						})
					}
				})
				this.specList = comInfo.specList
			},
			// 计算价格
			generatePice() {
				const comInfo = this.info
				let price = 0

				// 选择第一个价位
				if (this.skuListIndex >= 0) {
					price = comInfo.skuList[this.skuListIndex].price
					return price
				}


				if (this.skuListIndex === -1) {
					// 当没有人选择时候,选择加个区间展示
					if (comInfo.skuList.length === 0) {
						return price
					}

					if (comInfo.skuList.length > 1) {
						// 如果是多条就区间
						let prices = comInfo.skuList.map(function (item) {
							return item.price
						})
						prices = prices.sort(function (a, b) {
							return a - b
						})
						// 如果价格相同
						if (prices[0] !== prices[prices.length - 1]) {
							price = `${prices[0]}${this.priceSeparator}${prices[prices.length - 1]}`
						}
						// 当没有人选择时候,选择第一个数据展示
					}else {
						console.log("com info")
						price = comInfo.skuList[0].price
					}
				}

				return price
			},
			// 直接购买
			async submit(){
				// 如果没有登录
				if (!this.tui.checkLogin()) {
					return
				}

				// 如果没有选择规格,提示下
				if (!this.checkedComSkuInfo) {
					this.tui.toast("请选择商品规格")
					return false;
				}

				if (this.buyNum <= 0) {
					this.tui.toast("请选择商品数量")
					return false;
				}

				const resp = await this.wxhttp.route(api.cart.create,{
					comSkuId:parseInt(this.checkedComSkuInfo.id),
					num:parseInt(this.buyNum),
					typeId:parseInt(this.typeId),
				});
				// 添加失败
				if (resp.code !== 0) {
					this.tui.toast(resp.msg)
					return false
				}
				this.popupShow = false
				uni.navigateTo({
					url: '/pages/submitOrder/main?comMode=buy&cartIds=' + JSON.stringify([resp.data.info.id])
				})
			},
			matchSku() {
				const comInfo = this.info
				const specValueSign = this.specValueIdsChecked.sort(function (a, b) {
					return a - b
				})

				console.log("matchSku specValueIdsChecked",this.specValueIdsChecked)
				console.log("matchSku specValueSign",specValueSign)


				let comSkuInfo = null
				let skuListIndex = null
				const specValueSign_string = JSON.stringify(specValueSign)
				console.log("specValueSign",specValueSign)

				console.log("specValueSign_string",specValueSign_string)
				for (let i = 0; i < comInfo.skuList.length; i++) {
					console.log("comInfo.sku_list[i].specValueSign",comInfo.skuList[i].specValueSign)

					if (comInfo.skuList[i].specValueSign === specValueSign_string) {
						comSkuInfo = comInfo.skuList[i]
						skuListIndex = i
						break
					}
				}
				return {
					comSkuInfo,
					skuListIndex
				}
			},
			changeNum(e) {
				this.buyNum = e.value
			},
			navigateToUrl(url){
				uni.navigateTo({
					url: url
				})
			},
		},
		onPageScroll(e) {
			let scroll = e.scrollTop <= 0 ? 0 : e.scrollTop;
			let opcity = scroll / this.scrollH;
			if (this.opcity >= 1 && opcity >= 1) {
				return;
			}
			this.opcity = opcity;
			this.iconOpcity = 0.5 * (1 - opcity < 0 ? 0 : 1 - opcity)
		}
	}
</script>

<style>
	/* icon 也可以使用组件*/
	@import "../../static/style/icon.css";

	page {
		background: #f7f7f7;
	}

	.tui-operation .operation-button {
		position:relative;
		display:block;
		margin-left:0px;
		margin-right:0px;
		padding-left:0px;
		padding-right:0px;
		box-sizing:border-box;
		font-size:18px;
		text-align:center;
		text-decoration:none;
		line-height:1;
		border-radius:0px;
		-webkit-tap-highlight-color:transparent;
		overflow:hidden;
		color: #999;
		background-color: transparent;
	}


	.tui-operation .operation-button::after {
		display: block;
		border: none;
	}

	.container {
		padding-bottom: 110rpx;
	}

	.tui-header-box {
		width: 100%;
		position: fixed;
		left: 0;
		top: 0;
		z-index: 9998;
	}

	.tui-header {
		width: 100%;
		font-size: 18px;
		line-height: 18px;
		font-weight: 500;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.tui-header-icon {
		position: fixed;
		top: 0;
		left: 10px;
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		height: 32px;
		transform: translateZ(0);
		z-index: 99999;
	}



	.tui-header-icon .tui-badge {
		background: #e41f19 !important;
		position: absolute;
		right: -4px;
	}

	.tui-icon-ml {
		margin-left: 20rpx;
	}

	.tui-icon {
		border-radius: 16px;
	}


	.tui-icon-back {
		height: 32px !important;
		width: 32px !important;
		display: block !important;
	}

	.tui-header-icon .tui-icon-more-fill {
		height: 20px !important;
		width: 20px !important;
		padding: 6px !important;
		display: block !important;
	}

	.tui-banner-swiper {
		position: relative;
	}

	.tui-banner-swiper .tui-tag-class {
		position: absolute;
		color: #fff;
		bottom: 30rpx;
		right: 0;
	}

	.tui-slide-image {
		width: 100%;
		display: block;
	}

	/*顶部菜单*/

	.tui-menu-box {
		box-sizing: border-box;
	}

	.tui-menu-header {
		font-size: 34rpx;
		color: #fff;
		height: 32px;
		display: flex;
		align-items: center;
	}

	.tui-top-dropdown {
		z-index: 9999 !important;
	}

	.tui-menu-itembox {
		color: #fff;
		padding: 40rpx 10rpx 0 10rpx;
		box-sizing: border-box;
		display: flex;
		flex-wrap: wrap;
		font-size: 26rpx;
	}

	.tui-menu-item {
		width: 22%;
		height: 160rpx;
		border-radius: 24rpx;
		display: flex;
		align-items: center;
		flex-direction: column;
		justify-content: center;
		background: rgba(0, 0, 0, 0.4);
		margin-right: 4%;
		margin-bottom: 4%;
	}

	.tui-menu-item:nth-of-type(4n) {
		margin-right: 0;
	}

	.tui-badge-box {
		position: relative;
	}

	.tui-badge-box .tui-badge-class {
		position: absolute;
		top: -8px;
		right: -8px;
	}

	.tui-msg-badge {
		top: -10px;
	}

	.tui-icon-up {
		position: relative;
		display: inline-block;
		left: 50%;
		transform: translateX(-50%);
	}

	.tui-menu-text {
		padding-top: 12rpx;
	}

	.tui-opcity .tui-menu-text,
	.tui-opcity .tui-badge-box {
		opacity: 0.5;
		transition: opacity 0.2s ease-in-out;
	}

	/*顶部菜单*/

	/*内容 部分*/

	.tui-padding {
		padding: 0 30rpx;
		box-sizing: border-box;
	}

	/* #ifdef H5 */
	.tui-ptop {
		padding-top: 44px;
	}

	/* #endif */

	.tui-size {
		font-size: 24rpx;
		line-height: 24rpx;
	}

	.tui-gray {
		color: #999;
	}

	.tui-icon-red {
		color: #ff201f;
	}

	.tui-border-radius {
		border-bottom-left-radius: 24rpx;
		border-bottom-right-radius: 24rpx;
		overflow: hidden;
	}

	.tui-radius-all {
		border-radius: 24rpx;
		overflow: hidden;
	}

	.tui-mtop {
		margin-top: 26rpx;
	}

	.tui-pro-detail {
		box-sizing: border-box;
		color: #333;
	}

	.tui-product-title {
		background: #fff;
		padding: 30rpx 0;
	}

	.tui-pro-pricebox {
		display: flex;
		align-items: center;
		justify-content: space-between;
		color: #ff201f;
		font-size: 36rpx;
		font-weight: bold;
		line-height: 44rpx;
	}

	.tui-pro-price {
		display: flex;
		align-items: center;
	}

	.tui-pro-price .tui-tag-class {
		transform: scale(0.7);
		transform-origin: center center;
		line-height: 24rpx;
		font-weight: normal;
	}

	.tui-price {
		font-size: 58rpx;
	}

	.tui-original-price {
		font-size: 26rpx;
		line-height: 26rpx;
		padding: 10rpx 30rpx;
		box-sizing: border-box;
	}

	.tui-line-through {
		text-decoration: line-through;
	}

	.tui-collection {
		color: #333;
		display: flex;
		align-items: center;
		flex-direction: column;
		justify-content: center;
		height: 44rpx;
	}

	.tui-scale {
		transform: scale(0.7);
		transform-origin: center center;
		line-height: 24rpx;
		font-weight: normal;
	}

	.tui-icon-collection {
		line-height: 20px !important;
		margin-bottom: 0 !important;

	}

	.tui-pro-titbox {
		font-size: 32rpx;
		font-weight: 500;
		position: relative;
		padding: 0 150rpx 0 30rpx;
		box-sizing: border-box;
	}

	.tui-pro-title {
		padding-top: 20rpx;
	}

	.tui-share-btn {
		display: block;
		background: none;
		margin: 0;
		padding: 0;
		border-radius: 0;
	}

	.tui-tag-share {
		display: flex;
		align-items: center;
	}

	.tui-share-position {
		position: absolute;
		right: 0;
		top: 30rpx;
	}

	.tui-share-text {
		padding-left: 8rpx;
	}

	.tui-sub-title {
		padding: 20rpx 0;
	}

	.tui-sale-info {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding-top: 30rpx;
	}

	.tui-discount-box {
		background: #fff;
	}

	.tui-list-cell {
		position: relative;
		display: flex;
		align-items: center;
		font-size: 26rpx;
		line-height: 26rpx;
		padding: 36rpx 30rpx;
		box-sizing: border-box;
	}

	.tui-right {
		position: absolute;
		right: 30rpx;
		top: 30rpx;
	}

	.tui-top40 {
		top: 40rpx !important;
	}

	.tui-bold {
		font-weight: bold;
	}

	.tui-list-cell::after {
		content: '';
		position: absolute;
		border-bottom: 1rpx solid #eaeef1;
		-webkit-transform: scaleY(0.5);
		transform: scaleY(0.5);
		bottom: 0;
		right: 0;
		left: 126rpx;
	}

	.tui-last::after {
		border-bottom: 0 !important;
	}

	.tui-tag-coupon-box {
		display: flex;
		align-items: center;
	}

	.tui-tag-coupon-box .tui-tag-class {
		margin-right: 20rpx;
	}


	.tui-cell-title {
		width: 66rpx;
		padding-right: 30rpx;
		flex-shrink: 0;
	}

	.tui-promotion-box {
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		padding: 10rpx 0;
		width: 74%;
	}

	.tui-promotion-box .tui-tag-class {
		display: inline-block !important;
		transform: scale(0.8);
		transform-origin: 0 center;
	}

	/* .tui-inline-block {
		display: inline-block !important;
		transform: scale(0.8);
		transform-origin: 0 center;
	} */

	.tui-basic-info {
		background: #fff;
	}

	.tui-addr-box {
		width: 76%;
	}

	.tui-addr-item {
		padding: 10rpx;
		line-height: 34rpx;
	}

	.tui-guarantee {
		background: #fdfdfd;
		display: flex;
		flex-wrap: wrap;
		padding: 20rpx 30rpx 30rpx 30rpx;
		font-size: 24rpx;
	}

	.tui-guarantee-item {
		color: #999;
		padding-right: 30rpx;
		padding-top: 10rpx;
	}

	.tui-pl {
		padding-left: 4rpx;
	}

	.tui-cmt-box {
		background: #fff;
	}

	.tui-between {
		justify-content: space-between !important;
	}

	.tui-cmt-all {
		color: #ff201f;
		padding-right: 8rpx;
	}

	.tui-cmt-content {
		font-size: 26rpx;
	}

	.tui-cmt-user {
		display: flex;
		align-items: center;
	}

	.tui-acatar {
		width: 60rpx;
		height: 60rpx;
		border-radius: 30rpx;
		display: block;
		margin-right: 16rpx;
	}

	.tui-cmt {
		padding: 14rpx 0;
	}

	.tui-attr {
		font-size: 24rpx;
		color: #999;
		padding: 6rpx 0;
	}

	.tui-cmt-btn {
		padding: 50rpx 0 30rpx 0;
		box-sizing: border-box;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.tui-tag-cmt {
		min-width: 130rpx;
		padding: 20rpx 52rpx !important;
		font-size: 26rpx !important;
		display: inline-block;
	}

	.tui-nomore-box {
		padding-top: 10rpx;
	}

	.tui-product-img {
		transform: translateZ(0);
	}

	.tui-product-img image {
		width: 100%;
		display: block;
	}

	/*底部操作栏*/

	.tui-col-7 {
		width: 58.33333333%;
	}

	.tui-col-5 {
		width: 41.66666667%;
	}

	.tui-operation {
		width: 100%;
		height: 100rpx;
		/* box-sizing: border-box; */
		background: rgba(255, 255, 255, 0.98);
		position: fixed;
		display: flex;
		align-items: center;
		justify-content: space-between;
		z-index: 10;
		bottom: 0;
		left: 0;
		padding-bottom: env(safe-area-inset-bottom);
	}

	.tui-safearea-bottom {
		width: 100%;
		height: env(safe-area-inset-bottom);
	}

	.tui-operation::before {
		content: '';
		position: absolute;
		top: 0;
		right: 0;
		left: 0;
		border-top: 1rpx solid #eaeef1;
		-webkit-transform: scaleY(0.5);
		transform: scaleY(0.5);
	}

	.tui-operation-left {
		display: flex;
		align-items: center;
	}

	.tui-operation-item {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-direction: column;
		position: relative;
	}

	.tui-operation-text {
		font-size: 22rpx;
		color: #333;
	}

	.tui-opacity {
		opacity: 0.5;
	}

	.tui-scale-small {
		transform: scale(0.9);
		transform-origin: center center;
	}

	.tui-operation-right {
		height: 100rpx;
		/* box-sizing: border-box; */
		padding-top: 0;
	}

	.tui-right-flex {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.tui-btnbox-4 .tui-btn-class {
		width: 90% !important;
		display: block !important;
		font-size: 28rpx !important;
	}

	.tui-operation .tui-badge-class {
		position: absolute;
		top: -6rpx;
		/* #ifdef H5 */
		transform: translateX(50%)
			/* #endif  */
	}

	.tui-flex-1 {
		flex: 1;
	}

	/*底部操作栏*/

	/*底部选择弹层*/

	.tui-popup-class {
		border-top-left-radius: 24rpx;
		border-top-right-radius: 24rpx;
		padding-bottom: env(safe-area-inset-bottom);
	}

	.tui-popup-box {
		position: relative;
		padding: 30rpx 0 100rpx 0;
	}

	.tui-popup-btn {
		width: 100%;
		position: absolute;
		left: 0;
		bottom: 0;
	}

	.tui-popup-btn .tui-btn-class {
		width: 90% !important;
		display: block !important;
		font-size: 28rpx !important;
	}

	.tui-icon-close {
		position: absolute;
		top: 30rpx;
		right: 30rpx;
	}

	.tui-product-box {
		display: flex;
		align-items: flex-end;
		font-size: 24rpx;
		padding-bottom: 30rpx;
	}

	.tui-popup-img {
		height: 200rpx;
		width: 200rpx;
		border-radius: 24rpx;
		display: block;
	}

	.tui-popup-price {
		padding-left: 20rpx;
		padding-bottom: 8rpx;
	}

	.tui-amount {
		color: #ff201f;
		font-size: 36rpx;
	}

	.tui-number {
		font-size: 24rpx;
		line-height: 24rpx;
		padding-top: 12rpx;
		color: #999;
	}

	.tui-popup-scroll {
		height: 450rpx;
		font-size: 26rpx;
	}

	.tui-scrollview-box {
		padding: 0 30rpx 60rpx 30rpx;
		box-sizing: border-box;
	}

	.tui-attr-title {
		padding: 10rpx 0;
		color: #333;
	}

	.tui-attr-box {
		font-size: 0;
		padding: 20rpx 0;
	}

	.tui-attr-item {
		max-width: 100%;
		min-width: 200rpx;
		height: 64rpx;
		display: -webkit-inline-flex;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		background: #f7f7f7;
		padding: 0 26rpx;
		box-sizing: border-box;
		border-radius: 32rpx;
		margin-right: 20rpx;
		margin-bottom: 20rpx;
		font-size: 26rpx;
	}

	.tui-attr-active {
		background: #fcedea !important;
		color: #e41f19;
		font-weight: bold;
		position: relative;
	}

	.tui-attr-active::after {
		content: "";
		position: absolute;
		border: 1rpx solid #e41f19;
		width: 100%;
		height: 100%;
		border-radius: 40rpx;
		left: 0;
		top: 0;
	}

	.tui-number-box {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 20rpx 0 30rpx 0;
		box-sizing: border-box;
	}

	/*底部选择弹层*/
</style>
