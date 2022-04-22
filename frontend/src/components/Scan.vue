<template>
  <div style="width: 100%">
    <div style="width: 100%; overflow: hidden;">
      <div>Code: {{ code }}</div>
      <v-quagga class="videoInsert" :onDetected="logIt" :readerSize="{
        width: { max: 100 },
         height: {  max: 100 }
      }" :readerTypes="['ean_reader']"></v-quagga>
    </div>

  </div>

</template>

<script>
import Vue from "vue"
import VueQuagga from 'vue-quaggajs';

Vue.use(VueQuagga)
export default {
  name: 'Scan',
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
.videoInsert {
  position: absolute;
  right: 0;
  bottom: 0;
  min-width: 100%;
  min-height: 100%;
  width: auto;
  height: auto;
  z-index: -100;
  background-size: cover;
  overflow: hidden;
}
</style>
