import Vue from 'vue'
import Router from 'vue-router'
import RequestLog from '@/components/RequestLog'
import BootstrapVue from 'bootstrap-vue'

Vue.use(Router)
Vue.use(BootstrapVue)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'RequestLog',
      component: RequestLog
    }
  ]
})
