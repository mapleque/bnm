//index.js
//获取应用实例
const app = getApp()
const http = require('../../utils/http.js')
const image = require('../../utils/image.js')

Page({
  data: {
    businessDefaultIcon: image.business,
    loading: false,
    noMoreData: false,
    lastId: 0,
    profile: null,
    userInfo: null,
    hasUserInfo: false,
    businessList: []
  },
  myBusiness: function() {
    wx.navigateTo({
      url: '../business/index/index'
    })
  },
  validBusiness: function() {
    wx.navigateTo({
      url: '../business/valid/valid'
    })
  },
  toBusiness: function (e) {
    const item = e.currentTarget.dataset.item
    wx.navigateTo({
      url: '/pages/customer/business/business?id=' + item.id,
    })
  },
  loadProfile: function () {
    http.get('/business/profile', {
      success: ({ status, data, message }) => {
        if (status === 200 && data.status === 1) {
          this.setData({
            profile: data
          })
        }
      }
    })
  },
  loadBusinessList: function (lastId) {
    this.setData({loading: true})
    if (!lastId) { lastId = 0 }
    if (this.data.businessList.length === 0 || lastId > 0) {
      http.get('/customer/business?s=10&l=' + lastId, {
        success: ({ status, data, message }) => {
          this.setData({ loading: false })
          if (status === 200) {
            if (data.length == 0) {
              this.setData({
                noMoreData: true
              })
              return
            }
            const { businessList, lastId } = this.data
            this.setData({
              businessList: businessList.concat(data),
              lastId: data[data.length-1].id,
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
  onLoad: function () {
    this.loadBusinessList()
    this.loadProfile()
    if (app.globalData.userInfo) {
      this.setData({
        userInfo: app.globalData.userInfo,
        hasUserInfo: true
      })
    } else {
      // 由于 getUserInfo 是网络请求，可能会在 Page.onLoad 之后才返回
      // 所以此处加入 callback 以防止这种情况
      app.userInfoReadyCallback = res => {
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        })
      }
    }
  },
  onReachBottom: function () {
    if (this.data.loading || this.data.noMoreData) { return }
    this.loadBusinessList(this.data.lastId)
  },
})
