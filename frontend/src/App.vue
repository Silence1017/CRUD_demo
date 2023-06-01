<template>
  <div class="table-box">
    <!-- title -->
    <div class="title">
      <h2>最简单的 CRUD Demo</h2>
    </div>
    <!-- query -->
    <div class="query-box">
      <el-input class="query-input" v-model="queryInput" placeholder="请输入姓名搜索" @change="handleQueryName"/>
      <div class="btn-list">
        <el-button type="primary" @click="handleAdd">增加</el-button>
        <el-button type="danger" @click="handleDelList" v-if="multipleSelection.length > 0">删除多选</el-button>
      </div>
    </div>
    <!-- table -->
    <el-table border ref="multipleTableRef" :data="tableData" style="width: 100%"
      @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="姓名" width="120" />
      <el-table-column prop="phone" label="电话" width="120" />
      <el-table-column prop="email" label="邮箱" width="120" />
      <el-table-column prop="state" label="状态" width="120" />
      <el-table-column prop="address" label="地址" width="300" />

      <el-table-column fixed="right" label="操作" width="120">
        <template #default="scope">
          <el-button link type="primary" size="small" @click="handleRowDel(scope.row)" style="color: #F56C6C;">删除</el-button>
          <el-button link type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
        </template>
      </el-table-column>

    </el-table>

    <el-pagination
      background 
      layout="prev, pager, next" 
      :total="total" 
      style="display: flex;justify-content: center;margin-top: 10px;"
      v-model:current-page="curPage"
      @current-change="handleChangePage"
    />

    <!-- dialog -->
    <el-dialog v-model="dialogFormVisible" :title="dialogType === 'add' ? '新增':'编辑'">
      <el-form :model="tableForm">
        <el-form-item label="姓名" :label-width="80">
          <el-input v-model="tableForm.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="邮箱" :label-width="80">
          <el-input v-model="tableForm.email" autocomplete="off" />
        </el-form-item>
        <el-form-item label="电话" :label-width="80">
          <el-input v-model="tableForm.phone" autocomplete="off" />
        </el-form-item>
        <el-form-item label="状态" :label-width="80">
          <el-input v-model="tableForm.state" autocomplete="off" />
        </el-form-item>
        <el-form-item label="地址" :label-width="80">
          <el-input v-model="tableForm.address" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogConfirm">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>

import { ref } from 'vue';
import request from "./utils/request.js"

/* 数据 */
let queryInput = ref("")
let tableData = ref([])
// let tableDataCopy = Object.assign(tableData)
// 问题：使用Object.assign函数会影响结果
let tableDataCopy = ref([{
  id: "1",
  name: 'Tom1',
  phone: "13800138000",
  email: '123@qq.com',
  state: 'California',
  address: 'No. 189, Grove St, Los Angeles',
},
{
  id: "2",
  name: 'Tom2',
  phone: "13800138000",
  email: '123@qq.com',
  state: 'California',
  address: 'No. 189, Grove St, Los Angeles',
},
{
  id: "3",
  name: 'Tom3',
  phone: "13800138000",
  email: '123@qq.com',
  state: 'California',
  address: 'No. 189, Grove St, Los Angeles',
},
{
  id: "4",
  name: 'Tom4',
  phone: "13800138000",
  email: '123@qq.com',
  state: 'California',
  address: 'No. 189, Grove St, Los Angeles',
},])
let multipleSelection = ref([])
let dialogFormVisible = ref(false)
let tableForm = ref({
  name: '张三',
  email: "123@qq.com",
  phone: "13800138000",
  state: "在职",
  address: "广东省"
})
let dialogType = ref('add')

let total = ref(10)
let curPage = ref(1)

/* 方法 */

const getTableData = async (cur = 1) => {
  // 第一种请求方式
  let res = await request.get('/list', {
    pageSize: 10,
    pageNum: cur
  })
  console.log(res)
  // 第二种请求方式
  // let res = await request.get(`/list/?pageSize=10&pageNum=${cur}`)

  tableData.value = res.list
  total.value = res.total
  curPage.value = res.pageNum
}
getTableData(1)

// 请求分页
const handleChangePage = (val) => {
  getTableData(curPage.value)
}

// 搜索
const handleQueryName = async (val) => {
  // if(val.length > 0) {
  //   tableData.value = tableData.value.filter(item => (item.name).toLowerCase().match(val.toLowerCase()))
  // } else {
  //   tableData.value = tableDataCopy.value
  // }

  if(val.length > 0) {
    tableData.value = await request.get(`/list/${val}`)
  } else {
    await getTableData(curPage.value)
  }
}

// 编辑
const handleEdit = (row) => {
  dialogFormVisible.value = true
  dialogType.value = 'edit'
  tableForm.value = {...row}
}

// 删除一条
const handleRowDel = async ({ID}) => {
  // console.log(id)
  // // 1. 通过id获取到条目对应的索引值
  // let index = tableData.value.findIndex(item=>item.id === id)
  // // 2. 通过索引值进行删除
  // tableData.value.splice(index, 1)

  await request.delete(`/delete/${ID}`)
  await getTableData(curPage.value)

}

const handleDelList = () => {
  multipleSelection.value.forEach(ID => {
    handleRowDel({ID})
  })
  multipleSelection.value = []
}

// 选中
const handleSelectionChange = (val) => {
  // multipleSelection.value = val
  // console.log(val)
  multipleSelection.value = []

  val.forEach(item => {
    multipleSelection.value.push(item.ID)
  })
  console.log(multipleSelection.value)
}

// 新增
const handleAdd = () => {
  dialogFormVisible.value = true
  tableForm.value = {}
  dialogType.value = 'add'
}

// 确认
const dialogConfirm = async () => {
  dialogFormVisible.value = false

  // 1. 判断是新增还是编辑
  if(dialogType.value=== 'add') {
    // 1. 拿到数据
    // 2. 添加到table
    // tableData.value.push({
    //   id: (tableData.value.length + 1).toString(),
    //   ...tableForm.value
    // })
    // console.log(tableData)

    // 添加数据
    await request.post("/add", {
      ...tableForm.value
    })
    // 刷新数据
    await getTableData(curPage.value)
    
  } else {
    // // 1. 获取当前这条的索引
    // let index = tableData.value.findIndex(item => item.id === tableForm.value.id)
    // console.log(index)
    // // 2. 替换当前索引值对应的数据
    // tableData.value[index] = tableForm.value

    await request.put(`/update/${tableForm.value.ID}`, {
      ...tableForm.value
    })
    await getTableData(curPage.value)
  }
}
</script>

<style scoped>
.table-box {
  width: 800px;
  margin: 200px auto;
}

.title {
  text-align: center;
}

.query-box {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.query-input {
  width: 200px;
}
</style>
