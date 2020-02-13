<template>
  <el-form :id="'image'+imageId" label-position="left">
    <el-form-item label="Name">
      {{ image.name }}
    </el-form-item>
    <el-form-item v-if="image.isDeployment" label="Content">
      <codemirror :value="image.content" :options="{ readOnly: true, theme: 'monokai', mode: 'text/x-yaml' }" />
    </el-form-item>
    <el-form-item v-else label="Image">
      {{ image.content }}
    </el-form-item>
  </el-form>
</template>

<script>
import { getItem } from '@/api/image'

export default {
  name: 'ImageExpand',
  props: {
    imageId: {
      type: Number,
      default: 1
    }
  },
  data() {
    return {
      image: {
        name: '',
        content: '',
        isDeploymente: false
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.$nextTick(() => {
        const loading = this.$loading({
          lock: true,
          fullscreen: false,
          target: `#image${this.imageId}`
        })
        getItem(this.imageId).then((res) => {
          this.image = res
          loading.close()
        })
      })
    },
    handleCurrentChange(val) {
      this.currentPage = val
    }
  }
}
</script>

<style>
.el-row {
  margin-bottom: 20px;
  &:last-child {
    margin-bottom: 0;
  }
}
.el-form-item__content .CodeMirror {
    line-height: 25px;
}

.CodeMirror-code {
    padding: 8px;
}
</style>
