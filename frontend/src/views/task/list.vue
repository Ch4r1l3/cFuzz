<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="2">
        <router-link to="/task/create">
          <el-button type="primary">Create</el-button>
        </router-link>
      </el-col>
      <el-col :span="5">
        <el-input v-model="searchName" placeholder="name" />
      </el-col>
      <el-col :span="6">
        <el-button @click="search">Serach</el-button>
      </el-col>
    </el-row>
    <el-row>
      <el-table
        v-loading="listLoading"
        :data="items"
        element-loading-text="Loading"
        border
        fit
        highlight-current-row
      >
        <el-table-column width="48">
          <template slot-scope="scope">
            <div class="el-table__expand-icon" @click="expandHandle(scope.row)">
              <i class="el-icon el-icon-arrow-right" />
            </div>
          </template>
        </el-table-column>
        <el-table-column align="center" label="ID" width="95">
          <template slot-scope="scope">
            {{ scope.row.id }}
          </template>
        </el-table-column>
        <el-table-column label="Name">
          <template slot-scope="scope">
            {{ scope.row.name }}
          </template>
        </el-table-column>
        <el-table-column label="Status" width="110" align="center">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status | statusFilter">
              {{ scope.row.status.substring(4) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Crash Num" width="110" align="center">
          <template v-if="scope.row.status !== 'TaskCreated'" slot-scope="scope">
            {{ scope.row.crashNum }}
          </template>
        </el-table-column>
        <el-table-column label="Edit" width="95" align="center">
          <template v-if="scope.row.status === 'TaskCreated'" slot-scope="scope">
            <router-link :to="'/task/edit/'+scope.row.id">
              <el-button type="primary">
                Edit
              </el-button>
            </router-link>
          </template>
        </el-table-column>
        <el-table-column label="Delete" width="110" align="center">
          <template slot-scope="scope">
            <el-popconfirm
              confirm-button-text="OK"
              cancel-button-text="Cancel"
              icon="el-icon-info"
              icon-color="red"
              title="Delete it?"
              @onConfirm="deleteTask(scope.row)"
            >
              <el-button slot="reference" type="danger">
                Delete
              </el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
        <el-table-column label="Action" width="120" align="center">
          <template slot-scope="scope">
            <el-popconfirm
              v-if="scope.row.status==='TaskCreated'"
              confirm-button-text="OK"
              cancel-button-text="Cancel"
              icon="el-icon-info"
              icon-color="grey"
              title="Start it?"
              @onConfirm="startTask(scope.row)"
            >
              <el-button slot="reference" type="success" :loading="scope.row.loading">
                Start
              </el-button>
            </el-popconfirm>
            <el-popconfirm
              v-else-if="canStop(scope.row.status)"
              confirm-button-text="OK"
              cancel-button-text="Cancel"
              icon="el-icon-info"
              icon-color="red"
              title="Stop it?"
              @onConfirm="stopTask(scope.row)"
            >
              <el-button slot="reference" type="danger">
                Stop
              </el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
    <el-row>
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        layout="total, prev, pager, next, jumper"
        :total="count"
        @current-change="handleCurrentChange"
      />
    </el-row>
  </div>
</template>

<script>
import { getItemsCombine, deleteItem, startItem, stopItem } from '@/api/task'
import { pageSize } from '@/settings'
import { getOffset } from '@/utils'
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
      listLoading: true,
      items: [],
      count: 0,
      currentPage: 1,
      pageSize: pageSize,
      searchName: ''
    }
  },
  computed: {
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.items = []
      this.listLoading = true
      const offset = getOffset(this.currentPage, pageSize)
      getItemsCombine(offset, pageSize, this.searchName).then((data) => {
        data.data.forEach((item) => {
          item.loading = false
          this.items.push(parseServerItem(item))
        })
        this.count = data.count
        this.listLoading = false
      })
    },
    deleteTask(item) {
      deleteItem(item).then(() => {
        this.$message('delete success')
        if (this.items.length === 1 && this.currentPage > 1) {
          this.currentPage -= 1
        }
        this.fetchData()
      })
    },
    startTask(item) {
      item.loading = true
      startItem(item).then(() => {
        this.$message('start success')
        this.fetchData()
      }).finally(() => {
        item.loading = false
      })
    },
    stopTask(item) {
      stopItem(item).then(() => {
        this.$message('stop success')
        this.fetchData()
      })
    },
    handleCurrentChange(val) {
      this.currentPage = val
      this.fetchData()
    },
    getTimeString(val) {
      var t = new Date(val * 1000)
      return t.toLocaleString()
    },
    canStop(status) {
      return status === 'TaskStarted' || status === 'TaskInitializing' || status === 'TaskRunning'
    },
    expandHandle(item) {
      this.$router.push({ name: 'taskDetail', params: { id: item.id }})
    },
    search() {
      this.currentPage = 1
      this.fetchData()
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
</style>
