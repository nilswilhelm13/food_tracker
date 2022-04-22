<template>
  <div>
    <b-navbar toggleable="lg" type="dark" variant="info">
      <b-navbar-brand href="#">FoodTracker</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item v-if="!isAuthorized" href="/login">Login</b-nav-item>
          <b-nav-item v-if="isAuthorized" href="/">Dashboard</b-nav-item>
          <b-nav-item v-if="isAuthorized" href="/addfood">Add Food</b-nav-item>
          <b-nav-item v-if="isAuthorized" href="/logout">Logout</b-nav-item>
        </b-navbar-nav>


      </b-collapse>
    </b-navbar>
  </div>
</template>

<script>
import {BCollapse, BNavbar, BNavbarBrand, BNavbarNav, BNavbarToggle, BNavItem} from 'bootstrap-vue'
import {mapActions} from "vuex";

export default {
  name: "Nav",
  components: {
    BNavbar,
    BNavbarBrand,
    BNavbarToggle,
    BCollapse,
    BNavbarNav,
    BNavItem

  },
  created() {
    this.autoAuth()
  },
  computed: {
    isAuthorized: function () {
      return this.$store.state.isAuthorized
    }
  },
  methods: {
    ...mapActions([
      "logout",
      "autoAuth"
    ]),
  }

}
</script>

<style scoped>

</style>