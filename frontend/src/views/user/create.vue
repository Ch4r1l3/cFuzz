<template>
  <div class="app-container">
    <el-row>
      <el-col :span="10">
        <el-form ref="form" :rules="rules" :model="user" label-width="120px">
          <el-form-item label="Name" prop="username">
            <el-input v-model="user.username" />
          </el-form-item>
          <el-form-item label="Password" prop="password">
            <el-input v-model="user.password" show-password />
          </el-form-item>
          <el-form-item label="Confirm Password" prop="confirmPassword">
            <el-input v-model="user.confirmPassword" show-password />
          </el-form-item>
          <el-form-item>
            <el-button v-loading="loading" type="primary" @click="create">Create</el-button>
            <el-button @click="routerBack">Cancel</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { createItem } from '@/api/user'
export default {
  name: 'CreateUser',
  data() {
    var checkConfirm = (rule, value, callback) => {
      if (this.user.confirmPassword === '') {
        callback(new Error('please repeat the password'))
      } else if (this.user.confirmPassword !== this.user.password) {
        callback(new Error('two password are not equal'))
      } else {
        callback()
      }
    }
    return {
      user: {
        name: '',
        password: '',
        confirmPassword: ''
      },
      loading: false,
      rules: {
        username: [
          { required: true, trigger: 'blur' },
          { max: 20, trigger: 'blur' }
        ],
        password: [
          { required: true, trigger: 'blur' },
          { min: 6, max: 18, trigger: 'blur' }
        ],
        confirmPassword: [
          { validator: checkConfirm, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    routerBack() {
      this.$router.back(-1)
    },
    create() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.loading = true
          createItem(this.user).then(() => {
            this.$message('create success')
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

