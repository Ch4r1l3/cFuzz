<template>
  <div>
    <el-row>
      <el-table
        v-loading="listLoading"
        :data="items"
        element-loading-text="Loading"
        border
        fit
        highlight-current-row
        max-height="300"
        style="width: 100%"
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
        <el-table-column label="Action" width="150" align="center">
          <template slot-scope="scope">
            <el-button @click="emitChoose(scope.row)">
              Choose
            </el-button>
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
import { getCount, getSimpListPagination } from '@/api/deployment'
import { pageSize } from '@/settings'
import { getOffset } from '@/utils'
import DeploymentExpand from '@/components/DeploymentExpand'

export default {
  name: 'DeploymentList',
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
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      const offset = getOffset(this.currentPage, pageSize)
      getSimpListPagination(offset, pageSize).then((data) => {
        this.items = data
        getCount().then((res) => {
          this.count = res.count
          this.listLoading = false
        })
      })
    },
    handleCurrentChange(val) {
      this.currentPage = val
    },
    emitChoose(val) {
      this.$emit('choose', val)
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
