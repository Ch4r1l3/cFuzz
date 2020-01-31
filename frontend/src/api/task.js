import request from '@/utils/request'

export function getCount() {
  return request({
    url: '/api/task/count',
    method: 'get'
  })
}

export function getItemsPagination(offset, limit) {
  return request({
    url: `/api/task?offset=${offset}&limit=${limit}`,
    method: 'get'
  })
}

export function getItem(id) {
  return request({
    url: `/api/task/${id}`,
    method: 'get'
  })
}

export function createItem(item) {
  return request({
    url: '/api/task',
    method: 'post',
    data: item
  })
}

export function editItem(item) {
  return request({
    url: `/api/task/${item.id}`,
    method: 'put',
    data: item
  })
}

export function deleteItem(item) {
  return request({
    url: `/api/task/${item.id}`,
    method: 'delete'
  })
}

export function startItem(item) {
  return request({
    url: `/api/task/${item.id}/start`,
    method: 'post'
  })
}

export function stopItem(item) {
  return request({
    url: `/api/task/${item.id}/stop`,
    method: 'post'
  })
}

export function getCrashes(item) {
  return request({
    url: `/api/task/${item.id}/crash`,
    method: 'get'
  })
}

export function getResult(item) {
  return request({
    url: `/api/task/${item.id}/result`,
    method: 'get'
  })
}
