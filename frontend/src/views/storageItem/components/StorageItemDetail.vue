<template>
  <div class="storageItem-container">
    <el-row>
      <el-form ref="form" :model="storageItem" label-width="120px">
        <el-form-item label="Name">
          <el-input v-model="storageItem.name" />
        </el-form-item>
        <el-form-item label="Type">
          <el-select v-model="storageItem.type" placeholder="please select the type">
            <el-option label="Fuzzer" value="fuzzer" />
            <el-option label="Corpus" value="corpus" />
            <el-option label="Target" value="target" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!isEdit" label="Exist In Image">
          <el-switch v-model="storageItem.existsInImage" />
        </el-form-item>
        <el-form-item v-if="storageItem.existsInImage" label="Path">
          <el-input v-model="storageItem.path" />
        </el-form-item>
        <el-form-item v-else>
          <el-upload
            ref="upload"
            class="upload-demo"
            drag
            action="/api/storage_item"
            :auto-upload="false"
            :file-list="fileList"
            :data="storageItem"
            :limit="1"
            :on-success="uploadSuccess"
            :on-error="uploadError"
          >
            <i class="el-icon-upload" />
            <div class="el-upload__text">drag file here, or <em>choose file</em></div>
          </el-upload>
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
import { getItem, createItem, editItem } from '@/api/storageItem'
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
      storageItem: {
        name: '',
        type: '',
        existsInImage: true,
        path: ''
      },
      loading: false,
      fileList: []
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
        this.storageItem = data
        this.loading = false
      }).catch(error => {
        this.$message({
          message: error,
          type: 'warning'
        })
      })
    },
    create() {
      if (this.storageItem.existsInImage) {
        createItem(this.storageItem).then(() => {
          this.$message('create susscess')
          this.routerBack()
        }).catch((error) => {
          this.$message.error(error)
        })
      } else {
        if (this.storageItem.name.length === 0) {
          this.$message({
            message: 'name cannot be empty',
            type: 'warning'
          })
          return
        }
        console.log(this.$refs.upload.submit())
      }
    },
    edit() {
      if (this.storageItem.name.length === 0) {
        this.$message({
          message: 'name cannot be empty',
          type: 'warning'
        })
        return
      }
      this.loading = true
      editItem(this.storageItem).then(() => {
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
    },
    uploadSuccess() {
      this.$message('create susscess')
      this.routerBack()
    },
    uploadError() {
      this.$message.error('create failed')
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
</style>

