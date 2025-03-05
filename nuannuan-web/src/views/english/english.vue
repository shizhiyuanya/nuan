<script setup lang="ts">

import {computed, onMounted, ref} from 'vue'
import axios from "axios";
import {ElMessage} from "element-plus";

interface PaginationParams {
  page: number;
  pageSize: number;
  english?: string; // 可选参数，用于搜索
  chinese?: string; // 可选参数，用于搜索
}

interface English {
  english: string,
  chinese: string
}

const currentPage = ref(1)
const pageSize = ref(20)
const inputEnglish = ref('')
const inputChinese = ref('')
const totalPage = ref(10)

const currentShowSize = ref(pageSize.value)

const defaultEnglish: English = { english: "abandon", chinese: "abandon" };
// 使用默认值填充数组
const showEnglishList = ref<English[]>(
    Array.from({ length: pageSize.value }, () => ({ ...defaultEnglish }))
);


const handleEnglishPush = () => {


  axios.post('/api/english/push',{

    english: inputEnglish.value,
    chinese: inputChinese.value
  }).then(response => {
    console.log('english上传成功:', response.data);
    ElMessage.success("上传english成功！")
  }).catch(error => {
    console.error('上传english失败:', error);
  });
}
//
const getEnglish = () => {


  const paginationParams: PaginationParams = {
    page: currentPage.value, // 页码
    pageSize: pageSize.value, // 每页显示的条目数
  };

  axios.get('/api/english/get',{
    params: paginationParams
  }).then(response => {
    console.log('english获取成功:', response.data);
    // 我想要在这里清空一下showEnglishList

    for(let i = 0; i < response.data.data.length; i ++ ) {
      // Dto传回来的是大写
      if (showEnglishList.value[i]) {
        const item = response.data.data[i];
        currentShowSize.value = item.length
        if (item.English !== undefined && item.Chinese !== undefined) {
          showEnglishList.value[i].english = item.English as string;
          showEnglishList.value[i].chinese = item.Chinese as string;
        } else {
          console.error("English or Chinese is undefined for item at index", i);
        }
      } else {
        console.error("showEnglishList.value[", i, "] is undefined");
      }
    }

  }).catch(error => {
    console.error('english获取失败:', error);
  });
}
// 计算属性，返回要展示的数据子集
const paginatedShowEnglishList = computed(() => {
  return showEnglishList.value.slice(0, currentShowSize.value);
});

const getPage = () => {
  axios.get('/api/english/page').then(response => {
    console.log('返回response数据', response.data);
    totalPage.value = (response.data.totalPages as unknown as number) * 10 ;
  }).catch(error => {
    console.error('返回response数据失败: ', error);
  });
}
// 初始化一次
onMounted(() => {
  getPage()
  getEnglish()
})

const handlePageChange = (page: number) => {
  currentPage.value = page
  getEnglish()
}

</script>

<template>

  <el-input
      v-model="inputEnglish"
      style="width: 240px"
      placeholder="Please input"
      clearable
  />

  <el-input
      v-model="inputChinese"
      style="width: 240px"
      placeholder="Please input"
      clearable
  />

  <el-button @click="handleEnglishPush">提交</el-button>



  <div v-for="(englishItem, index) in showEnglishList" :key="index">
    <p class="showEnglish">English: {{ englishItem.english }}</p>
    <p class="showChinese">Chinese: {{ englishItem.chinese }}</p>
  </div>


  <div class="example-pagination-block">
    <div class="example-demonstration"></div>
    <el-pagination layout="prev, pager, next" :total="totalPage"  @current-change="handlePageChange" />
  </div>





</template>

<style scoped lang="scss">

.example-pagination-block + .example-pagination-block {
  margin-top: 10px;
}
.example-pagination-block .example-demonstration {
  margin-bottom: 16px;
}

.showEnglish {
  margin-left: 45%;
  display: inline-block;
}

.showChinese {
  display: inline-block;
  margin-left: 3%;
}
</style>




