import Vue from 'vue'
import App from './App.vue'
import {store} from './store/store'
import VueRouter from 'vue-router'
import {routes} from './routes/routes'

Vue.use(VueRouter);

export const router = new VueRouter({
	routes,
	mode: 'history'
})
Vue.config.productionTip = false

new Vue({
  render: h => h(App),
	store,
	router
}).$mount('#app')

import { BootstrapVue } from 'bootstrap-vue'
Vue.use(BootstrapVue)

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'