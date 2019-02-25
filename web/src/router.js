import Vue from 'vue';
import VueRouter from 'vue-router';
import has from 'lodash/has';
import store from './store';

import TheWizard from './views/TheWizard.vue';
import WizardContinueOrStart from './views/Wizard/WizardContinueOrStart.vue';
import WizardDiscoveryConfig from './views/Wizard/WizardDiscoveryConfig.vue';
import WizardConfigurationTabs from './views/Wizard/WizardConfigurationTabs.vue';
import WizardOverviewRun from './views/Wizard/WizardOverviewRun.vue';
import WizardExport from './views/Wizard/WizardExport.vue';
import NotFound from './views/NotFound.vue';
import ConformanceSuiteCallback from './views/ConformanceSuiteCallback/ConformanceSuiteCallback.vue';

Vue.use(VueRouter);

const router = new VueRouter({
  // Use the HTML5 history API, so that routes look normal
  // (e.g. `/about`) instead of using a hash (e.g. `/#/about`).
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'TheWizard',
      redirect: '/wizard/continue-or-start',
      component: TheWizard,
      children: [
        {
          path: '/wizard/continue-or-start',
          name: 'WizardContinueOrStart',
          component: WizardContinueOrStart,
        },
        {
          path: '/wizard/discovery-config',
          name: 'WizardDiscoveryConfig',
          component: WizardDiscoveryConfig,
        },
        {
          path: '/wizard/configuration',
          name: 'WizardConfigurationTabs',
          component: WizardConfigurationTabs,
        },
        {
          path: '/wizard/overview-run',
          name: 'WizardOverviewRun',
          component: WizardOverviewRun,
        },
        {
          path: '/wizard/export',
          name: 'WizardExport',
          component: WizardExport,
        },
      ],
    },
    {
      path: '/conformancesuite/callback',
      name: 'ConformanceSuiteCallback',
      component: ConformanceSuiteCallback,
    },
    // ---
    // Handle 404s
    // ---
    {
      path: '/404',
      name: 'NotFound',
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
 * Example: If an attempt to go to the `/wizard/overview-run` route is made whilst the current, `step`,
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
