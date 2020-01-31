<template>
  <el-form :id="'deploy'+deployId" label-position="left">
    <el-form-item label="Name">
      {{ deployment.name }}
    </el-form-item>
    <el-form-item label="Content">
      <codemirror :value="deployment.content" :options="{ readOnly: true, theme: 'monokai', mode: 'text/x-yaml' }" />
    </el-form-item>
  </el-form>
</template>

<script>
import { getItem } from '@/api/deployment'

export default {
  name: 'DeploymentExpand',
  props: {
    deployId: {
      type: Number,
      default: 1
    }
  },
  data() {
    return {
      deployment: {
        name: '',
        content: ''
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
          target: `#deploy${this.deployId}`
        })
        getItem(this.deployId).then((res) => {
          this.deployment = res
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
    padding: 8px;
}
</style>
