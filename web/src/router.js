import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  // linkExactActiveClass: 'ant-menu-item-selected',
  // linkActiveClass: 'active', // active class for non-exact links.
  // linkExactActiveClass: 'active', // active class for *exact* links.
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
      path: '/wizard',
      name: 'Wizard',
      component: () => import(/* webpackChunkName: "wizard" */ './components/Wizard'),
    },
    // {
    //   path: '/wizard/step1',
    //   name: 'Step1',
    //   component: () => import(/* webpackChunkName: "step1" */ './components/Wizard/Step1'),
    // },
    // {
    //   path: '/wizard/step2',
    //   name: 'Step2',
    //   component: () => import(/* webpackChunkName: "step2" */ './components/Wizard/Step2'),
    // },
    // {
    //   path: '/wizard/step3',
    //   name: 'Step3',
    //   component: () => import(/* webpackChunkName: "step3" */ './components/Wizard/Step3'),
    // },
    // {
    //   path: '/wizard/step4',
    //   name: 'Step4',
    //   component: () => import(/* webpackChunkName: "step4" */ './components/Wizard/Step4'),
    // },
    // {
    //   path: '/wizard/step5',
    //   name: 'Step5',
    //   component: () => import(/* webpackChunkName: "step5" */ './components/Wizard/Step5'),
    // },
    {
      path: '*',
      meta: { layout: 'clean' },
      component: () => import(/* webpackChunkName: "not-found" */ './components/NotFound'),
    },
  ],
});

export default router;
