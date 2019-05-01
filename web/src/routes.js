import TheWizard from './views/TheWizard.vue';
import WizardContinueOrStart from './views/Wizard/WizardContinueOrStart.vue';
import WizardImport, { MODES as WizardImportModes } from './views/Wizard/WizardImport.vue';
import WizardDiscoveryConfig from './views/Wizard/WizardDiscoveryConfig.vue';
import WizardConfigurationTabs from './views/Wizard/WizardConfigurationTabs.vue';
import WizardOverviewRun from './views/Wizard/WizardOverviewRun.vue';
import WizardExport from './views/Wizard/WizardExport.vue';
import NotFound from './views/NotFound.vue';
import ConformanceSuiteCallback from './views/ConformanceSuiteCallback/ConformanceSuiteCallback.vue';

export default {
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
          path: '/wizard/import/review',
          name: 'WizardImportReview',
          component: WizardImport,
          props: {
            mode: WizardImportModes.REVIEW,
          },
        },
        {
          path: '/wizard/import/rerun',
          name: 'WizardImportRerun',
          component: WizardImport,
          props: {
            mode: WizardImportModes.RERUN,
          },
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
};
