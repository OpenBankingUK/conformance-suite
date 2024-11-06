// src/auth.js
import Cookies from 'js-cookie';

export function isAuthenticated() {
  return !!Cookies.get('ob_jwt');
}

export function logout() {
  Cookies.remove('ob_jwt');
}