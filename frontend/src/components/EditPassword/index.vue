<template>
  <div class="user-container">
    <el-form ref="form" :rules="rules" :model="user" label-position="top">
      <el-form-item v-if="isNormal" label="Old Password" prop="oldpassword">
        <el-input v-model="user.oldpassword" show-password />
      </el-form-item>
      <el-form-item label="Password" prop="password">
        <el-input v-model="user.password" show-password />
      </el-form-item>
      <el-form-item label="Confirm Password" prop="confirmPassword">
        <el-input v-model="user.confirmPassword" show-password />
      </el-form-item>
      <el-form-item>
        <el-button v-loading="loading" type="primary" @click="edit">Edit</el-button>
        <el-button @click="onClose">Cancel</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { adminEditItem, editItem } from '@/api/user'

export default {
  name: 'EditPassword',
  props: {
    userId: {
      type: Number,
      default: 0
    },
    isNormal: {
      type: Boolean,
      default: true
    }
  },
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
    var checkOldPassword = (rule, value, callback) => {
      if (!this.isNormal) {
        callback()
        return
      }
      if (this.user.oldpassword === '') {
        callback(new Error('old password is required'))
      } else if (this.user.oldpassword.length < 6 || this.user.oldpassword.length > 18) {
        callback(new Error('old password must be between 6 and 18 characters'))
      } else {
        callback()
      }
    }
    return {
      loading: false,
      user: {
        oldpassword: '',
        password: '',
        confirmPassword: ''
      },
      rules: {
        oldpassword: [
          { validator: checkOldPassword, trigger: 'blur' }
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
  mounted: function() {
    window.addEventListener('pageshow', this.reset)
  },
  methods: {
    onClose() {
      this.$emit('onClose', true)
    },
    edit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          if (this.isNormal) {
            editItem(this.userId, this.user.oldpassword, this.user.password).then(() => {
              this.$message('edit success')
              this.onClose()
            })
          } else {
            adminEditItem(this.userId, this.user.password).then(() => {
              this.$message('edit success')
              this.onClose()
            })
          }
        }
      })
    },
    reset() {
      this.$refs.form.resetFields()
    }
  }
}
</script>

<style>
</style>
