<template>
  <div>
    <b-tabs>
      <b-tab title="Login">
        <login></login>
      </b-tab>
      <b-tab title="Register">
        <Register></Register>
      </b-tab>
    </b-tabs>
  </div>

</template>

<script>
import {BTab, BTabs} from "bootstrap-vue"
import Login from "@/components/Auth/Login";
import Register from "@/components/Auth/Register";

export default {
  name: "Auth",
  components: {Register, Login, BTabs, BTab}
}
</script>

<style scoped>

</style>