import request from '@/utils/request'

export function getItemsCombine(offset, limit, name) {
  return request({
    url: `/api/task?offset=${offset}&limit=${limit}&name=${name}`,
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

export function getCrashes(id) {
  return request({
    url: `/api/task/${id}/crash`,
    method: 'get'
  })
}

export function getResult(id) {
  return request({
    url: `/api/task/${id}/result`,
    method: 'get'
  })
}
