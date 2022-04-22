<template>
  <div>
    <div style="margin: 10px;" v-if="!scan">
      <b-button @click="onScan2('12345')"></b-button>
      <search class="search" @startScan="scan=true" @onChoose="setSearchFood" @scanned="onScan2"></search>
      <div>
        <div style="text-align: left">Selected:</div>

        <b-card style="background-color: rgba(143,255,120,0.42)" v-if="scanFood" @click="showDetails=!showDetails">

          <div>{{ scanFood.name }}</div>
          <div id="details" v-if="showDetails">
            <div v-if="scanFood.ean">EAN: {{ scanFood.ean }}</div>
            <div v-if="scanFood.isMeal">Nutrients per portion:</div>
            <div v-else>Nutrients per 100 grams:</div>
            <div>Energy: {{ scanFood.nutrition.energy.toFixed(2) }} kcal</div>
            <div>Protein: {{ scanFood.nutrition.protein.toFixed(2) }} g</div>
            <div>Carbs: {{ scanFood.nutrition.carbohydrate.toFixed(2) }} g</div>
            <div>Fat: {{ scanFood.nutrition.fat.toFixed(2) }} g</div>
          </div>

        </b-card>
      </div>

      <div v-if="scanFood">
        <label v-if="!scanFood.isMeal" for="amount">Amount in g/ml</label>
        <label v-else for="amount">Portions</label>
      </div>

      <b-form-input type="number" v-model="amount" id="amount"></b-form-input>
      <b-form-datepicker id="example-datepicker" v-model="date" class="mb-2"></b-form-datepicker>
      <b-button-group>
        <b-button variant="success" @click="lookupFood">Submit Meal</b-button>
      </b-button-group>

    </div>
    <not-exist-yet id="notExist" v-if="showModal" @yes="goToNewFood" @no="showModal=false"></not-exist-yet>
    <div v-if="showModal">pimmel</div>

  </div>


</template>

<script>
import {BButton, BButtonGroup, BCard, BFormDatepicker, BFormInput} from 'bootstrap-vue'

import Search from "@/components/Search";
import axios from "../../store/axios_instance";
import moment from "moment"
import {mapMutations} from 'vuex'
import NotExistYet from "@/components/NewFood/NotExistYet";



export default {
  name: "Eat",
  created() {
    this.date = moment().format()
  },
  data: () => {
    return {
      foodID: 0,
      amount: 0,
      showModal: false,
      date: "",
      showDetails: false,
      scan: false
    }
  },
  components: {


    NotExistYet,
    BFormInput,
    BButton,
    BButtonGroup,
    BFormDatepicker,
    BCard,
    Search,
  },
  computed: {
    scanFood() {
      return this.$store.state.scanFood
    }
  },
  methods: {
    ...mapMutations([
      'setSearchFood',
      'setResponse',
      'setError',
      'setEAN'
    ]),
    lookupFood() {
      axios.get("/foodlist/" + this.scanFood.id).then(res => {
        if (res.status !== 200) {
          this.showModal = true
        } else {
          this.postIntake()
        }
      }).catch(() => {
        this.setError("Could not get food")
      })
    },
    postIntake() {
      let data = {
        foodId: this.scanFood.id,
        amount: parseInt(this.amount),
        date: new Date(this.date)
      }
      console.log(data)
      axios.post("/intake", data).then(res => {
        this.setResponse(res)
      }).catch(() => {
        this.setError("Could not post food to intake")
      })
    },
    onScan2(barcode) {

      axios.get("foodlist/" + barcode + "?ean=true").then(res => {
        if (res.status === 200) {
          this.setSearchFood(res.data)
        } else {
          this.showModal = true
          this.setEAN(barcode)
        }
      }).catch(() => {
        this.showModal = true
        this.setEAN(barcode)
      })
    },
    goToNewFood() {
      console.log('pimmel')
      this.showModal = false
      this.$emit('newFood')
    },
  }
}
</script>
<style scoped>

#notExist {
  height: 100%;
  width: 100vw;
  background-color: rgba(39, 79, 39, 0.49);
  position: fixed;
  top: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.search {
}
</style>