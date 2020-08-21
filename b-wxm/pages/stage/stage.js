// pages/stage/stage.js
Component({
  /**
   * 组件的属性列表
   */
  properties: {
    stage: {
      type: String
    },
    dict: {
      type: Object,
      value: {
        "new": "待支付",
        "paid": "已支付",
        "transport": "已发货",
        "finish": "已完成",
        "c_w_cancel": "买家申请退款",
        "b_cancel": "卖家取消",
        "c_cancel": "买家取消",
        "repaid": "卖家退款",
      }
    },
    content: {
      type: String
    }
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

  }
})
