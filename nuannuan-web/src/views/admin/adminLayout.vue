<script setup lang="ts">
import router from "@/router";

const navigateToPictureAdmin = () => {
  router.push('/admin/picture')
}

const navigateToVideoAdmin = () => {
  router.push('/admin/video')
}

</script>

<template>
  <el-button class="header-button" @click="navigateToPictureAdmin">图片提交</el-button>
  <el-button class="header-button" @click="navigateToVideoAdmin">视频提交</el-button>
  <router-view></router-view>
</template>

<style scoped lang="scss">

</style>