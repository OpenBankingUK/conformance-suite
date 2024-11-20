<template>
  <div class="login-container">
    <form
      class="login-form"
      @submit.prevent="login">
      <input
        v-model="email"
        type="email"
        placeholder="Email"
        required >
      <input
        v-model="password"
        type="password"
        placeholder="Password"
        required >
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script>
import Cookies from 'js-cookie';
import axios from 'axios';

export default {
  data() {
    return {
      email: '',
      password: '',
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('/api/login', {
          email: this.email,
          password: this.password,
        });
        console.log(response.status);
        console.log(response.data);
        const jwt = response.data.token;
        Cookies.set('ob_jwt', jwt, { expires: 7 });
        this.$router.push('/');
      } catch (error) {
        console.error(error);
        alert('Invalid credentials');
      }
    },
  },
};
</script>

<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: #f5f5f5;
}

.login-form {
    display: flex;
    flex-direction: column;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
    background-color: #fff;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.login-form input {
    margin-bottom: 10px;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 3px;
}

.login-form button {
    padding: 10px;
    border: none;
    border-radius: 3px;
    background-color: #007bff;
    color: #fff;
    cursor: pointer;
}

.login-form button:hover {
    background-color: #0056b3;
}
</style>
