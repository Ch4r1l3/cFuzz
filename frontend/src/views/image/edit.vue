<template>
  <div class="app-container">
    <image-detail :is-edit="true" />
  </div>
</template>

<script>
import ImageDetail from './components/ImageDetail'
export default {
  name: 'EditImage',
  components: { ImageDetail }
}
</script>

