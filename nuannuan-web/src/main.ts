import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import createVuetify from '@/plugins/vuetify'; // 注意这里导入的是函数，不是配置对象

// tailwindcss引入
import '@/style.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(ElementPlus)

const vuetify = createVuetify(app); // 现在这里正确调用了函数
app.use(vuetify);

router.isReady().then(() => app.mount('#app'))
