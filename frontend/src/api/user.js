import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/api/user/login',
    method: 'post',
    data: data
  })
}

export function getInfo() {
  return request({
    url: '/api/user/info',
    method: 'get'
  })
}

export function getItems() {
  return request({
    url: '/api/user',
    method: 'get'
  })
}

export function getItemsCombine(offset, limit, name) {
  return request({
    url: `/api/user?offset=${offset}&limit=${limit}&name=${name}`,
    method: 'get'
  })
}

export function getItem(id) {
  return request({
    url: `/api/user/${id}`,
    method: 'get'
  })
}

export function createItem(item) {
  return request({
    url: '/api/user',
    method: 'post',
    data: item
  })
}

export function adminEditItem(id, password) {
  return request({
    url: `/api/user/${id}`,
    method: 'put',
    data: {
      newPassword: password
    }
  })
}

export function editItem(id, oldPassword, newPassword) {
  return request({
    url: `/api/user/${id}`,
    method: 'put',
    data: {
      oldPassword: oldPassword,
      newPassword: newPassword
    }
  })
}

export function deleteItem(item) {
  return request({
    url: `/api/user/${item.id}`,
    method: 'delete'
  })
}
