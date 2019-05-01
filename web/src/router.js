import Vue from 'vue';
import VueRouter from 'vue-router';
import has from 'lodash/has';
import store from './store';
import routes from './routes';

Vue.use(VueRouter);

const router = new VueRouter(routes);

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
