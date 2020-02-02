export function getServerItem(item) {
  const temp = Object.assign({}, item)
  var tempArguement = {}
  temp.arguments.forEach((arguement) => {
    tempArguement[arguement.key] = arguement.value
  })
  var tempEnvironment = []
  temp.environments.forEach((environment) => {
    tempEnvironment.push(environment.key + '=' + environment.value)
  })
  temp.arguments = tempArguement
  temp.environments = tempEnvironment
  return temp
}

export function parseServerItem(item) {
  var tempArguement = []
  Object.keys(item.arguments).forEach((key) => {
    tempArguement.push({ 'key': key, 'value': item.arguments[key] })
  })
  var tempEnvironment = []
  item.environments.forEach((val) => {
    const index = val.indexOf('=')
    if (index !== -1) {
      tempEnvironment.push({ 'key': val.substring(0, index), 'value': val.substring(index + 1) })
    }
  })
  item.arguments = tempArguement
  item.environments = tempEnvironment
  return item
}

export function parseServerStats(item) {
  var stats = []
  Object.keys(item).forEach((key) => {
    if (key !== 'stats' && key !== 'id' && key !== 'taskid') {
      stats.push({ 'key': key, 'value': item[key] })
    }
  })
  Object.keys(item.stats).forEach((key) => {
    stats.push({ 'key': key, 'value': item.stats[key] })
  })
  return stats
}
