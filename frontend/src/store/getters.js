const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  name: state => state.user.name,
  id: state => state.user.id,
  isAdmin: state => state.user.isAdmin,
  permission_routes: state => state.permission.routes
}
export default getters
