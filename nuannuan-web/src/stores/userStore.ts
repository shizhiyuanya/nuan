import { defineStore } from 'pinia';

// 定义用户信息的类型
interface User {
  id: string;
  username: string;
}

// 创建全局状态管理
export const useUserStore = defineStore('user', {
  state: () => ({
    currentUser: null as User | null, // 当前登录的用户
    // users: [] as User[], // 所有登录的用户（用于测试多用户）
  }),
  actions: {
    // 设置当前登录用户
    setCurrentUser(user: User) {
      this.currentUser = user;
    },
    // 添加用户到用户列表
    // addUser(user: User) {
    //   this.users.push(user);
    // },
    // 根据ID获取用户
    // getUserById(id: string): User | undefined {
    //   return this.users.find((user) => user.id === id);
    // },
  },
});