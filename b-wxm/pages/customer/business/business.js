// pages/customer/business/business.js
const http = require('../../../utils/http.js')
const image = require('../../../utils/image.js')
Page({

  /**
   * 页面的初始数据
   */
  data: {
    businessDefaultIcon: image.business,
    itemDefaultIcon: image.item,
    profile: {},
    loading: false,
    noMoreData: false,
    lastId: 0,
    itemList: [],
    galleryShow: false,
    galleryImages: [],
  },
  bid: 0,
  itemClick: function (e) {
    const item = e.currentTarget.dataset.item
    wx.navigateTo({
      url: `/pages/customer/item/item?bid=${this.bid}&itid=${item.id}`,
    })
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
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.bid = options.id
    this.loadBusinessPorfile()
    this.loadItemList()
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
    console.log('page show')
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
    this.loadItemList()
  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  },
  loadBusinessPorfile: function () {
    http.get(`/customer/business/${this.bid}`, {
      success: ({status, data, message}) => {
        if (status === 200) {
          this.setData({ profile: data })
        } else {
          wx.showToast({
            icon: 'error',
            title: message
          })
        }
      }
    })
  },
  loadItemList: function (lastId) {
    this.setData({loading: true})
    if (!lastId) { lastId = 0 }
    if (this.data.itemList.length === 0 || lastId > 0) {
      http.get(`/customer/business/${this.bid}/items?s=10&l=${lastId}`, {
        success: ({ status, data, message }) => {
          this.setData({loading: false})
          if (status === 200) {
            if (data.length === 0) {
              this.setData({
                noMoreData: true
              })
              return
            }
            const { itemList, lastId } = this.data
            this.setData({
              itemList: itemList.concat(data),
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
})