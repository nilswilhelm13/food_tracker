<template>
  <div>
    <b-form inline id="login">
      <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
        <b-form-input v-model="loginData.email" placeholder="email"></b-form-input>
      </b-input-group>
      <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
        <b-form-input v-model="loginData.password" placeholder="password"></b-form-input>
      </b-input-group>


      <b-button @click="login(loginData)" variant="primary">Login</b-button>
    </b-form>
  </div>
</template>

<script>
import {BButton, BForm, BFormInput, BInputGroup} from "bootstrap-vue"
import {mapActions} from 'vuex'

export default {
  name: "Login",
  data: function () {
    return {
      loginData: {
        email: "",
        password: ""
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
    ...mapActions([
      "login",
    ]),
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