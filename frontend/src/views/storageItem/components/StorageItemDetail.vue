<template>
  <div class="storageItem-container">
    <el-row>
      <el-form ref="form" :rules="rules" :model="storageItem" label-width="120px">
        <el-form-item label="Name" prop="name">
          <el-input v-model="storageItem.name" />
        </el-form-item>
        <el-form-item label="Type" prop="type">
          <el-select v-model="storageItem.type" placeholder="please select the type">
            <el-option label="Fuzzer" value="fuzzer" />
            <el-option label="Corpus" value="corpus" />
            <el-option label="Target" value="target" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!isEdit" label="Exist In Image">
          <el-switch v-model="storageItem.existsInImage" />
        </el-form-item>
        <el-form-item v-if="storageItem.existsInImage" label="Path" prop="path">
          <el-input v-model="storageItem.path" />
        </el-form-item>
        <el-form-item v-else prop="file">
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
import { getItem, createItem } from '@/api/storageItem'
export default {
  name: 'DeploymentDetail',
  props: {
    isEdit: {
      type: Boolean,
      defualt: false
    }
  },
  data() {
    var checkPath = (rule, value, callback) => {
      if (this.storageItem.existsInImage && value === '') {
        callback(new Error('please input the path'))
      } else {
        callback()
      }
    }
    var checkFile = (rule, value, callback) => {
      if (!this.storageItem.existsInImage && this.$refs.upload.uploadFiles.length === 0) {
        callback(new Error('please upload the file'))
      } else {
        callback()
      }
    }
    return {
      storageItem: {
        name: '',
        type: '',
        existsInImage: true,
        path: ''
      },
      loading: false,
      fileList: [],
      rules: {
        name: [
          { required: true, message: 'please input the name', trigger: 'blur' }
        ],
        type: [
          { required: true, message: 'please select the type', trigger: 'change' }
        ],
        path: [
          { validator: checkPath, trigger: 'blur' }
        ],
        file: [
          { validator: checkFile, trigger: 'change' }
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
        this.storageItem = data
        this.loading = false
      })
    },
    create() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          if (this.storageItem.existsInImage) {
            createItem(this.storageItem).then(() => {
              this.$message('create susscess')
              this.routerBack()
            })
          } else {
            this.$refs.upload.submit()
          }
        }
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

