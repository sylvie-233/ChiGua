import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import { permission } from './directives/permission'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(Antd)
app.directive('permission', permission)
app.mount('#app')
