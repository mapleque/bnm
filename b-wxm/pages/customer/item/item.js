// pages/customer/item/item.js
const http = require('../../../utils/http.js')
Page({

  /**
   * 页面的初始数据
   */
  data: {
    error: null,
    formData: {
      counts: 1,
      reciever: '',
      phone: '',
      address: '',
      additional:'',
    },
    item: {},
    rules: [
      {
        name: "counts",
        rules: { required: true, message: "商品数量不能为0" }
      },
      {
        name: "reciever",
        rules: { required: true, message: "收货人不能为空" }
      },
      {
        name: "address",
        rules: { required: true, message: "收货地址不能为空" }
      },
      {
        name: "phone",
        rules: { required: true, message: "联系方式不能为空" }
      },
    ],
    galleryShow: false,
    galleryImages: [],
  },
  itid: 0,
  bid: 0,
  loadItem: function () {
    http.get(`/customer/business/${this.bid}/items/${this.itid}`, {
      success: ({status, data, message}) => {
        if (status === 200) {
          this.setData({item: data})
        } else {
          wx.showToast({
            icon: 'error',
            title: message
          })
        }
      }
    })
  },
  hideGallery: function () {
    this.setData({
      galleryShow: false,
      galleryImages: [],
    })
  },
  showGallery: function (e) {
    this.setData({
      galleryShow: true,
      galleryImages: this.data.item.pic.split('$'),
    })
  },
  formInputChange: function (e) {
    const { field } = e.currentTarget.dataset
    this.setData({
      [`formData.${field}`]: e.detail.value
    })
  },
  submitForm: function () {
    this.selectComponent('#form').validate((valid, errors) => {
      console.log('valid', valid, errors)
      if (!valid) {
        const firstError = Object.keys(errors)
        if (firstError.length) {
          this.setData({
            error: errors[firstError[0]].message
          })
        }
      } else {
        const data = this.data.formData
        data.counts = data.counts*1
        http.post(`/customer/business/${this.bid}/items/${this.itid}`, data, {
          success: ({ status, message }) => {
            if (status === 200) {
              wx.showToast({
                type: 'success',
                title: '下单成功，请在订单列表中查看',
                mask: true
              })
            } else {
              this.setData({
                error: message
              })
            }
          },
        })
      }
    })
  },
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.itid = options.itid
    this.bid = options.bid
    this.loadItem()
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

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  }
})