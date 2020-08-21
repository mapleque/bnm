// pages/business/valid/valid.js
const http = require('../../../utils/http.js')
Page({

  /**
   * 页面的初始数据
   */
  data: {
    error: null,
    avatarSizeType: ['compressed'],
    formData: {},
    rules: [
      {
        name:"name",
        rules: { required: true, message: "商户名不能为空" }
      },
      {
        name: "desc",
        rules: { required: true, message: "商户简介不能为空" }
      },
      {
        name: "wxid",
        rules: { required: true, message: "商户主微信ID不能为空" }
      },
    ],
    files: [],
    uploadFile: (files) => {
      const urls = files.contents.map((arrBuf, i) => {
        const filename = files.tempFilePaths[i]
        const ext = filename.substr(filename.lastIndexOf('.') + 1)
        const base64 = wx.arrayBufferToBase64(arrBuf)
        return `data:image/${ext};base64,${base64}`
      })
      // 文件上传的函数，返回一个promise
      return new Promise((resolve, reject) => {
        resolve({ urls })
      })
    },
  },
  deleteFile(e) {
    const { files } = this.data
    const { index } = e.detail
    files.shift(index, 1)
    this.setData({ files })
  },
  uploadError(e) {
    console.log('upload error', e.detail)
  },
  uploadSuccess(e) {
    const { files } = this.data
    this.setData({ files: files.concat(e.detail.urls) })
  },
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    http.get('/business/profile', {
      success: ({ status, data, message }) => {
        if (status === 200) {
            // 初始化商户首页
            this.setData({
              formData: data
            })
        }
      }
    })
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

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

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
        if (this.data.files.length > 0) {
          data.avatar = this.data.files[0]
        }
        http.post('/business/profile', data, {
          success: ({status, message}) => {
            if (status === 200) {
              wx.showToast({
                type: 'success',
                title: '审核信息已提交，请等待管理员审核'
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
  }
})