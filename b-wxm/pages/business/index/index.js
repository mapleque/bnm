// pages/business/index/index.js
const http = require('../../../utils/http.js')
const image = require('../../../utils/image.js')

Page({

  /**
   * 页面的初始数据
   */
  data: {
    businessDefaultIcon: image.business,
    itemDefaultIcon: image.item,
    isItemList: true,
    isOrderList: false,
    loading: false,
    profile: {
      name: '商户名',
      desc: '商户简介',
    },
    noMoreItemData: false,
    noMoreOrderData: false,
    lastItemId: 0,
    lastOrderId: 0,
    itemList: [],
    orderList: [],
    actinos: {},
    galleryShow: false,
    galleryImages: [],
  },
  hideGallery: function () {
    this.setData({
      galleryShow: false,
      galleryImages: [],
    })
  },
  showGallery: function (e) {
    console.log(e)
    this.setData({
      galleryShow: true,
      galleryImages: e.currentTarget.dataset.pic.split('$'),
    })
  },
  toItemAdd: function () {
    wx.navigateTo({
      url: '/pages/business/item/add',
    })
  },
  toOrderList: function () {
    this.setData({ isOrderList: true, isItemList: false })
    if (this.data.orderList.length === 0) {
      this.loadOrderList()
    }
  },
  toItemList: function () {
    this.setData({ isOrderList: false, isItemList: true })
    if (this.data.itemList.length === 0) {
      this.loadItemList()
    }
  },
  actionClick: function (e) {
    console.log(e.detail.value)
    const { item, order, stage } = e.detail.value
    if (item) {
      this.changeItemStatus(item)
    } else {
      this.changeOrderStage(order, stage)
    }
  },
  changeOrderStage: function (order, stage) {
    order.stage = stage
    http.post('/business/orders/'+order.id, order, {
      success: ({status, message}) => {
        if (status !== 200) {
          wx.showToast({
            title: message,
            icon: 'error'
          })
        } else {
          const orderList = this.data.orderList
          for (var i = 0; i < itemList.length; i++) {
            if (orderList[i].id === order.id) {
              orderList[i].stage = order.stage
              break
            }
          }
          this.setData({
            orderList,
            actions: {},
          })
          wx.showToast({
            title: '操作成功',
            icon: 'success'
          })
        }
      }
    })
  },
  changeItemStatus: function (item) {
    item.status = item.status ? 0: 1
    http.post('/business/items/'+item.id, item, {
      success: ({status, message}) => {
        if (status !== 200) {
          wx.showToast({
            title: message,
            icon: 'error'
          })
        } else {
          const itemList = this.data.itemList
          for (var i = 0; i < itemList.length; i++) {
            if (itemList[i].id === item.id) {
              itemList[i].status = item.status
              break
            }
          }
          this.setData({
            itemList,
            actions: {},
          })
          wx.showToast({
            title: '操作成功',
            icon: 'success'
          })
        }
      }
    })
  },
  itemClick: function (e) {
    const item = e.currentTarget.dataset
    this.setData({
      actions: {
        show: true,
        title: "修改商品状态",
        groups: [{
          text: "上/下架",
          value: item,
        }],
      }
    })
  },
  orderClick: function (e) {
    const order = e.currentTarget.dataset.item
    console.log(order)
    const groups = []
    switch (order.stage) {
      case 'new':
        groups.push({
          text: '标记支付',
          value: {
            order,
            stage: 'paid'
          },
        })
        groups.push({
          text: '取消订单',
          value: {
            order,
            stage: 'b_cancel'
          },
        })
        break
      case 'c_w_cancel':
        gropus.push({
          text: '退款',
          value: {
            order,
            stage: 'repaid'
          }
        })
        break
      case 'paid':
        groups.push({
          text: '发货',
          value: {
            order,
            stage: 'transport'
          },
        })
        groups.push({
          text: '退款',
          value: {
            order,
            stage: 'repaid'
          },
        })
        break
    }
    this.setData({
      actions: {
        show: true,
        title: "修改订单状态",
        groups
      }
    })
  },
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    http.get('/business/profile', {
      success: ({ status, data, message }) => {
        this.setData({loading: false})
        if (status !== 200) {
          wx.showToast({ // 显示Toast
            title: message,
            icon: 'error'
          })
        } else {
          if (data.status === 0) {
            // 跳转到认证页
            wx.navigateTo({
              url: '../valid/valid'
            })
          } else {
            // 初始化商户首页
            this.setData({
              profile: data
            })
          }
        }
      }
    })
  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {
    if (this.data.loading) return
    const {
      isItemList,
      isOrderList,
      noMoreItemData,
      noMoreOrderData,
      lastItemId,
      lastOrderId
    } = this.data
    if (isOrderList && noMoreOrderData) return
    if (isItemList && noMoreItemData) return
    if (isOrderList) this.loadOrderList(lastOrderId)
    if (isItemList) this.loadItemList(lastItemId)
  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  },

  onShow: function () {
    this.setData({
      itemList: [],
      orderList: [],
    })
    if (this.data.isItemList) {
      this.loadItemList()
    }
    if (this.data.isOrderList) {
      this.loadOrderList()
    }
  },

  loadItemList: function (lastId) {
    if (this.data.loading) return
    this.setData({loading: true})
    if (!lastId) { lastId = 0 }
    if (this.data.itemList.length === 0 || lastId > 0) {
      http.get('/business/items?s=10&l='+lastId, {
        success: ({status, data, message}) => {
          this.setData({loading: false})
          if (status === 200) {
            if (data.length === 0) {
              this.setData({ noMoreItemData: true })
              return
            }
            const { itemList } = this.data
            this.setData({
              itemList: itemList.concat(data),
              lastItemId: data[data.length-1].id
            })
          } else {
            wx.showToast({
              icon: 'error',
              title: message
            })
          }
        }
      })
    }
  },
  loadOrderList: function (lastId) {
    if (this.data.loading) return
    this.setData({ loading: true })
    if (!lastId) { lastId = 0 }
    if (this.data.orderList.length === 0 || lastId > 0) {
      http.get('/business/orders?s=10&l=' + lastId, {
        success: ({ status, data, message }) => {
          this.setData({loading: false})
          if (status === 200) {
            if (data.length === 0) {
              this.setData({
                noMoreOrderData: true
              })
              return
            }
            const { orderList } = this.data
            this.setData({
              orderList: orderList.concat(data),
              lastOrderId: data[data.length-1].id
            })
          } else {
            wx.showToast({
              icon: 'error',
              title: message
            })
          }
        }
      })
    }
  },
})