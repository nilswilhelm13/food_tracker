<template>
  <div>
    <div class="form" v-if="!scan">

      <div class="inputContainer">
        <div>Name</div>
        <b-form-input style="width: 300px" v-model="food.name" placeholder="Enter Name of food"></b-form-input>
      </div>
      <div class="inputContainer">
        <div>EAN</div>
        <div style="display: flex">
          <b-form-input style="width: 260px" v-model="ean" placeholder="EAN if available"></b-form-input>
          <div id="barcode" @click="scan=true">
            <barcode></barcode>
          </div>
        </div>
      </div>
      <div class="inputContainer">
        <div>Calories</div>
        <b-form-input class="amount" v-model="food.nutrition.energy"></b-form-input>
      </div>
      <div class="inputContainer">
        <div>Carbs</div>
        <b-form-input class="amount" v-model="food.nutrition.carbohydrate"></b-form-input>
      </div>
      <div class="inputContainer">
        <div>Protein</div>
        <b-form-input class="amount" v-model="food.nutrition.protein"></b-form-input>
      </div>
      <div class="inputContainer">
        <div>Fat</div>
        <b-form-input class="amount" v-model="food.nutrition.fat"></b-form-input>
      </div>




      <b-button variant="success" class="bt" @click="postFood">Submit</b-button>


    </div>

    <scan v-if="scan" @scanned="onScan"></scan>
    <div v-if="showModal">Pimmel</div>
  </div>
</template>

<script>
import {BButton, BFormInput} from "bootstrap-vue"
import axios from "../../store/axios_instance";
import {mapMutations} from "vuex";
import Scan from "@/components/Scan";
import Barcode from "@/components/Icons/Barcode";

export default {
  name: "NewFood",
  components: {
    BFormInput,
    BButton,
    Scan,
    Barcode,


  },
  props: {
    id: String
  },
  data: function () {
    return {
      food: {
        id: this.id,
        name: "",
        nutrition:
            {
              energy: 0,
              fat: 0,
              carbohydrate: 0,
              protein: 0,
            }
      },
      scan: false,
      showModal: false
    }
  },
  computed: {
    ean () {
      return this.$store.state.ean
    }

  },
  methods: {
    ...mapMutations([
      'setResponse',
      'setError',
        "setEAN"
    ]),
    postFood: function () {
      const idToken = localStorage.getItem("token")
      const queryParams =
          "?auth=" + idToken
      let data = {
        id: this.id,
        name: this.food.name,
        ean: this.ean.trim(),
        nutrition:
            {
              energy: parseFloat(this.food.nutrition.energy),
              fat: parseFloat(this.food.nutrition.fat),
              carbohydrate: parseFloat(this.food.nutrition.carbohydrate),
              protein: parseFloat(this.food.nutrition.protein),
            }
      }
      axios.post("/foodlist/1" + queryParams, data).then(res => {
        console.log(res)
        this.setResponse(res)
        // Clear EAN in store
        this.setEAN("")

      }).catch(err => {
        this.setError(err)
      })
    },
    onScan(barcode) {
      this.showModal = true
      this.scan = false
      this.food.ean = barcode
    }
  }
}
</script>

<style scoped>
.bt {
  margin: 10px;

}

#barcode{
  margin-left: 5px;
}
.amount{
  width: 70px
}
.inputContainer{
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 5px;
}
.form{
  margin: 10px;
}

</style>