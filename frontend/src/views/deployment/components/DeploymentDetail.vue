<template>
  <div class="deployment-container">
    <el-row>
      <el-form ref="form" :rules="rules" :model="deployment" label-width="120px">
        <el-form-item label="Name" prop="name">
          <el-input v-model="deployment.name" />
        </el-form-item>
        <el-form-item label="Content" prop="content">
          <codemirror v-model="deployment.content" :options="{ theme: 'monokai', mode: 'text/x-yaml', lineNumbers: true }" />
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
      loading: false,
      rules: {
        name: [
          { required: true, message: 'please input the name', trigger: 'blur' }
        ],
        content: [
          { required: true, message: 'please input the content', trigger: 'change' }
        ]
      }
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
      })
    },
    create() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.loading = true
          createItem(this.deployment).then(() => {
            this.$message('create success')
            this.routerBack()
          }).finally(() => {
            this.loading = false
          })
        }
      })
    },
    edit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.loading = true
          editItem(this.deployment).then(() => {
            this.$message('edit success')
            this.routerBack()
          }).finally(() => {
            this.loading = false
          })
        }
      })
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
.el-form-item__content .CodeMirror {
    line-height: 25px;
}
</style>

