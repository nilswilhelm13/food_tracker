<template>
  <div>
    <div class="modal" v-if="response" @click="setResponse(null)">
      <b-alert variant="primary">
        {{ response.data }}
      </b-alert>
    </div>
    <div class="modal" v-if="error" @click="setError(null)">
      <b-alert show variant="danger">
        {{ error }}
      </b-alert>
    </div>
  </div>

</template>

<script>
import {mapMutations, mapState} from 'vuex'
import {BAlert} from "bootstrap-vue"

export default {
  name: "Modal",
  computed: mapState(['response', 'error']),
  methods: {
    ...mapMutations([
      'setResponse',
      'setError'
    ]),
  },
  components: {
    BAlert
  }

}
</script>

<style scoped>
.modal {
  height: 100%;
  width: 100vw;
  background-color: rgba(39, 79, 39, 0.49);
  position: fixed;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
</style>