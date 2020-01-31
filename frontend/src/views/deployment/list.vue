<template>
  <div class="app-container">
    <el-row>
      <router-link to="/deployment/create">
        <el-button type="primary">Create</el-button>
      </router-link>
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
        <el-table-column type="expand">
          <template :id="'deploy'+scope.row.id" slot-scope="scope">
            <deployment-expand :deploy-id="scope.row.id" />
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
        <el-table-column label="Edit" width="95" align="center">
          <template slot-scope="scope">
            <router-link :to="'/deployment/edit/'+scope.row.id">
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
              @onConfirm="deleteDeploy(scope.row)"
            >
              <el-button slot="reference" type="danger">
                Delete
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
import { getCount, getSimpListPagination, deleteItem } from '@/api/deployment'
import { pageSize } from '@/settings'
import { getOffset } from '@/utils'
import DeploymentExpand from '@/components/DeploymentExpand'

export default {
  components: { DeploymentExpand },
  data() {
    return {
      listLoading: true,
      items: [],
      count: 0,
      currentPage: 1,
      pageSize: pageSize
    }
  },
  computed: {
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      const offset = getOffset(this.currentPage, pageSize)
      getSimpListPagination(offset, pageSize).then((data) => {
        data.forEach((item, index) => {
          data[index].content = ''
        })
        this.items = data
        getCount().then((res) => {
          this.count = res.count
          this.listLoading = false
        })
      })
    },
    deleteDeploy(item) {
      deleteItem(item).then(() => {
        this.$message('delete success')
        this.fetchData()
      })
    },
    handleCurrentChange(val) {
      this.currentPage = val
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
