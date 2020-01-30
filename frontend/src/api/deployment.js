import request from '@/utils/request'

export function getSimpList() {
  return request({
    url: '/api/deployment/simplist',
    method: 'get'
  })
}

export function getCount() {
  return request({
    url: '/api/deployment/count',
    method: 'get'
  })
}

export function getSimpListPagination(offset, limit) {
  return request({
    url: `/api/deployment/simplist?offset=${offset}&limit=${limit}`,
    method: 'get'
  })
}

export function getItem(id) {
  return request({
    url: `/api/deployment/${id}`,
    method: 'get'
  })
}

export function createItem(item) {
  return request({
    url: '/api/deployment',
    method: 'post',
    data: item
  })
}

export function editItem(item) {
  return request({
    url: `/api/deployment/${item.id}`,
    method: 'put',
    data: item
  })
}

export function deleteItem(item) {
  return request({
    url: `/api/deployment/${item.id}`,
    method: 'delete'
  })
}
