import { createRouter, createWebHistory } from 'vue-router'
import AdminView from "@/views/admin/adminPicture.vue";
import HomeView from "@/views/home/home.vue";
import EnglishView from "@/views/english/english.vue";
import VideoView from "@/views/video/video.vue";
import AdminChildren from "@/router/admin.ts"
import AdminLayout from "@/views/admin/adminLayout.vue";
import UserLayout from '@/views/user/userLayout.vue';
import UserChildren from "@/router/user.ts"
import RichTextEditor from "@/views/wangEditor/richTextEditor.vue";
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/home',
      name: 'home',
      component: HomeView,
    },
    {
      path:'/admin',
      name:'admin',
      component: AdminLayout,
      redirect: '/admin/picture',
      // 二级路由
      children: AdminChildren
    },
    {
      path:'/english',
      name:'english',
      component: EnglishView,
    },
    {
      path:'/video',
      name:'video',
      component: VideoView,
    },
    {
      path: '/user',
      name: 'user',
      component: UserLayout,
      redirect: '/user/login',
      children: UserChildren
    },
    {
      path: '/editor',
      name: 'editor',
      component: RichTextEditor,
    },
  ],
})

export default router
