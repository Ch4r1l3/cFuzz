import request from '@/utils/request'

export function getItemsCombine(offset, limit, name) {
  return request({
    url: `/api/storage_item?offset=${offset}&limit=${limit}&name=${name}`,
    method: 'get'
  })
}

export function getItemsByType(type) {
  return request({
    url: `/api/storage_item/${type}`,
    method: 'get'
  })
}

export function getItemsByTypeCombine(type, offset, limit, name) {
  return request({
    url: `/api/storage_item/${type}?offset=${offset}&limit=${limit}&name=${name}`,
    method: 'get'
  })
}

export function getItem(id) {
  return request({
    url: `/api/storage_item/${id}`,
    method: 'get'
  })
}

export function createItem(item) {
  return request({
    url: '/api/storage_item/exist',
    method: 'post',
    data: item
  })
}

export function editItem(item) {
  return request({
    url: `/api/storage_item/${item.id}`,
    method: 'put',
    data: item
  })
}

export function deleteItem(item) {
  return request({
    url: `/api/storage_item/${item.id}`,
    method: 'delete'
  })
}
