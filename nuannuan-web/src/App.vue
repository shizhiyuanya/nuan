<script setup lang="ts">
import { RouterView } from 'vue-router'
import "@/styles/container.scss"
import {ref} from "vue";
import {ElMessage, type UploadFile} from 'element-plus'
import { Delete, Download, Plus, ZoomIn } from '@element-plus/icons-vue'
import axios from "axios";
import router from "@/router";

const navigateToAdmin = () => {
  router.push('/admin/picture'); // 使用路由名称
}

const navigateToHome = () => {
  router.push('/home')
}
const navigateToEnglish = () => {
  router.push('/english')
}

const navigateToVideo = () => {
  router.push('/video')
}

const navigateToUser = () => {
  router.push('/user')
}

const richTextEditor = () => {
  router.push('/editor')
}
</script>

<template>


    <div class="common-layout">
      <el-container class="container-style">
        <el-header class="header-style">

          <el-button @click="navigateToHome" class="header-button">首页</el-button>

          <el-button class="header-button">直播</el-button>

          <el-button class="header-button">番剧</el-button>

          <el-button class="header-button" @click="richTextEditor">富文本编辑器</el-button>

          <el-button class="header-button" @click="navigateToEnglish">英语</el-button>

          <el-button class="header-button" @click="navigateToAdmin">后台</el-button>

          <el-button class="header-button" @click="navigateToVideo">视频</el-button>

          <el-button class="header-button" @click="navigateToUser">用户</el-button>

        </el-header>
        <el-container>
          <el-container>
            <el-main>
              <RouterView />
            </el-main>
          </el-container>
        </el-container>
        <el-footer>Footer</el-footer>
      </el-container>

    </div>





</template>

<style scoped lang="scss">

.dialog-style {
  height: 50vh;
  width: 50vw;
  .big-image-style {
    height: 50vh;
    width: 48vw;
    background-size: cover;
  }
  .small-image-style {

  }

}

</style>









