<template>
    <div class="container">
      <div class="login-wrapper">
        <div class="header">
          用户登录
        </div>
        <!--表单包裹-->
        <div class="form-wrapper">
          <!--用户名-->
          <input type="text" v-model="username" placeholder="请输入用户名" class="input-item">
          <!--密码-->
          <input type="password" v-model="telephone" placeholder="请输入手机号码" class="input-item">

          <input type="password" v-model="password" placeholder="请输入密码" class="input-item">
          <!--登录-->
          <input type="submit" id="login" class="btn" value="注册" @click="register">
        </div>
        <!-- <div class="msg"> 已有账号？<a href="">立即登录</a> </div> -->
  
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue'
  import { useRouter } from "vue-router";

  import axios from "axios";
  const username = ref<string>()
  const password = ref<string>()
  const telephone = ref<string>()
  const router = useRouter()
  // 获取url 可以获取的到
  import { ElMessage } from 'element-plus';

  const register = () => {
    console.log("用户注册：", username, password)


  
    // 向后端发送登录请求
    axios.post('/api/register', {
    Username: username.value,
    Telephone: telephone.value,
    Password: password.value
    }).then(response => {
    console.log('注册成功:', response.data);
    ElMessage.success("注册成功！")
    }).catch(error => {
    ElMessage.error("注册失败！")
    console.error('注册失败:', error);
    });
  }
  </script>
  
  <style scoped lang="scss">
  // @import "style/login.scss";
  </style>