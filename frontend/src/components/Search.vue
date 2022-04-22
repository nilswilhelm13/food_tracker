<template>
  <div>
    <div v-if="!scan">
      <div style="display: flex">
        <b-form-input v-model="searchTerm" @change="search" placeholder="Search Food"></b-form-input>
        <div @click="scan=startScan" style="margin-left: 5px">
          <barcode></barcode>
        </div>

      </div>
      <div style="text-align: left" v-if="searchResults.length > 0">Search Results:</div>
      <b-list-group class="list">
        <b-list-group-item :key="result.id" @click="onChoose(result)" v-for="result in searchResults">
          {{ result.name }}
        </b-list-group-item>
      </b-list-group>
    </div>
    <camera-view v-if="scan" @scanned="onScan"></camera-view>
  </div>

</template>

<script>
import axios from "@/store/axios_instance";

import {BFormInput, BListGroup, BListGroupItem,} from 'bootstrap-vue'
import Barcode from "@/components/Icons/Barcode";
import CameraView from "@/components/Scan";
import {mapActions} from "vuex"

export default {
  name: "Search",
  components: {
    CameraView,
    Barcode,
    BListGroupItem,
    BListGroup,
    BFormInput,

  },
  data: () => {
    return {
      searchResults: [],
      searchTerm: null,
      scan: false,
    }
  },
  methods: {
    ...mapActions(['getFoodByEAN']),
    search() {
      if (this.searchTerm) {
        axios.get("/search/" + this.searchTerm).then(res => {
          this.searchResults = res.data
          console.log(res)
        })
      }
    },
    onChoose(result) {
      console.log(result)
      this.$emit('onChoose', result)
    },
    onScan(barcode) {
      this.$emit('scanned', barcode)
      this.scan = false
    },
    startScan() {
      this.scan = true
      this.$emit('startScan')
    },
    getFoodByEAN: ({commit}, id) => {
      axios.get("foodlist/" + id + "?ean=true").then(res => {
        if (res.status === 200) {
          commit("setSearchFood", res.data)
          this.$emit('scanned', res.data)
        } else {
          this.$emit('newFood', id)
        }
        // commit('setError', "Product with EAN:" + id + " not found")
      }).catch(() => {
        commit('setError', "Could not get product with EAN:" + id)
      })
    }
  }
}
</script>

<style scoped>
.list {
  overflow: hidden;
  overflow-y: scroll;
}
</style>