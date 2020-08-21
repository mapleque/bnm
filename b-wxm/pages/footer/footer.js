// pages/customer/footer/footer.js
const app = getApp()
Component({
  /**
   * 组件的属性列表
   */
  properties: {

  },

  /**
   * 组件的初始数据
   */
  data: {
  },

  /**
   * 组件的方法列表
   */
  methods: {
    toAddress: function () {
      wx.navigateTo({
        url: '/pages/customer/address/address',
      })
    },
    toOrder: function () {
      wx.navigateTo({
        url: '/pages/customer/order/order',
      })
    }
  },

  /**
   * 组件的生命周期
   */
  // 在组件实例进入页面节点树时执行
  attached: function () {
  },
  // 在组件实例被从页面节点树移除时执行
  detached: function () {
  },
})
