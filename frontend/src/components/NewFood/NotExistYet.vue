<template>
  <div class="notExist">
    <b-card>
      The food with EAN: {{ ean }} does not exist in our database. <br>
      Do you want to create it?
      <b-button-group>
        <b-button variant="success" @click="yes">Yes</b-button>
        <b-button variant="danger" @click="no">No</b-button>
      </b-button-group>
    </b-card>

  </div>
</template>

<script>
import {BButton, BButtonGroup, BCard} from "bootstrap-vue"
import {mapState} from "vuex"

export default {
  name: "NotExistYet",
  components: {
    BButton,
    BButtonGroup,
    BCard
  },
  methods: {
    no() {
      this.$emit('no')
    },
    yes() {
      this.$emit('yes')
      console.log('yes')
    }
  },
  computed: {
    ...mapState(['ean'])
  }
}
</script>

<style scoped>

</style>