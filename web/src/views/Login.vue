<template>
    <div class="login-container">
        <form @submit.prevent="login" class="login-form">
            <input type="email" v-model="email" placeholder="Email" required />
            <input type="password" v-model="password" placeholder="Password" required />
            <button type="submit">Login</button>
        </form>
    </div>
</template>

<script>
import Cookies from 'js-cookie';

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
                // const response = await axios.post('/api/login', {
                //   email: this.email,
                //   password: this.password,
                // });
                // const jwt = response.data.jwt;
                Cookies.set('ob_jwt', this.email, { expires: 7 });
                this.$router.push('/');
            } catch (error) {
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