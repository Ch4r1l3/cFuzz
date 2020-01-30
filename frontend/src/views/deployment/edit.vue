<template>
  <div class="app-container">
    <deployment-detail :is-edit="true" />
  </div>
</template>

<script>
import DeploymentDetail from './components/DeploymentDetail'
export default {
  name: 'EditDeployment',
  components: { DeploymentDetail }
}
</script>

