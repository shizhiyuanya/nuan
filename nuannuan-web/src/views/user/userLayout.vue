<script setup lang="ts">
import router from "@/router";

const navigateToUserLogin = () => {
  router.push('/user/login')
}

const navigateToUserRegister = () => {
  router.push('/user/register')
}

const navigateToUserCommunication = () => {
  router.push('/user/communication')
}

</script>

<template>
  <el-button class="header-button" @click="navigateToUserLogin">用户登录</el-button>
  <el-button class="header-button" @click="navigateToUserRegister">用户注册</el-button>
  <el-button class="header-button" @click="navigateToUserCommunication">用户交谈</el-button>
  <router-view></router-view>
</template>

<style scoped lang="scss">

</style>