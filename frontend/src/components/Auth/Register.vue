<template>
    <div>
        <b-form inline id="login">
            <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
                <b-form-input v-model="userData.email" placeholder="email"></b-form-input>
            </b-input-group>
            <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
                <b-form-input v-model="userData.username" placeholder="username"></b-form-input>
            </b-input-group>
            <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
                <b-form-input v-model="userData.password" placeholder="password"></b-form-input>
            </b-input-group>


            <b-button @click="register" variant="primary">Create Account</b-button>
        </b-form>
    </div>
</template>

<script>
    import {BForm, BInputGroup, BFormInput, BButton} from "bootstrap-vue"
    import axios from "../../store/axios_instance";

    export default {
        name: "Register",
        data: function () {
            return {
                userData: {
                    email: "",
                    password: "",
                    username:""
                }
            }
        },
        components: {
            BFormInput,
            BInputGroup,
            BForm,
            BButton
        },
        methods: {
            register: function () {
                console.log(this.userData)
                axios.post("/register", this.userData).then(res => {
                    console.log(res)
                    // const now = new Date();
                    // const expirationDate = new Date(
                    //     now.getTime() + res.data.expiresIn * 1000
                    // );
                    // localStorage.setItem("token", res.data.token);
                    // localStorage.setItem("userId", res.data.userId);
                    // localStorage.setItem("expirationDate", expirationDate);
                })
            }
        }
    }
</script>

<style scoped>
    #login {
        margin-top: 100px;

    }

    @media (min-width: 1000px) {
        #login {
            margin-left: 20px;

        }
    }

</style>