import Vue from 'vue'
import Router from 'vue-router'
import FileTranfer from '@/components/FileTranfer'


Vue.use(Router)

export default new Router({
  routes: [
    // {
    //   path: '/',
    //   name: 'HelloWorld',
    //   component: HelloWorld
    // },
    {
      path: '/',
      name: 'FileTranfer',
      component: FileTranfer
    }
  ]
})
