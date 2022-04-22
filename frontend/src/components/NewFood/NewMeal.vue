<template>
  <div class="buildMeal">
    <div v-if="!scan">
      <search id="search" @startScan="scan=true" @onChoose="addIngredient" @scanned="onScan"></search>
      <div>Ingredients:</div>
      <b-list-group v-if="!scan">
        <b-list-group-item class="ingredient" :key="ingredient.food.id" v-for="ingredient in ingredients">
          {{ ingredient.food.name }}
          <b-form-input @change="calculate" style="width: 70px" v-model="ingredient.amount" placeholder="Amount"></b-form-input>
          <div class="deleteBtn" @click="removeIngredient(ingredient.food.id)">
            <svg width="1em" height="1em" viewBox="0 0 16 16" class="bi bi-trash" fill="currentColor"
                 xmlns="http://www.w3.org/2000/svg">
              <path
                  d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/>
              <path fill-rule="evenodd"
                    d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4L4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/>
            </svg>
          </div>
        </b-list-group-item>
      </b-list-group>
      <div class="meal" v-if="!scan">
        <div class="col" style="flex: 6 0">
          <div class="row">
            <b-form-input v-model="name" placeholder="Enter Name of Meal"></b-form-input>
          </div>
          <div class="row">
            <div class="macro">Energy <b-badge variant="success">{{energy.toFixed(2)}}kcal</b-badge></div>
            <div class="macro">Carbs <b-badge variant="success">{{carbs.toFixed(2)}}g</b-badge></div>
            <div class="macro">Protein <b-badge variant="success">{{protein.toFixed(2)}}g</b-badge></div>
            <div class="macro">Fat <b-badge variant="success">{{fat.toFixed(2)}}g</b-badge></div>
          </div>
        </div>
        <div class="col" style="flex: 1 0">
          <b-button id="bt" @click="postMeal">Post</b-button>
        </div>
      </div>
    </div>
    <scan v-if="scan" @scanned="onScan"></scan>
    <not-exist-yet id="notExist" v-if="showModal" @yes="goToNewFood" @no="showModal=false"></not-exist-yet>
  </div>
</template>

<script>

import {BListGroup, BListGroupItem, BBadge, BButton, BFormInput} from "bootstrap-vue"
import Search from "@/components/Search";
import axios from "@/store/axios_instance"
import {mapActions, mapMutations, mapState} from "vuex";
import Scan from "@/components/Scan";
import NotExistYet from "@/components/NewFood/NotExistYet";
export default {
  name: "NewMeal",
  components: {
    NotExistYet,
    BListGroup,
    BListGroupItem,
    BBadge,
    Search,
    BButton,
    BFormInput,
    Scan

  },
  computed: mapState(['scanFood']),
  data: () => {
    return {
      ingredients: [],
      energy: 0,
      carbs: 0,
      protein: 0,
      fat: 0,
      name: "",
      scan: false,
      showModal: false
    }
  },
  methods: {
    ...mapMutations([
      'setResponse',
      'setError'
    ]),
    ...mapActions([
      'getFoodByEAN'
    ]),
    addIngredient(result) {
      let ingredient = {food: result}
      this.ingredients.push(ingredient)
      console.log(this.ingredients)
    },
    calculate(){
      let countE = 0
      let countC = 0
      let countP = 0
      let countF = 0
      this.ingredients.forEach(ingridient => {
        countE += ingridient.food.nutrition.energy * ingridient.amount / 100
        countC += ingridient.food.nutrition.carbohydrate * ingridient.amount / 100
        countP += ingridient.food.nutrition.protein * ingridient.amount / 100
        countF += ingridient.food.nutrition.fat * ingridient.amount / 100
      })
      this.energy = countE
      this.carbs = countC
      this.protein = countP
      this.fat = countF
    },
    removeIngredient(id){
      for(let i = 0; i < this.ingredients.length; i++) {
        if(this.ingredients[i].food.id === id) {
          this.ingredients.splice(i, 1)
          this.calculate()
          return
        }
      }
    },
    postMeal(){
      let  list = []
      this.ingredients.forEach(ingredient => {
        let i = {
          food: {
            id: ingredient.food.id
          },
          amount: parseFloat(ingredient.amount)
        }
        list.push(i)
      })
      let data = {
        ingredients: list,
        name: this.name
      }
      axios.post("/meal/1", data).then(res => {
        this.setResponse(res)
      }).catch(err => {
        this.setError(err)
      })
    },
    onScan(barcode) {
      axios.get("foodlist/" + barcode + "?ean=true").then(res => {
        if (res.status === 200) {
          let ingredient = {food: res.data}
          this.ingredients.push(ingredient)
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
      this.showModal = false
      this.$emit('newFood')
    },
  }
}
</script>

<style scoped>
  .ingredient{
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .meal{
    display: flex;
    border: 1px solid lightgray;
    justify-content: space-between;
    align-items: center;
    margin: 10px 0;
    border-radius: 5px;
    padding: 5px;
    width: 95%;
    background-color: white;
  }
  .macro{
    font-size: 1.1em;
    display: flex;
    flex-direction: column;
    margin: 0 10px;
  }
  .buildMeal{
    margin: 10px;
  }
  #search {
    margin-bottom: 10px;
  }
  .row{
    display: flex;
    flex-direction: row;
  }
  .col{
    display: flex;
    flex-direction: column;
  }
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
</style>