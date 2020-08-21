// pages/customer/address/address.js
const http = require('../../../utils/http.js')
Page({

  /**
   * 页面的初始数据
   */
  data: {
    loading: false,
    lastId: 0,
    addressList: [],
    noMoreData: false,
    actions: {}
  },
  actionClick: function (e) {
    const { address } = e.detail.value
    this.deleteAddress(address.id)
  },
  itemClick: function (e) {
    const address = e.currentTarget.dataset.item
    this.setData({
      actions: {
        show: true,
        title: '请选择相关操作',
        groups: [{
          text: '删除',
          value: {
            address
          }
        }]
      }
    })
  },
  deleteAddress: function (id) {
    http.del(`/customer/address/${id}`, {
      success: ({status, message}) => {
        if (status === 200) {
          this.setData({
            addressList: this.data.addressList.filter(item => item.id !== id),
            actions: {}
          })
        } else {
          wx.showToast({
            title: message,
            icon: 'error'
          })
        }
      }
    })
  },
  loadAddressList: function (lastId) {
    this.setData({loading: true})
    if (!lastId) lastId = 0
    http.get(`/customer/address?s=10&l=${lastId}`, {
      success: ({status, data, message}) => {
        this.setData({loading: false})
        if (status === 200) {
          if (data.length === 0) {
            this.setData({noMoreData: true})
            return
          } else {
            this.setData({
              addressList: this.data.addressList.concat(data),
              lastId: data[data.length-1].id
            })
          }
        } else {
          wx.showToast({
            title: message,
            icon: 'error'
          })
        }
      }
    })
  },
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.loadAddressList()
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
    this.loadAddressList(this.data.lastId)
  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  }
})