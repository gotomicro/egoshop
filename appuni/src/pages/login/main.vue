<template>
	<view class="container">
		<view class="tui-bg-box">
			<image src="../../static/images/login/bg_login.png" class="tui-bg-img"></image>
			<image src="../../static/images/my/mine_def_touxiang_3x.png" class="tui-photo"></image>
			<view class="tui-login-name">希寻名宠</view>
		</view>
		<view class="tui-login-from">
			<button class="btn-success tui-btn-submit" hover-class="btn-hover" open-type="getUserInfo" @getuserinfo="doLogin">微信登录</button>
			<view class="tui-protocol" hover-class="opcity" :hover-stay-time="150">点击"登录"即表示已同意
				<text class="tui-protocol-red" @tap="protocol">《用户协议》</text>
			</view>
		</view>
	</view>
</template>

<script>
	import {
		mapState,
		mapMutations,
		mapActions
	} from 'vuex'
	import tuiIcon from "@/components/icon/icon"
	import tuiButton from "@/components/button/button"
	const util = require('../../utils/util.js')
	import {api} from "../../../config/api"
	var qcloud = require("wafer2-client-sdk/index.js");
	export default {
		components: {
			tuiIcon,
			tuiButton
		},
		mounted() {
			qcloud.setLoginUrl(api.auth.login.url);
			this.url = this.$root.$mp.query.url;
		},
		data() {
			return {
				disabled: false,
				btnText: "获取验证码",
				mobile: "",
				type: "primary",
				code: "",
				url: "",
			}
		},
		methods: {
			doLogin() {
				uni.showLoading({
					title: "登录中...", //提示的内容,
					mask: true, //显示透明蒙层，防止触摸穿透,
					success: res => {}
				});
				// const session = qcloud.Session.get();
				var that = this
				// if (session) {
				// 	// 第二次登录
				// 	// 或者本地已经有登录态
				// 	// 可使用本函数更新登录态
				// 	qcloud.loginWithCode({
				// 		success: res => {
				// 			// this.setData({ userInfo: res, logged: true });
				// 			uni.setStorageSync("key", "value");
				// 		},
				// 		fail: err => {
				// 			console.error(err);
				// 		}
				// 	});
				// } else {
					// 首次登录
					qcloud.login({
						success: res => {
							uni.hideLoading();
							uni.setStorageSync("userInfo", res);
							var openId = res.openId;
							uni.setStorageSync("openId", openId);
							that.tui.setToken(qcloud.Session.get().skey)
							if (this.url !== "") {
								uni.navigateTo({
									url: "/"+this.url
								})
								return
							}
							uni.navigateBack();
						},
						fail: err => {
							console.log(err);
							uni.hideLoading();
							if (this.url !== "") {
								uni.navigateTo({
									url: "/"+this.url
								})
								return
							}
							uni.navigateBack();
						}
					});
				// }
			},
			protocol: function() {
				uni.navigateTo({
					url: '../about/about'
				})
			}
		}
	}
</script>

<style>
	page {
		background: #fff;
	}

	.tui-bg-box {
		width: 100%;
		box-sizing: border-box;
		position: relative;
		padding-top: 44rpx;
	}

	.tui-photo {
		height: 138rpx;
		width: 138rpx;
		display: block;
		margin: 10rpx auto 0 auto;
		border-radius: 50%;
		position: relative;
		z-index: 2;
	}

	.tui-login-name {
		width: 128rpx;
		height: 40rpx;
		font-size: 30rpx;
		color: #fff;
		margin: 36rpx auto 0 auto;
		text-align: center;
		position: relative;
		z-index: 2;
	}

	.tui-bg-img {
		width: 100%;
		height: 346rpx;
		display: block;
		position: absolute;
		top: 0;
	}

	.tui-login-from {
		width: 100%;
		padding: 128rpx 104rpx 0 104rpx;
		box-sizing: border-box;
	}

	.tui-input {
		font-size: 32rpx;
		flex: 1;
		display: inline-block;
		padding-left: 32rpx;
		box-sizing: border-box;
		overflow: hidden;
	}

	.tui-line-cell {
		padding: 27rpx 0;
		display: -webkit-flex;
		display: flex;
		-webkiit-align-items: center;
		align-items: center;
		position: relative;
		box-sizing: border-box;
		overflow: hidden;
	}

	.tui-line-cell::after {
		content: '';
		position: absolute;
		border-bottom: 1rpx solid #e0e0e0;
		-webkit-transform: scaleY(0.5);
		transform: scaleY(0.5);
		bottom: 0;
		right: 0;
		left: 0;
	}

	.tui-top28 {
		margin-top: 28rpx;
	}

	.tui-btn-class {
		width: 196rpx !important;
		height: 54rpx !important;
		border-radius: 27rpx !important;
		font-size: 28rpx !important;
		padding: 0 !important;
		line-height: 54rpx !important;
	}

	.tui-btn-submit {
		margin-top: 100rpx;
	}

	.tui-protocol {
		color: #333;
		font-size: 24rpx;
		text-align: center;
		width: 100%;
		margin-top: 29rpx;
	}

	.tui-protocol-red {
		color: #f54f46;
	}
</style>
