<template>
  <div class="app-container">
    <el-row>
      <el-page-header @back="routerBack" />
      <el-form id="detail-form" ref="form" :model="task" label-width="120px" style="margin-top: 30px">
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
        <el-form-item label="Image ID:">
          {{ task.imageID }}
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
      </el-form>
      <el-form id="detail-tables" :model="task" label-position="top" label-width="120px">
        <el-form-item v-if="task.arguments.length!==0" label="Arguments:">
          <el-table :data="task.arguments" :max-height="300">
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
          <el-table :data="task.environments" :max-height="300">
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
        <el-form-item v-if="crashes.length !== 0" id="detail-crash" label="Crashes:">
          <el-table :data="crashes" :max-height="300">
            <el-table-column label="ID" align="center">
              <template v-slot:default="crash">
                {{ crash.row.id }}
              </template>
            </el-table-column>
            <el-table-column label="fileName" align="center">
              <template v-slot:default="crash">
                {{ crash.row.fileName }}
              </template>
            </el-table-column>
            <el-table-column label="ReproduceAble" align="center">
              <template v-slot:default="crash">
                {{ crash.row.reproduceAble }}
              </template>
            </el-table-column>
            <el-table-column label="Action" align="center">
              <template v-slot:default="crash">
                <el-link :href="`api/crash/${crash.row.id}`" target="_blank" type="primary">Download</el-link>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
        <el-form-item v-if="stats.length !== 0" id="detail-stats" label="Stats:">
          <el-table :data="stats" :max-height="300">
            <el-table-column label="Key" align="center">
              <template v-slot:default="stat">
                {{ stat.row.key }}
              </template>
            </el-table-column>
            <el-table-column label="Value" align="center">
              <template v-slot:default="stat">
                {{ stat.row.value }}
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>

      </el-form>
    </el-row>
  </div>
</template>

<script>
import { getItem, getCrashes, getResult, downloadCrash } from '@/api/task'
import { parseServerItem, parseServerStats } from '@/utils/task'
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
        imageID: 0,
        fuzzCycleTime: 60,
        environments: [],
        arguments: [],
        fuzzerID: 0,
        corpusID: 0,
        targetID: 0
      },
      crashes: [],
      stats: []
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
      this.$nextTick(() => {
        const loading1 = this.$loading({
          lock: true,
          fullscreen: false,
          target: '#detail-form'
        })
        getItem(id).then((data) => {
          this.task = parseServerItem(data)
        }).finally(() => {
          loading1.close()
        })
        const loading2 = this.$loading({
          lock: true,
          fullscreen: false,
          target: '#detail-crash'
        })
        getCrashes(id).then((res) => {
          this.crashes = res.data
        }).finally(() => {
          loading2.close()
        })
        const loading3 = this.$loading({
          lock: true,
          fullscreen: false,
          target: '#detail-stats'
        })
        getResult(id).then((res) => {
          this.stats = parseServerStats(res)
        }).finally(() => {
          loading3.close()
        })
      })
    },
    getTimeString(val) {
      var t = new Date(val * 1000)
      return t.toLocaleString()
    },
    downloadCrash(item) {
      downloadCrash(item.id)
    }
  }
}
</script>

<style>
#detail-tables label{
    width: 120px;
    padding: 0 12px 0 0;
    text-align: right;
}
</style>

