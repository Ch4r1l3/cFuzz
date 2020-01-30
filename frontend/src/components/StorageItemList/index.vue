<template>
  <div class="app-container">
    <el-row>
      <el-table
        v-loading="listLoading"
        :data="showItems"
        element-loading-text="Loading"
        border
        fit
        highlight-current-row
        max-height="300"
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
        :total="items.length"
        @current-change="handleCurrentChange"
      />
    </el-row>
  </div>
</template>

<script>
import { getItemsByType } from '@/api/storageItem'
import { getOffset } from '@/utils'
import { pageSize } from '@/settings'

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
      getItemsByType(this.storageItemType).then((data) => {
        this.items = data
        this.listLoading = false
      })
    },
    chooseStorageItem(item) {
      this.$emit('choose', item)
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
