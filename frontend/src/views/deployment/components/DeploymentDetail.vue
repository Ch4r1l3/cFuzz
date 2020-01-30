<template>
  <div class="deployment-container">
    <el-row>
      <el-form ref="form" :model="deployment" label-width="120px">
        <el-form-item label="Name">
          <el-input v-model="deployment.name" />
        </el-form-item>
        <el-form-item label="Content">
          <codemirror v-model="deployment.content" />
        </el-form-item>
        <el-form-item>
          <el-button v-if="isEdit" v-loading="loading" type="primary" @click="edit">Edit</el-button>
          <el-button v-else v-loading="loading" type="primary" @click="create">Create</el-button>
          <el-button @click="routerBack">Cancel</el-button>
        </el-form-item>
      </el-form>
    </el-row>
  </div>
</template>

<script>
import { getItem, createItem, editItem } from '@/api/deployment'
export default {
  name: 'DeploymentDetail',
  props: {
    isEdit: {
      type: Boolean,
      defualt: false
    }
  },
  data() {
    return {
      deployment: {
        name: '',
        content: ''
      },
      loading: false
    }
  },
  created() {
    if (this.isEdit) {
      const id = this.$route.params && this.$route.params.id
      this.get(id)
    }
  },
  methods: {
    routerBack() {
      this.$router.back(-1)
    },
    get(id) {
      this.loading = true
      getItem(id).then((data) => {
        this.deployment = data
        this.loading = false
      }).catch(error => {
        this.$message({
          message: error,
          type: 'warning'
        })
      })
    },
    create() {
      if (this.deployment.name.length === 0) {
        this.$message({
          message: 'name cannot be empty',
          type: 'warning'
        })
        return
      }
      this.loading = true
      createItem(this.deployment).then(() => {
        this.loading = false
        this.$message('create success')
        this.routerBack()
      }).catch((error) => {
        this.loading = false
        this.$message({
          message: error,
          type: 'warning'
        })
      })
    },
    edit() {
      if (this.deployment.name.length === 0) {
        this.$message({
          message: 'name cannot be empty',
          type: 'warning'
        })
        return
      }
      this.loading = true
      editItem(this.deployment).then(() => {
        this.loading = false
        this.$message('edit success')
        this.routerBack()
      }).catch((error) => {
        this.loading = false
        this.$message({
          message: error,
          type: 'warning'
        })
      })
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
</style>

