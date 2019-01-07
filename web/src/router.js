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
      name: 'Wizard',
      redirect: '/wizard/continue-or-start',
      component: () => import(/* webpackChunkName: "wizard" */ './components/Wizard'),
      children: [
        {
          path: '/wizard/continue-or-start',
          name: 'ContinueOrStart',
          component: () => import(/* webpackChunkName: "wizard/continue-or-start" */ './components/Wizard/ContinueOrStart'),
        },
        {
          path: '/wizard/discovery-config',
          name: 'DiscoveryConfig',
          component: () => import(/* webpackChunkName: "wizard/discovery-config" */ './components/Wizard/DiscoveryConfig'),
        },
        {
          path: '/wizard/configuration',
          name: 'Configuration',
          component: () => import(/* webpackChunkName: "wizard/configuration" */ './components/Wizard/Configuration'),
        },
        {
          path: '/wizard/run-overview',
          name: 'RunOverview',
          component: () => import(/* webpackChunkName: "wizard/run-overview" */ './components/Wizard/RunOverview'),
        },
        {
          path: '/wizard/summary',
          name: 'Summary',
          component: () => import(/* webpackChunkName: "wizard/summary" */ './components/Wizard/Summary'),
        },
        {
          path: '/wizard/export',
          name: 'Export',
          component: () => import(/* webpackChunkName: "wizard/export" */ './components/Wizard/Export'),
        },
      ],
    },
    {
      path: '*',
      name: 'NotFound',
      component: () => import(/* webpackChunkName: "not-found" */ './components/NotFound'),
    },
  ],
});

export default router;
