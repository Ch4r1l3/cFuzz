<template>
  <div class="app-container">
    <storage-item-detail :is-edit="true" />
  </div>
</template>

<script>
import StorageItemDetail from './components/StorageItemDetail'
export default {
  name: 'EditStorageItem',
  components: { StorageItemDetail }
}
</script>

