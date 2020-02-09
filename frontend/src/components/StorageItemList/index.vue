<template>
  <div class="task-container">
    <el-row :gutter="20">
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
        max-height="320"
        style="width: 100%"
      >
        <el-table-column type="expand">
          <template slot-scope="scope">
            <el-form label-position="left">
              <el-form-item label="Exist in Image">
                <i v-if="scope.row.existsInImage" class="el-icon-check" />
                <i v-else class="el-icon-close" />
              </el-form-item>
              <el-form-item v-if="scope.row.existsInImage" label="Path">
                {{ scope.row.path }}
              </el-form-item>
            </el-form>
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
            <el-button slot="reference" @click="chooseStorageItem(scope.row)">
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
import { getItemsByTypeCombine } from '@/api/storageItem'
import { getOffset } from '@/utils'

export default {
  name: 'StorageItemList',
  props: {
    storageItemType: {
      type: String,
      default: 'fuzzer'
    }
  },
  data() {
    return {
      listLoading: true,
      items: [],
      currentPage: 1,
      pageSize: 4,
      count: 0,
      searchName: ''
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      const offset = getOffset(this.currentPage, this.pageSize)
      getItemsByTypeCombine(this.storageItemType, offset, this.pageSize, this.searchName).then((data) => {
        this.items = data.data
        this.count = data.count
        this.listLoading = false
      })
    },
    chooseStorageItem(item) {
      this.$emit('choose', item)
    },
    handleCurrentChange(val) {
      this.currentPage = val
      this.fetchData()
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
