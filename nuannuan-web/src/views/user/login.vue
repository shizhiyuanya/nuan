<template>
  <div class="container">
    <div class="login-wrapper">
      <div class="header">用户登录</div>
      <div class="form-wrapper">
        <input type="text" v-model="username" placeholder="请输入用户名" class="input-item" />
        <input type="password" v-model="password" placeholder="请输入密码" class="input-item" />
        <input type="submit" id="login" class="btn" value="登录" @click="login" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/stores/userStore'; // 导入 Pinia store

const username = ref<string>('');
const password = ref<string>('');
// const router = useRouter();
const userStore = useUserStore(); // 使用 Pinia store

const login = async () => {
  try {
    const response = await axios.post('/api/login', {
      Username: username.value,
      Password: password.value,
    });

    if (response.data.success) {
      const user = { id: response.data.id, username: username.value };
      userStore.setCurrentUser(user); // 保存当前用户到 Pinia
      ElMessage.success('登录成功！');
      // router.push('/chat'); // 跳转到聊天页面
    } else {
      ElMessage.error('登录失败：' + response.data.id);
    }
  } catch (error) {
    ElMessage.error('登录失败：网络错误或服务器异常');
    console.error('登录失败:', error);
  }
};
</script>

<style scoped>
/* 样式代码 */
</style>