import Vue from 'vue';
import VueRouter from 'vue-router';
import Login from '../components/Login';
import Landing from '../components/Landing';
import Profile from '../components/Profile';
import Reporter from '../components/Reporter';
import NotFound from '../components/NotFound';
import Config from '../components/Config';
import store from '../store';

Vue.use(VueRouter);

const ifNotAuthenticated = (to, from, next) => {
  if (!store.getters['user/isSignedIn']) return next();
  return next('/');
};

const ifAuthenticated = (to, from, next) => {
  if (store.getters['user/isSignedIn']) return next();
  return next('/login');
};

const router = new VueRouter({
  mode: 'history',
  linkExactActiveClass: 'ant-menu-item-selected',
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: { layout: 'clean' },
      beforeEnter: ifNotAuthenticated,
    },
    {
      path: '/',
      name: 'Landing',
      component: Landing,
      beforeEnter: ifAuthenticated,
    },
    {
      path: '/config',
      name: 'Config',
      component: Config,
      beforeEnter: ifAuthenticated,
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile,
      beforeEnter: ifAuthenticated,
    },
    {
      path: '/reports',
      name: 'Reporter',
      component: Reporter,
      beforeEnter: ifAuthenticated,
    },
    {
      path: '*',
      meta: { layout: 'clean' },
      component: NotFound,
    },
  ],
});

export default router;
