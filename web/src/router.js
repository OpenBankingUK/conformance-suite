import Vue from 'vue';
import VueRouter from 'vue-router';
import has from 'lodash/has';
import store from './store/';

import Wizard from './components/Wizard.vue';
import ContinueOrStart from './components/Wizard/ContinueOrStart.vue';
import DiscoveryConfig from './components/Wizard/DiscoveryConfig.vue';
import ConfigurationTabs from './components/Wizard/ConfigurationTabs.vue';
import RunOverview from './components/Wizard/RunOverview.vue';
import Summary from './components/Wizard/Summary.vue';
import Export from './components/Wizard/Export.vue';
import NotFound from './components/NotFound.vue';

Vue.use(VueRouter);

const router = new VueRouter({
  // Use the HTML5 history API, so that routes look normal
  // (e.g. `/about`) instead of using a hash (e.g. `/#/about`).
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'Wizard',
      redirect: '/wizard/continue-or-start',
      component: Wizard,
      children: [
        {
          path: '/wizard/continue-or-start',
          name: 'ContinueOrStart',
          component: ContinueOrStart,
        },
        {
          path: '/wizard/discovery-config',
          name: 'DiscoveryConfig',
          component: DiscoveryConfig,
        },
        {
          path: '/wizard/configuration',
          name: 'Configuration',
          component: ConfigurationTabs,
        },
        {
          path: '/wizard/run-overview',
          name: 'RunOverview',
          component: RunOverview,
        },
        {
          path: '/wizard/summary',
          name: 'Summary',
          component: Summary,
        },
        {
          path: '/wizard/export',
          name: 'Export',
          component: Export,
        },
      ],
    },
    // ---
    // Handle 404s
    // ---
    {
      path: '/404',
      name: '404',
      component: NotFound,
    },
    {
      path: '*',
      redirect: '404',
    },
  ],
});

/**
 * Prevent access to a route if the previous step has not been completed.
 * Example: If an attempt to go to the `/wizard/run-overview` route is made whilst the current, `step`,
 * is `1` we redirect to landing page (`/`). This tends to happen when the User refreshes the page.
 */
router.beforeEach((to, from, next) => {
  const { path } = to;
  const navigation = store.getters['config/navigation'];

  // If it is not a wizard-related path (e.g., `/404`), ignore it.
  if (!has(navigation, path)) {
    return next();
  }

  const viewable = navigation[path];
  if (viewable) {
    return next();
  }

  return next('/');
});

export default router;
