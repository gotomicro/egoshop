<template>
	<view class="container">
		<view class="tui-box">
			<tui-list-cell :arrow="true" :last="true" :radius="true" @click="chooseAddr">
				<view class="tui-address">
					<view v-if="address.id > 0">
						<view class="tui-userinfo">
							<text class="tui-name">{{address.name}}</text> {{address.mobile}}
						</view>
						<view class="tui-addr">
							<view class="tui-addr-tag" v-if="address.isDefault">默认</view>
							<text>{{address.region+address.detail}}</text>
						</view>
					</view>
					<view class="tui-none-addr" v-else>
						<image src="/static/images/index/map.png" class="tui-addr-img" mode="widthFix"></image>
						<text>选择收货地址</text>
					</view>
				</view>
				<view class="tui-bg-img"></view>
			</tui-list-cell>
			<view class="tui-top tui-goods-info">
				<tui-list-cell :hover="false" :lineLeft="false">
					<view class="tui-goods-title">
						商品信息
					</view>
				</tui-list-cell>
				<block v-for="(item,index) in listData" :key="index">
					<tui-list-cell :hover="false" padding="0">
						<view class="tui-goods-item">
							<image :src="item.cover" class="tui-goods-img"></image>
							<view class="tui-goods-center">
								<view class="tui-goods-name">{{item.title}}</view>
								<view class="tui-goods-attr">{{item.subTitle}}</view>
							</view>
							<view class="tui-price-right">
								<view>￥{{item.price}}</view>
								<view>x{{item.num}}</view>
							</view>
						</view>
					</tui-list-cell>
				</block>
				<tui-list-cell :hover="false">
					<view class="tui-padding tui-flex">
						<view>商品总额</view>
						<view>￥{{calculate.payAmount}}</view>
					</view>
				</tui-list-cell>
				<tui-list-cell :arrow="hasCoupon" :hover="hasCoupon" >
					<view class="tui-padding tui-flex">
						<view>优惠券</view>
						<view :class="{'tui-color-red':hasCoupon}">{{hasCoupon?"满5减1":'没有可用优惠券'}}</view>
					</view>
				</tui-list-cell>
				<tui-list-cell :hover="false">
					<view class="tui-padding tui-flex">
						<view>配送费</view>
						<view>￥0.00</view>
					</view>
				</tui-list-cell>
				<tui-list-cell :hover="false" :lineLeft="false" padding="0">
					<view class="tui-remark-box tui-padding tui-flex">
						<view>订单备注</view>
						<input type="text" class="tui-remark" placeholder="选填: 请先和商家协商一致"
							   placeholder-class="tui-phcolor" @input="bindRemark"/>
					</view>
				</tui-list-cell>
				<tui-list-cell :hover="false" :last="true">
					<view class="tui-padding tui-flex tui-total-flex">
						<view class="tui-flex-end tui-color-red">
							<view class="tui-black">合计： </view>
							<view class="tui-size-26">￥</view>
							<view class="tui-price-large">{{calculate.payAmount}}</view>
						</view>
					</view>
				</tui-list-cell>
			</view>

			<view class="tui-top">
				<tui-list-cell :last="true" :hover="insufficient" :radius="true" :arrow="insufficient">
					<view class="tui-flex">
						<view class="tui-balance">余额支付<text class="tui-gray">(￥0.00)</text></view>
						<switch color="#30CC67" class="tui-scale-small" v-show="!insufficient" />
						<view class="tui-pr-30 tui-light-dark" v-show="insufficient">余额不足, 去充值</view>
					</view>
				</tui-list-cell>
			</view>
		</view>
		<view class="tui-safe-area"></view>
		<view class="tui-tabbar">
			<view class="tui-flex-end tui-color-red tui-pr-20">
				<view class="tui-black">实付金额: </view>
				<view class="tui-size-26">￥</view>
				<view class="tui-price-large">{{calculate.payAmount}}</view>
			</view>
			<view class="tui-pr25">
				<tui-button width="200rpx" height="70rpx" type="danger" shape="circle" @click="btnPay">确认支付</tui-button>
			</view>
		</view>

	</view>
</template>

<script>
	import tuiButton from "@/components/extend/button/button"
	import tuiListCell from "@/components/list-cell/list-cell"
	import tuiBottomPopup from "@/components/bottom-popup/bottom-popup"
	import {api} from "../../../config/api"

	export default {
		components: {
			tuiButton,
			tuiListCell,
			tuiBottomPopup
		},
		data() {
			return {
				hasCoupon: false,
				insufficient: true,
				address: {}, // 地址
				calculate: {}, // 费用计算
				listData: {}, // 购物车
				cartIds: {}, // 购物车id
				remark: "", // 提示信息
			}
		},
		async onShow() {
			if (uni.getStorageSync("addressId")) {
				this.addressId = uni.getStorageSync("addressId");
			}
			// 初始化用户地址
			this.initAddress()
			// 购物车id
			let cartIds = this.$root.$mp.query.cartIds;
			// cartIds = "[1]"
			if (cartIds === undefined) {
				uni.showToast({
					title: "购物车id有问题", //提示的内容,
					duration: 2000, //延迟时间,
					icon: "none",
					mask: true, //显示透明蒙层，防止触摸穿透,
					success: res => {
					}
				});
				return
			}
			// 获得购物车id
			this.cartIds = JSON.parse(cartIds);
			console.log("cartIds",this.cartIds)

			// 先初始化购物车
			const isInitCart = await this.initCartList()
			if (!isInitCart) {
				uni.showToast({
					title: "支付商品状态已变", //提示的内容,
					icon: "none", //图标,
					duration: 1500, //延迟时间,
					mask: false, //显示透明蒙层，防止触摸穿透,
					success: res => {}
				});
				return
			}
			this.initCalculate()

		},
		methods: {
			async initAddress() {
				let resp
				if (this.addressId !== undefined && this.addressId > 0) {
					resp = await this.wxhttp.route(api.address.info,{
						id: this.addressId,
					});
				}else {
					resp = await this.wxhttp.route(api.address.default,{});
				}

				if (resp.code !== 0) {
					return false
				}
				this.addressId= resp.data.id
				this.address = resp.data
				return true
			},
			// 计算费用
			async initCalculate() {
				const resp = await this.wxhttp.route(api.pay.calculate,{
					cartIds: this.cartIds,
					addressId: this.addressId
				});
				if (resp.code !== 0) {
					this.tui.toast(resp.msg,1500,"error")
					return false
				}
				this.calculate = resp.data
			},

			async initCartList() {
				const resp = await this.wxhttp.route(api.cart.list,{
					cartIds: this.cartIds,
				});

				if (resp.code !== 0) {
					return false
				}
				this.listData = resp.data.list
				return true
			},
			// 订单备注
			bindRemark(e) {
				this.remark = e.detail.value
			},

			async btnPay() {
				const that = this

				if (!this.addressId) {
					uni.showToast({
						title: "请选择收货地址", //提示的内容,
						icon: "error", //图标,
						duration: 1500, //延迟时间,
						mask: false, //显示透明蒙层，防止触摸穿透,
						success: res => {
						}
					});
					return
				}
				// 创建订单
				const resp = await this.wxhttp.route(api.order.create, {
					'comMode': this.comMode,
					'addressId': this.addressId,
					'cartIds': this.cartIds,
					'remark': this.remark,
				});

				if (resp.code !== 0) {
					uni.showToast({
						title: "创建订单失败,错误信息" + resp.msg, //提示的内容,
						icon: "error", //图标,
						duration: 1500, //延迟时间,
						mask: false, //显示透明蒙层，防止触摸穿透,
						success: res => {
						}
					});
					return false
				}
				// 购买
				const resp2 = await this.wxhttp.route(api.order.pay, {
					'orderType': "comBuy",
					'paySn': resp.data.paySn,
					'paymentCode': "wechat",
					'paymentChannel': "wechat_mini",
				});
				if (resp2.code !== 0) {
					uni.showToast({
						title: "支付失败:" + resp.msg, //提示的内容,
						icon: "error", //图标,
						duration: 1500, //延迟时间,
						mask: false, //显示透明蒙层，防止触摸穿透,
						success: res => {
						}
					});
					uni.navigateBack({delta: this.delta})
					return false
				}
				uni.requestPayment({
					'timeStamp': resp2.data.timeStamp,
					'nonceStr': resp2.data.nonceStr,
					'package': resp2.data.package,
					'signType': resp2.data.signType,
					'paySign': resp2.data.paySign,
					'success': function () {
						uni.redirectTo({
							url: `/pages/payResult/main?payAmount=${resp2.data.payAmount}&orderId=${resp.data.orderId}&paySn=${resp2.data.paySn}`
						})
					},
					'fail': function (res) {
						uni.showToast({
							title: "支付被取消:" + resp.msg, //提示的内容,
							icon: "error", //图标,
							duration: 1500, //延迟时间,
							mask: false, //显示透明蒙层，防止触摸穿透,
							success: res => {
							}
						});
						setTimeout(function () {
							uni.redirectTo({
								url: `/pages/submitOrder/main?id=${resp2.data.orderfId}`
							})
						}, 1000)
					}
				})
			},
			chooseAddr() {
				uni.navigateTo({
					url: "/pages/address/main"
				})
			},
			// btnPay() {
			// 	uni.navigateTo({
			// 		url: "/pages/success/main"
			// 	})
			// }
		}
	}
</script>

<style>
	.container {
		padding-bottom: 98rpx;
	}

	.tui-box {
		padding: 20rpx 0 118rpx;
		box-sizing: border-box;
	}

	.tui-address {
		min-height: 80rpx;
		padding: 10rpx 0;
		box-sizing: border-box;
		position: relative;
	}

	.tui-userinfo {
		font-size: 30rpx;
		font-weight: 500;
		line-height: 30rpx;
		padding-bottom: 12rpx;
	}

	.tui-name {
		padding-right: 40rpx;
	}

	.tui-addr {
		font-size: 24rpx;
		word-break: break-all;
		padding-right: 25rpx;
	}

	.tui-addr-tag {
		padding: 5rpx 8rpx;
		flex-shrink: 0;
		background: #EB0909;
		color: #fff;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		font-size: 25rpx;
		line-height: 25rpx;
		transform: scale(0.8);
		transform-origin: 0 center;
		border-radius: 6rpx;
	}

	.tui-bg-img {
		position: absolute;
		width: 100%;
		height: 8rpx;
		left: 0;
		bottom: 0;
		background: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAL0AAAAECAMAAADszM6/AAAAOVBMVEUAAAAVqfH/fp//vWH/vWEVqfH/fp8VqfH/fp//vWEVqfH/fp8VqfH/fp//vWH/vWEVqfH/fp//vWHpE7b6AAAAEHRSTlMA6urqqlVVFRUVq6upqVZUDT4vVAAAAEZJREFUKM/t0CcOACAQRFF6r3v/w6IQJGwyDsPT882IQzQE0E3chToByjG5LwMgLZN3TQATmdypCciBya0cgOT3/h//9PgF49kd+6lTSIIAAAAASUVORK5CYII=") repeat;
	}

	.tui-top {
		margin-top: 20rpx;
		overflow: hidden;
	}

	.tui-goods-title {
		font-size: 28rpx;
		display: flex;
		align-items: center;
	}

	.tui-padding {
		box-sizing: border-box;
	}

	.tui-goods-item {
		width: 100%;
		padding: 20rpx 30rpx;
		box-sizing: border-box;
		display: flex;
		justify-content: space-between;
	}

	.tui-goods-img {
		width: 180rpx;
		height: 180rpx;
		display: block;
		flex-shrink: 0;
	}

	.tui-goods-center {
		flex: 1;
		padding: 20rpx 8rpx;
		box-sizing: border-box;
	}

	.tui-goods-name {
		max-width: 310rpx;
		word-break: break-all;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
		font-size: 26rpx;
		line-height: 32rpx;
	}

	.tui-goods-attr {
		font-size: 22rpx;
		color: #888888;
		line-height: 32rpx;
		padding-top: 20rpx;
		word-break: break-all;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
	}

	.tui-price-right {
		text-align: right;
		font-size: 24rpx;
		color: #888888;
		line-height: 30rpx;
		padding-top: 20rpx;
	}

	.tui-flex {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: space-between;
		font-size: 26rpx;
	}
	.tui-total-flex{
		justify-content: flex-end;
	}

	.tui-color-red,
	.tui-invoice-text {
		color: #E41F19;
		padding-right: 30rpx;
	}

	.tui-balance {
		font-size: 28rpx;
		font-weight: 500;
	}

	.tui-black {
		color: #222;
		line-height: 30rpx;
	}

	.tui-gray {
		color: #888888;
		font-weight: 400;
	}

	.tui-light-dark {
		color: #666;
	}

	.tui-goods-price {
		display: flex;
		align-items: center;
		padding-top: 20rpx;
	}

	.tui-size-26 {
		font-size: 26rpx;
		line-height: 26rpx;
	}

	.tui-price-large {
		font-size: 34rpx;
		line-height: 32rpx;
		font-weight: 600;
	}

	.tui-flex-end {
		display: flex;
		align-items: flex-end;
		padding-right: 0;
	}

	.tui-phcolor {
		color: #B3B3B3;
		font-size: 26rpx;
	}

	/* #ifndef H5 */
	.tui-remark-box {
		padding: 22rpx 30rpx;
	}

	/* #endif */
	/* #ifdef H5 */
	.tui-remark-box {
		padding: 26rpx 30rpx;
	}

	/* #endif */

	.tui-remark {
		flex: 1;
		font-size: 26rpx;
		padding-left: 64rpx;
	}

	.tui-scale-small {
		transform: scale(0.8);
		transform-origin: 100% center;
	}

	.tui-scale-small .wx-switch-input {
		margin: 0 !important;
	}

	/* #ifdef H5 */
	>>>uni-switch .uni-switch-input {
		margin-right: 0 !important;
	}

	/* #endif */
	.tui-tabbar {
		width: 100%;
		height: 98rpx;
		background: #fff;
		position: fixed;
		left: 0;
		bottom: 0;
		display: flex;
		align-items: center;
		justify-content: flex-end;
		font-size: 26rpx;
		box-shadow: 0 0 1px rgba(0, 0, 0, .3);
		padding-bottom: env(safe-area-inset-bottom);
		z-index: 999;
	}

	.tui-pr-30 {
		padding-right: 30rpx;
	}

	.tui-pr-20 {
		padding-right: 20rpx;
	}

	.tui-none-addr {
		height: 80rpx;
		padding-left: 5rpx;
		display: flex;
		align-items: center;
	}

	.tui-addr-img {
		width: 36rpx;
		height: 46rpx;
		display: block;
		margin-right: 15rpx;
	}


	.tui-pr25 {
		padding-right: 25rpx;
	}

	.tui-safe-area {
		height: 1rpx;
		padding-bottom: env(safe-area-inset-bottom);
	}
</style>
