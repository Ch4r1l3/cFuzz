<template>
  <div class="app-container">
    <el-row>
      <router-link to="/storageItem/create">
        <el-button type="primary">Create</el-button>
      </router-link>
    </el-row>
    <el-row>
      <el-table
        v-loading="listLoading"
        :data="showItems"
        element-loading-text="Loading"
        border
        fit
        highlight-current-row
      >
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
        <el-table-column label="Path" width="100" align="center">
          <template slot-scope="scope">
            <el-tooltip v-if="scope.row.existsInImage" class="item" effect="dark" :content="scope.row.path" placement="top">
              <i class="el-icon-more-outline" />
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="Type" width="85" align="center">
          <template slot-scope="scope">
            <el-tag :type="scope.row.type | typeFilter">
              {{ scope.row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Exist in Image" width="115" align="center">
          <template slot-scope="scope">
            <i v-if="scope.row.existsInImage" class="el-icon-check" />
            <i v-else class="el-icon-close" />
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
              @onConfirm="deleteStorageItem(scope.row)"
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
        :total="items.length"
        @current-change="handleCurrentChange"
      />
    </el-row>
  </div>
</template>

<script>
import { getItems, deleteItem } from '@/api/storageItem'
import { getOffset } from '@/utils'
import { pageSize } from '@/settings'

export default {
  filters: {
    typeFilter(type) {
      const typeMap = {
        fuzzer: 'success',
        target: 'info',
        corpus: 'danger'
      }
      return typeMap[type]
    }
  },
  data() {
    return {
      listLoading: true,
      items: [],
      currentPage: 1,
      pageSize: pageSize
    }
  },
  computed: {
    showItems: function() {
      const offset = getOffset(this.currentPage, pageSize)
      return this.items.slice(offset, offset + pageSize)
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      getItems().then((data) => {
        this.items = data
        this.listLoading = false
      })
    },
    deleteStorageItem(item) {
      deleteItem(item).then(() => {
        this.$message('delete success')
        this.fetchData()
      }).catch((error) => {
        this.$message.error(error)
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
