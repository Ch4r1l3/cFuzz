<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="2">
        <router-link to="/image/create">
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
        <el-table-column v-if="isAdmin" align="center" label="User ID" width="95">
          <template slot-scope="scope">
            {{ scope.row.userID }}
          </template>
        </el-table-column>
        <el-table-column label="Name">
          <template slot-scope="scope">
            {{ scope.row.name }}
          </template>
        </el-table-column>
        <el-table-column label="Edit" width="95" align="center">
          <template slot-scope="scope">
            <router-link :to="'/image/edit/'+scope.row.id">
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
              @onConfirm="deleteImage(scope.row)"
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
import { mapGetters } from 'vuex'
import { getItemsCombine, deleteItem } from '@/api/image'
import { pageSize } from '@/settings'
import { getOffset } from '@/utils'
import ImageExpand from '@/components/ImageExpand'

export default {
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
  computed: {
    ...mapGetters([
      'isAdmin'
    ])
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
    deleteImage(item) {
      deleteItem(item).then(() => {
        this.$message('delete success')
        if (this.items.length === 1 && this.currentPage > 1) {
          this.currentPage -= 1
        }
        this.fetchData()
      })
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
