import request from '@/utils/request'

export function getItemsCombine(offset, limit, name) {
  return request({
    url: `/api/image?offset=${offset}&limit=${limit}&name=${name}`,
    method: 'get'
  })
}

export function getItem(id) {
  return request({
    url: `/api/image/${id}`,
    method: 'get'
  })
}

export function createItem(item) {
  return request({
    url: '/api/image',
    method: 'post',
    data: item
  })
}

export function editItem(item) {
  return request({
    url: `/api/image/${item.id}`,
    method: 'put',
    data: item
  })
}

export function deleteItem(item) {
  return request({
    url: `/api/image/${item.id}`,
    method: 'delete'
  })
}
