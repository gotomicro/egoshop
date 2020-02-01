<template>
	<view class="container">
		<tui-tabs :tabs="tabs" :isFixed="scrollTop>=0" :currentTab="currentTab" selectedColor="#E41F19" sliderBgColor="#E41F19"
		 @change="change"></tui-tabs>
		<!--选项卡逻辑自己实现即可，此处未做处理-->
		<view :class="{'tui-order-list':scrollTop>=0}" v-if="listData.length > 0">
			<view class="tui-order-item" v-for="(item,orderIndex) in listData" :key="orderIndex">
				<tui-list-cell :hover="false" :lineLeft="false">
					<view class="tui-goods-title">
						<view>订单号：{{item.order.sn}}</view>
						<view class="tui-order-status">{{item.stateDesc}}</view>
					</view>
				</tui-list-cell>
				<block v-for="(item2,index2) in item.orderGoods" :key="index2">
					<tui-list-cell padding="0" @click="detail">
						<view class="tui-goods-item">
							<image :src="item2.cover" class="tui-goods-img"></image>
							<view class="tui-goods-center">
								<view class="tui-goods-name">{{item2.title}}</view>
								<view class="tui-goods-attr">{{item2.subTitle}}</view>
							</view>
							<view class="tui-price-right">
								<view>￥{{item2.price}}</view>
								<view>x{{item2.num}}</view>
							</view>
						</view>
					</tui-list-cell>
				</block>
				<tui-list-cell :hover="false" :last="true">
					<view class="tui-goods-price">
						<view>共{{item.order.goodsNum}}件商品 合计：</view>
						<view class="tui-size-24">￥</view>
						<view class="tui-price-large">{{item.order.amount}}</view>
					</view>
				</tui-list-cell>
				<view class="tui-order-btn">
					<view class="tui-btn-ml"  v-if="item.isCancel">
						<tui-button type="black" :plain="true" width="148rpx" height="56rpx" :size="26" shape="circle" @click="bindDeleteOrder(item.order.id)">删除订单</tui-button>
					</view>
					<view class="tui-btn-ml" v-if="item.isEvaluate">
						<tui-button type="black" :plain="true" width="148rpx" height="56rpx" :size="26" shape="circle">评价晒单</tui-button>
					</view>
					<view class="tui-btn-ml" v-if="item.isPay">
						<tui-button type="danger" :plain="true" width="148rpx" height="56rpx" :size="26" shape="circle">再次购买</tui-button>
					</view>
					<view class="tui-btn-ml" v-else>
						<tui-button type="danger" :plain="true" width="148rpx" height="56rpx" :size="26" shape="circle">确认支付</tui-button>
					</view>
				</view>
			</view>
		</view>
		<view :class="{'tui-order-list':scrollTop>=0}" v-else>
			<tui-tips :fixed="false" imgUrl="/static/images/toast/img_noorder_3x.png">暂无订单</tui-tips>
		</view>

		<!--加载loadding-->
		<tui-loadmore :visible="loadding" :index="3" type="red"></tui-loadmore>
		<tui-nomore :visible="!pullUpOn" bgcolor="#fafafa"></tui-nomore>
		<!--加载loadding-->

	</view>
</template>

<script>
	import tuiTabs from "@/components/tui-tabs/tui-tabs"
	import tuiButton from "@/components/extend/button/button"
	import tuiLoadmore from "@/components/loadmore/loadmore"
	import tuiNomore from "@/components/nomore/nomore"
	import tuiListCell from "@/components/list-cell/list-cell"
	import tuiSticky from "@/components/sticky/sticky"
	import tuiTips from "@/components/extend/tips/tips"
	import {api} from "../../../config/api"

	export default {
		components: {
			tuiTabs,
			tuiButton,
			tuiLoadmore,
			tuiNomore,
			tuiListCell,
			tuiSticky,
			tuiTips,
		},
		data() {
			return {
				tabs: [{
					id: 'all',
					name: "全部"
				}, {
					id: 'stateNew',
					name: "待付款"
				}, {
					id: 'statePay',
					name: "待发货"
				}, {
					id: 'stateSend',
					name: "待收货"
				}, {
					id: 'stateSuccess',
					name: "待评价"
				}],
				currentTab: 0,
				pageIndex: 1,
				loadding: false,
				pullUpOn: true,
				scrollTop: 0,
				listData: [],
			}
		},
		onShow() {
			this.getList()
		},
		methods: {
			change(e) {
				this.currentTab = e.index
				this.getList()
			},
			detail() {
				uni.navigateTo({
					url: '../orderDetail/orderDetail'
				})
			},
			// 获取订单列表
			async getList() {
				const resp = await this.wxhttp.route(api.order.list, {"stateType":this.tabs[this.currentTab].id});
				if (resp.code !== 0) {
					uni.showToast({
						title: resp.msg,
						icon: "error",
						duration: 1500
					});
					return false
				}
				this.listData = resp.data.list;
			},
			async bindDeleteOrder(id) {
				const resp = await this.wxhttp.route(api.order.delete, {"id":id});
				if (resp.code !== 0) {
					uni.showToast({
						title: resp.msg,
						icon: "error",
						duration: 1500
					});
					return false
				}
				this.tui.toast("删除成功")
				this.getList()
			}
		},
		onPullDownRefresh() {
			setTimeout(() => {
				uni.stopPullDownRefresh()
			}, 200);
		},
		onReachBottom() {
			//只是测试效果，逻辑以实际数据为准
			this.loadding = true
			this.pullUpOn = true
			setTimeout(() => {
				this.loadding = false
				this.pullUpOn = false
			}, 1000)
		},
		onPageScroll(e) {
			this.scrollTop = e.scrollTop;
		}
	}
</script>

<style>
	.container {
		padding-bottom: env(safe-area-inset-bottom);
	}

	.tui-order-list {
		margin-top: 80rpx;
	}

	.tui-order-item {
		margin-top: 20rpx;
		border-radius: 10rpx;
		overflow: hidden;
	}

	.tui-goods-title {
		width: 100%;
		font-size: 28rpx;
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.tui-order-status {
		color: #888;
		font-size: 26rpx;
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

	.tui-color-red {
		color: #E41F19;
		padding-right: 30rpx;
	}

	.tui-goods-price {
		width: 100%;
		display: flex;
		align-items: flex-end;
		justify-content: flex-end;
		font-size: 24rpx;
	}

	.tui-size-24 {
		font-size: 24rpx;
		line-height: 24rpx;
	}

	.tui-price-large {
		font-size: 32rpx;
		line-height: 30rpx;
		font-weight: 500;
	}

	.tui-order-btn {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: flex-end;
		background: #fff;
		padding: 10rpx 30rpx 20rpx;
		box-sizing: border-box;
	}

	.tui-btn-ml {
		margin-left: 20rpx;
	}
</style>
