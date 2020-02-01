<template>
	<view class="tui-container">
		<!--选项卡逻辑自己实现即可，此处未做处理-->
		<view :class="{'tui-order-list':scrollTop>=0}" v-if="listData.length > 0">
			<view class="tui-news-view">
				<block v-for="(item,index) in listData" :key="index">
					<tui-list-cell :index="index" @click="detail(item.goodsId,item.typeId)">
						<view class="tui-news-flex tui-flex-start">
							<view class="tui-news-picbox tui-w220 tui-h165">
								<image :src="item.cover" mode="widthFix" class="tui-block"></image>
							</view>
							<view class="tui-news-tbox tui-flex-column tui-flex-between tui-h165 tui-pl-20">
								<view class="tui-news-title">{{item.name}}</view>
								<view class="tui-sub-box">
									<view class="tui-sub-source">{{item.createdAt}}</view>
									<view class="tui-sub-cmt">
<!--										<view>{{item.cmtsNum}}评论</view>-->
										<view class="tui-scale">
											<tui-tag size="small" :plain="true" shape="circleRight">{{item.typeName}}</tui-tag>
										</view>
									</view>
								</view>
							</view>
						</view>
					</tui-list-cell>
				</block>
			</view>
		</view>
		<view :class="{'tui-order-list':scrollTop>=0}" v-else>
			<tui-tips :fixed="false" imgUrl="/static/images/toast/img_noorder_3x.png">暂无订单</tui-tips>
		</view>

		<!--加载loadding-->
		<tui-loadmore :visible="loading" :index="3" type="red"></tui-loadmore>
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
	import tuiTag from "@/components/tag/tag"
	import {api} from "../../../config/api"
	import {relativeTime} from "../../utils/common"

	export default {
		components: {
			tuiTabs,
			tuiButton,
			tuiLoadmore,
			tuiNomore,
			tuiListCell,
			tuiSticky,
			tuiTips,
			tuiTag,
		},
		data() {
			return {
				pageIndex: 1,
				scrollTop: 0,
				listData: [],
				tid: 2, // 表示收藏类型

				loading: false,
				pullUpOn: true,
				currentPage: 1, // 当前页面
				totalPage: 0, // 总页面
			}
		},
		onShow() {
			this.getListData(true)
		},
		methods: {
			// 获取收藏列表
			async getListData(first) {
				if (this.tid === undefined) {
					return
				}
				const resp = await this.wxhttp.route(api.userGoods.list, {}, {tid:parseInt(this.tid),typeId:0,currentPage:this.currentPage});
				if (resp.code !== 0) {
					uni.showToast({
						title: resp.msg,
						icon: "error",
						duration: 1500
					});
					return false
				}

				this.totalPage = resp.data.total;
				if (first) {
					this.listData = resp.data.list.map(function (item, index) {
						item.createdAt = relativeTime(item.createdAt)
						return item;
					});
				} else {
					//上拉加载跟多
					this.listData = this.listData.concat(resp.data.list.map(function (item, index) {
						item.createdAt = relativeTime(item.createdAt)
						return item;
					}));
				}
			},

			change(e) {
				this.currentTab = e.index
				this.getListData(true)
			},
			detail(id,typeId) {
				uni.navigateTo({
					url: '/pages/cominfo/main?id='+id
				})
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
		onPullDownRefresh: function() {
			this.pullUpOn = true;
			this.loading = false

			this.currentPage = 1;
			this.getListData(true);
			//刷新完成后关闭
			uni.stopPullDownRefresh()
		},

		onReachBottom: function() {
			if (!this.pullUpOn) return;
			this.loading = true;

			// 如果页面已经到达最大页面
			const tmpPage = this.currentPage + 1;
			if (tmpPage > this.totalPage) {
				this.loadding = false
				return false;
			}
			this.currentPage = tmpPage
			this.getListData();
			this.loading = false
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

	page {
		background: #f2f2f2;
	}


	.tui-order-list {
		margin-top: 80rpx;
	}

	page {
		background: #f2f2f2;
	}

	.tui-container {
		display: flex;
		flex-direction: column;
		box-sizing: border-box;
		padding-bottom: env(safe-area-inset-bottom);
	}

	.tui-searchbox {
		padding: 16rpx 20rpx;
		box-sizing: border-box;
	}

	.tui-rolling-search {
		width: 100%;
		height: 60rpx;
		border-radius: 35rpx;
		padding: 0 40rpx 0 30rpx;
		box-sizing: border-box;
		background: #fff;
		display: flex;
		align-items: center;
		flex-wrap: nowrap;
		color: #999;
	}

	.tui-swiper {
		font-size: 26rpx;
		height: 60rpx;
		flex: 1;
		padding-left: 12rpx;
	}

	.tui-swiper-item {
		display: flex;
		align-items: center;
	}

	.tui-hot-item {
		line-height: 26rpx;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.tui-banner-swiper {
		width: 100%;
		height: 300rpx;
		position: relative;
	}

	.tui-slide-image {
		width: 100%;
		height: 300rpx;
		display: block;
	}

	.tui-banner-title {
		width: 100%;
		height: 100rpx;
		position: absolute;
		z-index: 9999;
		color: #fff;
		bottom: 0;
		padding: 0 30rpx;
		padding-top: 25rpx;
		font-size: 34rpx;
		font-weight: bold;
		background: linear-gradient(to bottom, transparent, rgba(0, 0, 0, 0.7));
		box-sizing: border-box;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* #ifdef APP-PLUS || MP */
	.tui-banner-swiper .wx-swiper-dots.wx-swiper-dots-horizontal {
		width: 100%;
		top: 280rpx;
		text-align: right;
		padding-right: 30rpx;
		box-sizing: border-box;
	}

	.tui-banner-swiper .wx-swiper-dot {
		width: 8rpx;
		height: 8rpx;
		display: inline-flex;
		background: none;
		justify-content: space-between;
	}

	.tui-banner-swiper .wx-swiper-dot::before {
		content: '';
		flex-grow: 1;
		background: rgba(255, 255, 255, 0.9);
		border-radius: 8rpx;
	}

	.tui-banner-swiper .wx-swiper-dot-active::before {
		background: #5677fc;
	}

	.tui-banner-swiper .wx-swiper-dot.wx-swiper-dot-active {
		width: 18rpx;
	}

	/* #endif */

	/* #ifdef H5 */
	>>>.tui-banner-swiper .uni-swiper-dots.uni-swiper-dots-horizontal {
		width: 100%;
		top: 280rpx;
		text-align: right;
		padding-right: 30rpx;
		box-sizing: border-box;
	}

	>>>.tui-banner-swiper .uni-swiper-dot {
		width: 8rpx;
		height: 8rpx;
		display: inline-flex;
		background: none;
		justify-content: space-between;
	}

	>>>.tui-banner-swiper .uni-swiper-dot::before {
		content: '';
		flex-grow: 1;
		background: rgba(255, 255, 255, 0.9);
		border-radius: 8rpx;
	}

	>>>.tui-banner-swiper .uni-swiper-dot-active::before {
		background: #5677fc;
	}

	>>>.tui-banner-swiper .uni-swiper-dot.uni-swiper-dot-active {
		width: 18rpx;
	}

	/* #endif */

	.tui-news-flex {
		width: 100%;
		display: flex;
	}

	.tui-flex-start {
		align-items: flex-start !important;
	}

	.tui-flex-center {
		align-items: center !important;
	}

	.tui-flex-column {
		flex-direction: column !important;
	}

	.tui-flex-between {
		justify-content: space-between !important;
	}

	.tui-news-cell {
		display: flex;
		padding: 20rpx 30rpx;
	}

	.tui-news-tbox {
		flex: 1;
		width: 100%;
		box-sizing: border-box;
		display: flex;
	}

	.tui-news-picbox {
		display: flex;
		position: relative;
	}

	.tui-w220 {
		width: 220rpx;
	}

	.tui-h165 {
		height: 165rpx;
	}

	.tui-block {
		display: block;
	}

	.tui-w-full {
		width: 100%;
	}

	.tui-one-third {
		width: 33%;
	}

	.tui-news-title {
		width: 100%;
		font-size: 34rpx;
		word-break: break-all;
		word-wrap: break-word;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
		box-sizing: border-box;
	}

	.tui-pl-20 {
		padding-left: 20rpx;
	}

	.tui-pt20 {
		padding-top: 20rpx;
	}

	.tui-sub-box {
		display: flex;
		align-items: center;
		justify-content: space-between;
		color: #999;
		box-sizing: border-box;
		line-height: 24rpx;
	}

	.tui-sub-source {
		font-size: 26rpx;
	}

	.tui-sub-cmt {
		font-size: 24rpx;
		line-height: 24rpx;
		display: flex;
		align-items: center;
	}

	.tui-tag {
		padding: 2rpx 6rpx !important;
		margin-left: 10rpx;
	}

	.tui-scale {
		transform: scale(0.6);
		transform-origin: center center;
	}

	.tui-btm-badge {
		position: absolute;
		right: 0;
		bottom: 0;
		font-size: 24rpx;
		color: #fff;
		padding: 2rpx 12rpx;
		background: rgba(0, 0, 0, 0.6);
		z-index: 20;
		transform: scale(0.8);
		transform-origin: 100% 100%;
	}

	.tui-video {
		position: absolute;
		z-index: 10;
		display: flex;
		align-items: center;
		justify-content: center;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		transform-origin: 0 0;
	}

	.tui-icon {
		background: rgba(0, 0, 0, 0.5);
		border-radius: 50%;
		padding: 26rpx;
	}

	.tui-icon-box .tui-icon {
		background: none;
		padding: 0;
		border-radius: 0;
	}
</style>
