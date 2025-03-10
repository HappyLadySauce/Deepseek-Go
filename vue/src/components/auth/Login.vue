<script setup lang="ts">
import useLogin from '@/hooks/auth/useLogin'

const { form, rules, login } = useLogin()

</script>

<template>
    <div class="login-container">
        <h1>登录</h1>
        <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" @submit.prevent="login">
            <el-form-item label="用户名或者邮箱" prop="username">
                <el-input v-model="form.username" placeholder="请输入用户名或者邮箱"></el-input>
            </el-form-item>
            <el-form-item label="密码" prop="password">
                <el-input v-model="form.password" placeholder="请输入密码"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" native-type="submit" tive-type="submit">登录</el-button> 
            </el-form-item>
        </el-form>
    </div>
</template>



<style scoped>

</style>
