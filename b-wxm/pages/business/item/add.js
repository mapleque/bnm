// pages/business/item/add.js
const http = require('../../../utils/http.js')
Page({

  /**
   * 页面的初始数据
   */
  data: {
    error: null,
    formData: {},
    rules: [
      {
        name: "name",
        rules: { required: true, message: "商品名不能为空" }
      },
      {
        name: "desc",
        rules: { required: true, message: "商品简介不能为空" }
      },
      {
        name: "price",
        rules: { required: true, message: "商品价格不能为空" }
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
      console.log(urls)
      // 文件上传的函数，返回一个promise
      return new Promise((resolve, reject) => {
        resolve({ urls })
      })
    },
  },
  deleteFile(e) {
    console.log(this.data.files)
    const { files } = this.data
    const { index } = e.detail
    files.shift(index, 1)
    this.setData({files})
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
        data.price = data.price*100
        data.pic = this.data.files.join('$')
        http.post('/business/items', data, {
          success: ({ status, message }) => {
            if (status === 200) {
              wx.showToast({
                type: 'success',
                title: '商品添加成功，请在商品列表中查看',
                mask: true,
                complete: () => {
                  wx.navigateBack({})
                }
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