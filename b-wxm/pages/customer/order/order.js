// pages/customer/order/order.js
const http = require('../../../utils/http.js')
const image = require('../../../utils/image.js')
Page({

  /**
   * 页面的初始数据
   */
  data: {
    itemDefaultIcon: image.item,
    loading: false,
    lastId: 0,
    noMoreData: false,
    orderList: [],
    actinos: {}
  },
  actionClick: function (e) {
    const { order, stage } = e.detail.value
    this.changeOrderStage(order, stage)
  },
  changeOrderStage: function (order, stage) {
    order.stage = stage
    http.post('/customer/orders/' + order.id, order, {
      success: ({ status, message }) => {
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
  orderClick: function (e) {
    const order = e.currentTarget.dataset.item
    const groups = []
    switch (order.stage) {
      case 'new':
        groups.push({
          text: '取消订单',
          value: {
            order,
            stage: 'c_cancel'
          }
        })
        break
      case 'transport':
        groups.push({
          text: '确认收货',
          value: {
            order,
            stage: 'paid'
          },
        })
        break
      case 'paid':
        groups.push({
          text: '申请退款',
          value: {
            order,
            stage: 'c_w_cancel'
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
  loadOrderList: function (lastId) {
    this.setData({loading: true})
    if (!lastId) { lastId = 0 }
    if (this.data.orderList.length === 0 || lastId > 0) {
      http.get('/customer/orders?s=10&l=' + lastId, {
        success: ({ status, data, message }) => {
          this.setData({ loading: false })
          if (status === 200) {
            if (data.length === 0) {
              this.setData({ noMoreData: true })
              return
            }
            const { orderList, lastId } = this.data
            this.setData({
              orderList: orderList.concat(data.map(item => {
                const ts = new Date(item.create_at)
                item.create_at = `${ts.getFullYear()}${ts.getMonth()+1}${ts.getDate()}${ts.getHours()}${ts.getMinutes()}${ts.getSeconds()}`
                return item
              })),
              lastId: data[data.length-1].id
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
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.loadOrderList()
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {
    if (this.data.loading || this.data.noMoreData) return
    this.loadOrderList(this.data.lastId)
  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  }
})