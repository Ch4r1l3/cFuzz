<template>
  <div class="task-container">
    <el-row>
      <el-form ref="form" :rules="rules" :model="task" label-width="120px">
        <el-form-item label="Name" prop="name">
          <el-input v-model="task.name" />
        </el-form-item>
        <el-form-item label="Time (second)" prop="time">
          <el-input-number v-model.number="task.time" :min="1" type="number" />
        </el-form-item>
        <el-form-item label="UseDeployment">
          <el-switch v-model="useDeployment" />
        </el-form-item>
        <el-form-item v-if="!useDeployment" label="Image" prop="image">
          <el-input v-model="task.image" />
        </el-form-item>
        <el-form-item v-else label="Deployment ID" prop="deploymentID">
          <el-input v-model.number="task.deploymentID" placeholder="input Deployment ID" class="id-input">
            <el-button slot="suffix" style="margin-right: 5px" type="text" @click="deployDialogVisible=true"> Choose </el-button>
          </el-input>
        </el-form-item>
        <el-form-item label="FuzzCycleTime (second)" prop="fuzzCycleTime">
          <el-input-number v-model.number="task.fuzzCycleTime" :min="1" />
        </el-form-item>
        <el-form-item label="Fuzzer ID" prop="fuzzerID">
          <el-input v-model.number="task.fuzzerID" placeholder="input Fuzzer ID" class="id-input">
            <el-button slot="suffix" style="margin-right: 5px" type="text" @click="storageDialogVisible.fuzzer = true"> Choose </el-button>
          </el-input>
        </el-form-item>
        <el-form-item label="Corpus ID" prop="corpusID">
          <el-input v-model.number="task.corpusID" placeholder="input Corpus ID" class="id-input">
            <el-button slot="suffix" style="margin-right: 5px" type="text" @click="storageDialogVisible.corpus = true"> Choose </el-button>
          </el-input>
        </el-form-item>
        <el-form-item label="Target ID" prop="targetID">
          <el-input v-model.number="task.targetID" placeholder="input Target ID" class="id-input">
            <el-button slot="suffix" style="margin-right: 5px" type="text" @click="storageDialogVisible.target = true"> Choose </el-button>
          </el-input>
        </el-form-item>
        <el-form-item label="Arguments">
          <div v-for="(argument,index) in task.arguments" :key="'argument'+index" class="input-kv">
            <el-input v-model="argument.key" placeholder="Key" />
            <i>-</i>
            <el-input v-model="argument.value" placeholder="Value" />
            <el-button type="text" icon="el-icon-close" style="color: #F56C6C" @click="deleteArgument(index)" />
          </div>
          <el-button type="text" @click="addArgument">Add</el-button>
        </el-form-item>
        <el-form-item label="Environments">
          <div v-for="(environment,index) in task.environments" :key="'environment'+index" class="input-kv">
            <el-input v-model="environment.key" placeholder="Key" />
            <i>-</i>
            <el-input v-model="environment.value" placeholder="Value" />
            <el-button type="text" icon="el-icon-close" style="color: #F56C6C" @click="deleteEnvironment(index)" />
          </div>
          <el-button type="text" @click="addEnvironment">Add</el-button>
        </el-form-item>

        <el-form-item>
          <el-button v-if="isEdit" v-loading="loading" type="primary" @click="edit">Edit</el-button>
          <el-button v-else v-loading="loading" type="primary" @click="create">Create</el-button>
          <el-button @click="routerBack">Cancel</el-button>
        </el-form-item>
      </el-form>
    </el-row>
    <el-dialog
      title="Choose Deployment"
      :visible="deployDialogVisible"
      :before-close="() => {deployDialogVisible = false}"
    >
      <deployment-list @choose="chooseDeploy" />
    </el-dialog>
    <el-dialog
      title="Choose Fuzzer"
      :visible="storageDialogVisible.fuzzer"
      :before-close="() => {storageDialogVisible.fuzzer = false}"
    >
      <storage-item-list storage-item-type="fuzzer" @choose="(item) => chooseStorageItem('fuzzer', item)" />
    </el-dialog>
    <el-dialog
      title="Choose Corpus"
      :visible="storageDialogVisible.corpus"
      :before-close="() => {storageDialogVisible.corpus = false}"
    >
      <storage-item-list storage-item-type="corpus" @choose="(item) => chooseStorageItem('corpus', item)" />
    </el-dialog>
    <el-dialog
      title="Choose Target"
      :visible="storageDialogVisible.target"
      :before-close="() => {storageDialogVisible.target = false}"
    >
      <storage-item-list storage-item-type="target" @choose="(item) => chooseStorageItem('target', item)" />
    </el-dialog>
  </div>
</template>

<script>
import { getItem, createItem, editItem } from '@/api/task'
import DeploymentList from '@/components/DeploymentList'
import StorageItemList from '@/components/StorageItemList'
import { getServerItem, parseServerItem } from '@/utils/task'
export default {
  name: 'TaskCreateEdit',
  components: { DeploymentList, StorageItemList },
  props: {
    isEdit: {
      type: Boolean,
      defualt: false
    }
  },
  data() {
    var checkImage = (rule, value, callback) => {
      if (!this.useDeployment && this.task.image === '') {
        callback(new Error('image is required'))
      } else {
        callback()
      }
    }
    var checkDeployment = (rule, value, callback) => {
      if (this.useDeployment && !this.task.deploymentID) {
        callback(new Error('deploymentID is required'))
      } else if (this.useDeployment && this.task.deploymentID.constructor.name !== 'Number') {
        callback(new Error('deploymentID is not a number'))
      } else if (this.useDeployment && this.task.deploymentID <= 0) {
        callback(new Error('deployment ID should large than zero'))
      } else {
        callback()
      }
    }
    return {
      task: {
        name: '',
        time: 3600,
        deploymentID: 0,
        image: '',
        fuzzCycleTime: 60,
        environments: [],
        arguments: [],
        fuzzerID: 0,
        corpusID: 0,
        targetID: 0
      },
      useDeployment: false,
      loading: false,
      deployDialogVisible: false,
      storageDialogVisible: {
        fuzzer: false,
        corpus: false,
        target: false
      },
      rules: {
        name: [
          { required: true, trigger: 'blur' }
        ],
        time: [
          { type: 'number', required: true, trigger: 'change' }
        ],
        image: [
          { validator: checkImage, trigger: 'change' }
        ],
        deploymentID: [
          { type: 'number', validator: checkDeployment }
        ],
        fuzzCycleTime: [
          { type: 'number', required: true }
        ],
        fuzzerID: [
          { type: 'number', required: true },
          { type: 'number', min: 1 }
        ],
        targetID: [
          { type: 'number', required: true },
          { type: 'number', min: 1 }
        ],
        corpusID: [
          { type: 'number', required: true },
          { type: 'number', min: 1 }
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
        this.task = parseServerItem(data)
        this.useDeployment = this.task.deploymentID !== 0
        this.loading = false
      })
    },
    create() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.loading = true
          const temp = getServerItem(this.task)
          if (this.useDeployment) {
            temp.image = ''
          } else {
            temp.deploymentID = 0
          }
          createItem(temp).then(() => {
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
          const temp = getServerItem(this.task)
          if (this.useDeployment) {
            temp.image = ''
          } else {
            temp.deploymentID = 0
          }
          editItem(temp).then(() => {
            this.$message('edit success')
            this.routerBack()
          }).finally(() => {
            this.loading = false
          })
        }
      })
    },
    addArgument() {
      this.task.arguments.push({ key: '', value: '' })
    },
    deleteArgument(index) {
      this.task.arguments.splice(index, 1)
    },
    addEnvironment() {
      this.task.environments.push({ key: '', value: '' })
    },
    deleteEnvironment(index) {
      this.task.environments.splice(index, 1)
    },
    chooseDeploy(deploy) {
      this.deployDialogVisible = false
      this.task.deploymentID = deploy.id
    },
    chooseStorageItem(type, storageItem) {
      this.storageDialogVisible[type] = false
      this.task[type + 'ID'] = storageItem.id
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
.el-form .el-input {
  width: 350px;
}
.el-form .id-input {
  width: 200px;
}
.input-kv .el-input {
  width: 130px;
}
.input-kv {
  margin-top: 10px;
}
</style>

