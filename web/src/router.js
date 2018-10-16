import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  linkExactActiveClass: 'ant-menu-item-selected',
  routes: [
    {
      path: '/',
      name: 'Landing',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "landing" */ './components/Landing'),
    },
    {
      path: '/config',
      name: 'Config',
      component: () => import(/* webpackChunkName: "config" */ './components/Config'),
    },
    {
      path: '/reports',
      name: 'Reporter',
      component: () => import(/* webpackChunkName: "reporter" */ './components/Reporter'),
    },
    {
      path: '*',
      meta: { layout: 'clean' },
      component: () => import(/* webpackChunkName: "not-found" */ './components/NotFound'),
    },
  ],
});

export default router;
