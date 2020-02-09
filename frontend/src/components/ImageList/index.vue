<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="5">
        <el-input v-model="searchName" placeholder="name" />
      </el-col>
      <el-col :span="6">
        <el-button @click="searchName">Serach</el-button>
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
        max-height="300"
        style="width: 100%"
      >
        <el-table-column type="expand">
          <template :id="'image'+scope.row.id" slot-scope="scope">
            <image-expand :image-id="scope.row.id" />
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
import { getItemsCombine } from '@/api/image'
import { pageSize } from '@/settings'
import { getOffset } from '@/utils'
import ImageExpand from '@/components/ImageExpand'

export default {
  name: 'ImageList',
  components: { ImageExpand },
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
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      const offset = getOffset(this.currentPage, pageSize)
      getItemsCombine(offset, pageSize, this.searchName).then((data) => {
        this.items = data.data
        this.count = data.count
        this.listLoading = false
      })
    },
    handleCurrentChange(val) {
      this.currentPage = val
    },
    emitChoose(val) {
      this.$emit('choose', val)
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
