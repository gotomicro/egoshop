<template>
	<view class="tui-safe-area">
		<view class="tui-address" v-if="listData.length!==0">
			<block v-for="(item,index) in listData" :key="index">
				<tui-list-cell padding="0">
					<view class="tui-address-flex">
						<view class="tui-address-left"  @click="selAddress(item.id)">
							<view class="tui-address-main">
								<view class="tui-address-name tui-ellipsis">{{item.name}}</view>
								<view class="tui-address-tel">{{item.mobile}}</view>
							</view>
							<view class="tui-address-detail">
								<view class="tui-address-label" v-if="item.isDefault">默认</view>
								<view class="tui-address-label">{{item.typeName}}</view>
								<text>{{item.region}}{{item.detail}}</text>
							</view>
						</view>
						<view class="tui-address-imgbox" @click="updateAddr(item.id)">
							<image class="tui-address-img" src="/static/images/mall/my/icon_addr_edit.png" />
						</view>
					</view>
				</tui-list-cell>
			</block>
		</view>
		<view :class="{'tui-order-list':scrollTop>=0}" v-else>
			<tui-tips :fixed="false" imgUrl="/static/images/toast/img_noorder_3x.png">暂无地址</tui-tips>
		</view>
		<!-- 新增地址 -->
		<view class="tui-address-new">
			<tui-button type="danger" height="88rpx" @click="createAddr">+ 新增收货地址</tui-button>
		</view>
	</view>
</template>

<script>
	import tuiButton from "@/components/extend/button/button"
	import tuiListCell from "@/components/list-cell/list-cell"
	import tuiTips from "@/components/extend/tips/tips"

	import {api} from "../../../config/api"

	export default {
		components: {
			tuiButton,
			tuiListCell,
			tuiTips,
		},
		data() {
			return {
				addressList: [],
				listData: [],
			}
		},
		onLoad: function(options) {

		},
		onShow() {
			// 获取所有地址列表
			this.getAddressList();
		},
		methods: {
			async getAddressList() {
				const resp = await this.wxhttp.route(api.address.list,{});
				if (resp.code === 0 ) {
					this.listData = resp.data.list;
				}
			},
			createAddr() {
				uni.navigateTo({
					url: "/pages/editAddress/main"
				})
			},
			updateAddr(id) {
				uni.navigateTo({
					url: "/pages/editAddress/main?id=" + id
				})
			},
			selAddress(id) {
				uni.setStorageSync("addressId", id);
				// uni.redirectTo({ url: "/pages/order/main?id=" + id });
				uni.navigateBack({
					delta: 1 //返回的页面数，如果 delta 大于现有页面数，则返回到首页,
				});
			},
		}
	}
</script>

<style>
	.tui-order-list {
		margin-top: 80rpx;
	}
	.tui-address {
		width: 100%;
		padding-top: 20rpx;
		padding-bottom: 160rpx;
	}
	.tui-address-flex {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.tui-address-main {
		width: 600rpx;
		height: 70rpx;
		display: flex;
		font-size: 30rpx;
		line-height: 86rpx;
		padding-left: 30rpx;
	}

	.tui-address-name {
		width: 120rpx;
		height: 60rpx;
	}

	.tui-address-tel {
		margin-left: 10rpx;
	}

	.tui-address-detail {
		font-size: 24rpx;
		word-break: break-all;
		padding-bottom: 25rpx;
		padding-left: 25rpx;
		padding-right: 120rpx;
	}

	.tui-address-label {
		padding: 5rpx 8rpx;
		flex-shrink: 0;
		background: #e41f19;
		border-radius: 6rpx;
		color: #fff;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		font-size: 25rpx;
		line-height: 25rpx;
		transform: scale(0.8);
		transform-origin: center center;
		margin-right: 6rpx;
	}

	.tui-address-imgbox {
		width: 80rpx;
		height: 100rpx;
		position: absolute;
		display: flex;
		justify-content: center;
		align-items: center;
		right: 10rpx;
	}

	.tui-address-img {
		width: 36rpx;
		height: 36rpx;
	}

	.tui-address-new {
		width: 100%;
		position: fixed;
		left: 0;
		bottom: 0;
		z-index: 999;
		padding: 20rpx 25rpx 30rpx;
		box-sizing: border-box;
		background: #fafafa;
	}

	.tui-safe-area {
		padding-bottom: env(safe-area-inset-bottom);
	}
</style>
