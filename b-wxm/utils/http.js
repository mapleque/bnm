const host = 'http://192.168.1.5:8080'
// const host = 'http://localhost:8080'
const version = '1.0.0'

// 当多个请求同时需要登录的时候，
// 需要保证登陆api只调用一次。
//
// 这里用一个锁和队列来控制，
// 首先，所有callback都先进队列。
// 当登陆api正在被调用的时候，直接返回不处理。
// 当登陆api成功返回后，会逐个执行队列中的callback。
let loginLock = false
const loginCallbackQueue = []
const login = (success) => {
  loginCallbackQueue.push(success)
  if (loginLock) { return }
  loginLock = true

  // 为了避免login接口有问题时陷入死循环
  // 这里记录一下login的请求时间，一般来说10秒之内不需要请求多次
  const lastTime = wx.getStorageSync('last_login_time')
  const now = new Date()
  if (now - lastTime < 10000) {
    loginLock = false
    return
  }
  wx.setStorageSync('last_login_time', now)
  // 登录
  wx.login({
    success: res => {
      // 发送 res.code 到后台换取 openId, sessionKey, unionId
      post('/login', { code: res.code }, {
        success: ({ status, data, message }) => {
          if (status === 200) {
            wx.setStorageSync('token', data)
            while (loginCallbackQueue.length > 0) {
              const cb = loginCallbackQueue.shift()
              cb()
            }
          } else {
            // tip login fail
            wx.showToast({ // 显示Toast
              title: message,
              icon: 'error',
              duration: 1500
            })
          }
          loginLock = false
        }
      })
    }
  })
}

// 这里封装了针对401返回码的自动登录逻辑
const request = (uri, method, data, cb) => {
  wx.request({
    url: getUrl(uri),
    data,
    method,
    header: getHeader(),
    success: res => {
      if (res.statusCode == 200) {
        if (res.data.status === 401) {
          login(() => { request(uri, method, data, cb) })
          return
        }
        typeof cb.success === 'function' && cb.success(res.data)
      }
    },
    fail: cb.fail,
  })
}

const post = (uri, data, cb) => {
  request(uri,'POST', data, cb)
}

const get = (uri, cb) => {
  request(uri, 'GET', null, cb)
}

const del = (uri, cb) => {
  request(uri, 'DELETE', null, cb)
}

const getUrl = uri => {
  let url = host+uri
  if (uri.indexOf('?') >= 0) {
    url+='&'
  } else {
    url+='?'
  }
  url+='v='+version
  return url
}

const getHeader = () => {
  const header = {}
  header['Content-Type'] = 'application/json'
  header['Session-Client'] = 'bnm-b-wx'
  header['Session-Version'] = version
  const token = wx.getStorageSync('token')
  if (token) {
    header['Session-Key'] = token
  }
  return header
}

module.exports = {
  post: post,
  get: get,
  del: del,
}