<template>
  <div class="app-container">
    <el-row>
      <el-page-header @back="routerBack" />
      <el-form ref="form" :model="task" label-width="120px" style="margin-top: 30px">
        <el-form-item label="ID:">
          {{ task.id }}
        </el-form-item>
        <el-form-item label="Name:">
          {{ task.name }}
        </el-form-item>
        <el-form-item label="Status">
          <el-tag :type="task.status | statusFilter">
            {{ task.status.substring(4) }}
          </el-tag>
        </el-form-item>
        <el-form-item label="Time (second):">
          {{ task.time }}
        </el-form-item>
        <el-form-item v-if="task.image !== ''" label="Image:">
          {{ task.image }}
        </el-form-item>
        <el-form-item v-else label="Deployment ID:">
          {{ task.deploymentID }}
        </el-form-item>
        <el-form-item label="FuzzCycleTime:">
          {{ task.fuzzCycleTime }}
        </el-form-item>
        <el-form-item label="Fuzzer ID:">
          {{ task.fuzzerID }}
        </el-form-item>
        <el-form-item label="Corpus ID:">
          {{ task.corpusID }}
        </el-form-item>
        <el-form-item label="Target ID:">
          {{ task.targetID }}
        </el-form-item>
        <el-form-item v-if="task.errorMsg !== ''" label="Error Message:">
          {{ task.errorMsg }}
        </el-form-item>
        <el-form-item v-if="task.startedAt !== 0" label="Started At:">
          {{ getTimeString(task.startedAt) }}
        </el-form-item>
        <el-form-item v-if="task.arguments.length!==0" label="Arguments:">
          <el-table :data="task.arguments">
            <el-table-column label="Key" align="center">
              <template v-slot:default="env">
                {{ env.row.key }}
              </template>
            </el-table-column>
            <el-table-column label="Value" align="center">
              <template v-slot:default="env">
                {{ env.row.value }}
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
        <el-form-item v-if="task.environments.length!==0" label="Environments:">
          <el-table :data="task.environments">
            <el-table-column label="Key" align="center">
              <template v-slot:default="env">
                {{ env.row.key }}
              </template>
            </el-table-column>
            <el-table-column label="Value" align="center">
              <template v-slot:default="env">
                {{ env.row.value }}
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
    </el-row>
  </div>
</template>

<script>
import { getItem } from '@/api/task'
import { parseServerItem } from '@/utils/task'
export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        TaskCreated: 'info',
        TaskStarted: '',
        TaskInitializing: '',
        TaskStopped: 'danger',
        TaskError: 'warning',
        TaskRunning: 'success'
      }
      return statusMap[status]
    }
  },
  data() {
    return {
      task: {
        name: '',
        time: 3600,
        status: '',
        deploymentID: 0,
        image: '',
        fuzzCycleTime: 60,
        environments: [],
        arguments: [],
        fuzzerID: 0,
        corpusID: 0,
        targetID: 0
      },
      crashes: [],
      stats: [],
      loading: false
    }
  },
  created() {
    const id = this.$route.params && this.$route.params.id
    this.get(id)
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
    getTimeString(val) {
      var t = new Date(val * 1000)
      return t.toLocaleString()
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
</style>

