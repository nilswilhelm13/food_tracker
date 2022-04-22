<template>
  <div>Logging out ...</div>
</template>

<script>
import {mapActions} from 'vuex'
export default {
name: "Logout",
  created() {
    this.logout()
    this.$router.push('login')
  },
  methods: {
    ...mapActions([
      "logout"
    ]),
  }
}
</script>

<style scoped>

</style>