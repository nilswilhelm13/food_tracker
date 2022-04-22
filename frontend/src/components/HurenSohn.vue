<template>
  <div style="width: 100%">


  </div>

</template>

<script>
import Vue from "vue"
import VueQuagga from 'vue-quaggajs';

Vue.use(VueQuagga)
export default {
  name: 'HurenSohn',
  data: () => {
    return {
      code: "",
      detecteds: [],
      w: 0,
      h: 0,
    }
  },
  created() {
    this.w = window.innerWidth;
    this.h = window.innerHeight;
  },
  mounted() {
    this.logIt({
      codeResult: {
        code: "123456789"
      }
    })
  },
  components: {},
  methods: {
    logIt(data) {
      console.log('detected', data)
      this.code = data.codeResult.code
      this.$emit("scanned", this.code)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
/*video {*/
/*  width: 100%;*/
/*}*/

</style>
