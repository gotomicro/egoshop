<template>
	<view class="container">
		<!-- #ifdef MP || H5-->
		<view class="tui-edit-goods">
			<view>购物车共<text class="tui-goods-num">{{cartNum}}</text>件商品</view>
			<view>
				<tui-button type="gray" :plain="true" shape="circle" width="160rpx" height="60rpx" :size="24" @click="editGoods">{{isEdit?"完成":"编辑商品"}}</tui-button>
			</view>
		</view>
		<!-- #endif -->
		<checkbox-group>
			<view class="tui-cart-cell  tui-mtop" v-for="(item,index) in listData" :key="index">
				<tui-swipe-action :actions="actions" @click="handlerButton" :params="item">
					<template v-slot:content>
						<!-- 生效购物车 -->
						<view class="tui-goods-item" v-if="item.status === 1">
							<label class="tui-checkbox" @click="checkBox(index, item.id)">
								<checkbox color="#fff" :checked="cartIds[index]" ></checkbox>
							</label>
							<image :src="item.cover" class="tui-goods-img" />
							<view class="tui-goods-info">
								<view class="tui-goods-title">
									{{item.title}}
								</view>
								<view class="tui-goods-model">
									{{item.subTitle}}
								</view>
								<view class="tui-price-box">
									<view class="tui-goods-price">{{item.price}}</view>
									<view class="tui-scale">
										<tui-numberbox :value="item.num" :height="40" :width="74" :min="1" :index="index" disabled=true></tui-numberbox>
									</view>
								</view>
							</view>
						</view>
						<!-- 失效购物车 -->
						<view class="tui-goods-item" v-if="item.status === 2">
							<view class="tui-checkbox tui-invalid-pr">
								<view class="tui-invalid-text">失效</view>
							</view>
							<image :src="item.cover" class="tui-goods-img" />
							<view class="tui-goods-info">
								<view class="tui-goods-title tui-gray">
									{{item.title}}
								</view>
								<view class="tui-price-box tui-flex-center">
									<view class="tui-goods-invalid">产品失效</view>
									<view class="tui-btn-alike">
										<tui-button type="white" :plain="true" shape="circle" width="120rpx" height="48rpx" :size="24">找相似</tui-button>
									</view>
								</view>
							</view>
						</view>

					</template>
				</tui-swipe-action>
			</view>
		</checkbox-group>

		<!--tabbar-->
		<view class="tui-tabbar">
			<view class="tui-checkAll">
				<checkbox-group>
					<label class="tui-checkbox" @click="allCheckBox">
						<checkbox color="#fff" :checked="isAllCheck" ></checkbox>
						<text class="tui-checkbox-pl">全选({{isCheckedNumber}})</text>
					</label>
				</checkbox-group>
				<view class="tui-total-price" v-if="!isEdit">合计:<text class="tui-bold">￥{{allPrice}}</text> </view>
			</view>
			<view>
				<tui-button width="200rpx" height="70rpx" :size="30" type="danger" shape="circle" v-if="!isEdit" @click="orderDown">去结算({{isCheckedNumber}})</tui-button>
				<tui-button width="200rpx" height="70rpx" :size="30" type="danger" shape="circle" :plain="true" v-else>删除</tui-button>
			</view>
		</view>
		<!--加载loadding-->
		<tui-loadmore :visible="loadding" :index="3" type="red"></tui-loadmore>
	</view>
</template>

<script>
	import tuiSwipeAction from "@/components/swipe-action/swipe-action"
	import tuiButton from "@/components/extend/button/button"
	import tuiNumberbox from "@/components/numberbox/numberbox"
	import tuiIcon from "@/components/icon/icon"
	import tuiDivider from "@/components/divider/divider"
	import tuiLoadmore from "@/components/loadmore/loadmore"
	import {api} from "../../../config/api"

	export default {
		components: {
			tuiSwipeAction,
			tuiButton,
			tuiNumberbox,
			tuiIcon,
			tuiDivider,
			tuiLoadmore
		},
		data() {
			return {
				listData: [],
				totalPage: 0,
				cartNum: 0,
				cartIds: [],
				isAllCheck: false,
				dataList: [{
					id: 1,
					buyNum:2
				}, {
					id: 2,
					buyNum:1
				}],
				actions: [{
						name: '收藏',
						width: 64,
						color: '#fff',
						fontsize: 28,
						background: '#FFC600'
					},
					{
						name: '看相似',
						color: '#fff',
						fontsize: 28,
						width: 64,
						background: '#FF7035'
					},
					{
						name: '删除',
						color: '#fff',
						fontsize: 28,
						width: 64,
						background: '#F82400'
					}
				],
				actions2: [{
						name: '看相似',
						color: '#fff',
						fontsize: 28,
						width: 64,
						background: '#FF7035'
					},
					{
						name: '删除',
						color: '#fff',
						fontsize: 28,
						width: 64,
						background: '#F82400'
					}
				],
				isEdit: false,
				pageIndex: 1,
				loadding: false,
				pullUpOn: true
			}
		},
		onShow: function() {
			this.getListData()
			this.cartTotal()
		},
		methods: {
			async getListData() {
				const resp = await this.wxhttp.route(api.cart.list,{});
				if (resp.code === 0) {
					this.listData = resp.data.list;
					this.totalPage = resp.data.total;
				}
			},
			// 购物车数量
			async cartTotal() {
				const resp = await this.wxhttp.route(api.cart.totalNum,{});
				if (resp.code !== 0) {
					return false
				}
				this.cartNum = resp.data
			},

			//勾选添加按钮
			checkBox(index, id) {
				if (this.cartIds[index]) {
					this.$set(this.cartIds, index, false);
				} else {
					this.$set(this.cartIds, index, id);
				}
				console.log("this.cartIds",this.cartIds)
			},
			// 勾选全部按钮
			allCheckBox(feed) {
				// 先清空
				this.cartIds = [];
				if (this.isAllCheck) {
					this.isAllCheck = false;
				} else {
					this.isAllCheck = true;

					//循环遍历所有的商品id
					for (let i = 0; i < this.listData.length; i++) {
						const element = this.listData[i];
						if (element.status === 1) {
							this.cartIds.push(element.id);
						}
					}
				}
			},
			// 下订单
			async orderDown() {
				if (this.cartIds.length === 0) {
					this.tui.toast("请选择商品");
					return false;
				}
				// 去掉数组中空的false的
				var newGoodsIds = [];
				for (let i = 0; i < this.cartIds.length; i++) {
					const element = this.cartIds[i];
					if (element) {
						newGoodsIds.push(element);
					}
				}

				uni.navigateTo({
					url: '/pages/submitOrder/main?comMode=cart&cartIds=' + JSON.stringify(newGoodsIds)
				})
			},

			handlerButton: function(e) {
				let index = e.index;
				let item = e.item;
				this.tui.toast(`商品id：${item.id}，按钮index：${index}`);
			},
			editGoods: function() {
				// #ifdef H5 || MP
				this.isEdit = !this.isEdit;
				// #endif
			},
			detail: function() {
				uni.navigateTo({
					url: '/pages/productDetail/main'
				})
			},
			btnPay(){
				uni.navigateTo({
					url: '/pages/submitOrder/main'
				})
			}
		},

		computed: {
			isCheckedNumber() {
				let number = 0;
				for (let i = 0; i < this.cartIds.length; i++) {
					if (this.cartIds[i]) {
						number++;
					}
				}
				if (number === this.listData.length && number !== 0) {
					this.isAllCheck = true;
				} else {
					this.isAllCheck = false;
				}
				return number;
			},
			allPrice() {
				var calPrice = 0;
				for (let i = 0; i < this.cartIds.length; i++) {
					if (this.cartIds[i]) {
						calPrice = calPrice + this.listData[i].price * this.listData[i].num;
					}
				}
				return calPrice;
			}
		},
		//
		// onPullDownRefresh() {
		// 	setTimeout(() => {
		// 		uni.stopPullDownRefresh()
		// 	}, 200)
		// },
		// onPullDownRefresh: function() {
		// 	let loadData = JSON.parse(JSON.stringify(this.productList));
		// 	loadData = loadData.splice(0, 10)
		// 	this.productList = loadData;
		// 	this.pageIndex = 1;
		// 	this.pullUpOn = true;
		// 	this.loadding = false
		// 	uni.stopPullDownRefresh()
		// },
		// onReachBottom: function() {
		// 	if (!this.pullUpOn) return;
		// 	this.loadding = true;
		// 	if (this.pageIndex === 4) {
		// 		this.loadding = false;
		// 		this.pullUpOn = false
		// 	} else {
		// 		let loadData = JSON.parse(JSON.stringify(this.productList));
		// 		loadData = loadData.splice(0, 10)
		// 		if (this.pageIndex === 1) {
		// 			loadData = loadData.reverse();
		// 		}
		// 		this.productList = this.productList.concat(loadData);
		// 		this.pageIndex = this.pageIndex + 1;
		// 		this.loadding = false
		// 	}
		// },
		onNavigationBarButtonTap(e) {
			this.isEdit = !this.isEdit;
			let text = this.isEdit ? "完成" : "编辑";
			// #ifdef APP-PLUS
			let webView = this.$mp.page.$getAppWebview();
			webView.setTitleNViewButtonStyle(0, {
				text: text
			});
			// #endif
		}
	}

</script>

<style>
	.container {
		padding-bottom: 120rpx;
	}

	.tui-mtop {
		margin-top: 24rpx;
	}

	.tui-edit-goods {
		width: 100%;
		border-radius: 12rpx;
		overflow: hidden;
		padding: 24rpx 30rpx 0 30rpx;
		box-sizing: border-box;
		display: flex;
		justify-content: space-between;
		align-items: center;
		color: #333;
		font-size: 24rpx;
	}

	.tui-goods-num {
		font-weight: bold;
		color: #e41f19;
	}

	.tui-cart-cell {
		width: 100%;
		border-radius: 12rpx;
		background: #FFFFFF;
		padding: 40rpx 0;
		overflow: hidden;
	}

	.tui-goods-item {
		display: flex;
		padding: 0 30rpx;
		box-sizing: border-box;
	}

	.tui-checkbox {
		/* width: 40rpx; */
		padding-right: 30rpx;
		display: flex;
		align-items: center;
	}

	/* #ifdef APP-PLUS || MP */
	.tui-checkbox .wx-checkbox-input {
		width: 36rpx;
		height: 36rpx;
		border-radius: 50%;
		margin: 0;
	}

	.tui-checkbox .wx-checkbox-input.wx-checkbox-input-checked {
		background: #F82400;
		width: 40rpx;
		height: 40rpx;
		border: none;
	}

	/* #endif */
	/* #ifdef H5 */
	>>>.tui-checkbox .uni-checkbox-input {
		width: 36rpx;
		height: 36rpx;
		border-radius: 50% !important;
		margin: 0;
	}

	>>>.tui-checkbox .uni-checkbox-input.uni-checkbox-input-checked {
		background: #F82400;
		width: 40rpx;
		height: 40rpx;
		border: none;
	}

	/* #endif */
	.tui-goods-img {
		width: 220rpx;
		height: 220rpx !important;
		border-radius: 12rpx;
		flex-shrink: 0;
		display: block;
	}

	.tui-goods-info {
		width: 100%;
		padding-left: 20rpx;
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		justify-content: space-between;
		box-sizing: border-box;
		overflow: hidden;
	}

	.tui-goods-title {
		white-space: normal;
		word-break: break-all;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
		font-size: 24rpx;
		color: #333;
	}

	.tui-goods-model {
		max-width: 100%;
		color: #333;
		background: #F5F5F5;
		border-radius: 40rpx;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 16rpx;
		box-sizing: border-box;
	}

	.tui-model-text {
		max-width: 100%;
		transform: scale(0.9);
		transform-origin: 0 center;
		font-size: 24rpx;
		line-height: 32rpx;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.tui-price-box {
		width: 100%;
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
	}

	.tui-goods-price {
		font-size: 34rpx;
		font-weight: 500;
		color: #e41f19;
	}

	.tui-scale {
		transform: scale(0.8);
		transform-origin: 100% 100%;
	}

	.tui-activity {
		font-size: 24rpx;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 30rpx 20rpx 100rpx;
		box-sizing: border-box;
	}

	.tui-buy {
		display: flex;
		align-items: center
	}

	.tui-bold {
		font-weight: bold;
	}

	.tui-sub-info {
		max-width: 532rpx;
		font-size: 24rpx;
		line-height: 24rpx;
		padding: 20rpx 30rpx 10rpx 30rpx;
		box-sizing: border-box;
		color: #333;
		transform: scale(0.8);
		transform-origin: 100% center;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		margin-left: auto
	}

	.tui-invalid-text {
		width: 66rpx;
		margin-right: 4rpx;
		text-align: center;
		font-size: 24rpx;
		color: #fff;
		background: rgba(0, 0, 0, .3);
		transform: scale(0.8);
		transform-origin: center center;
		border-radius: 4rpx;
		flex-shrink: 0;
	}

	.tui-invalid-pr {
		padding-right: 0 !important;
	}

	.tui-gray {
		color: #B2B2B2 !important;
	}

	.tui-goods-invalid {
		color: #555;
		font-size: 24rpx;
	}

	.tui-flex-center {
		align-items: center !important;
	}

	.tui-invalid-ptop {
		padding-top: 40rpx;
	}

	.tui-tabbar {
		width: 100%;
		height: 100rpx;
		background: #fff;
		position: fixed;
		left: 0;
		bottom: 0;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 30rpx;
		box-sizing: border-box;
		font-size: 24rpx;
		z-index: 99999;
	}

	.tui-tabbar::before {
		content: '';
		width: 100%;
		border-top: 1rpx solid #d9d9d9;
		position: absolute;
		top: 0;
		left: 0;
		-webkit-transform: scaleY(0.5);
		transform: scaleY(0.5);
	}

	.tui-checkAll {
		display: flex;
		align-items: center;
	}

	.tui-checkbox-pl {
		padding-left: 12rpx;
	}

	.tui-total-price {
		font-size: 30rpx !important;
	}

	/*猜你喜欢*/
	.tui-youlike {
		padding-left: 12rpx
	}

	.tui-product-list {
		display: flex;
		justify-content: space-between;
		flex-direction: row;
		flex-wrap: wrap;
		box-sizing: border-box;
	}

	.tui-product-container {
		flex: 1;
		margin-right: 2%;
	}

	.tui-product-container:last-child {
		margin-right: 0;
	}

	.tui-pro-item {
		width: 100%;
		margin-bottom: 4%;
		background: #fff;
		box-sizing: border-box;
		border-radius: 12rpx;
		overflow: hidden;
	}

	.tui-pro-img {
		width: 100%;
		display: block;
	}

	.tui-pro-content {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		box-sizing: border-box;
		padding: 20rpx;
	}

	.tui-pro-tit {
		color: #2e2e2e;
		font-size: 26rpx;
		word-break: break-all;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
	}

	.tui-pro-price {
		padding-top: 18rpx;
	}

	.tui-sale-price {
		font-size: 34rpx;
		font-weight: 500;
		color: #e41f19;
	}

	.tui-factory-price {
		font-size: 24rpx;
		color: #a0a0a0;
		text-decoration: line-through;
		padding-left: 12rpx;
	}

	.tui-pro-pay {
		padding-top: 10rpx;
		font-size: 24rpx;
		color: #656565;
	}
</style>
