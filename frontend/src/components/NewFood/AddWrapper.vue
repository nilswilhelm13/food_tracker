<template>
  <div>
    <b-tabs v-model="tabIndex">
      <b-tab title="Eat">
        <eat @newFood="tabIndex=1"></eat>
      </b-tab>
      <b-tab title="New Food">
        <new-food></new-food>
      </b-tab>
      <b-tab title="New Meal">
       <new-meal @newFood="tabIndex=1"></new-meal>
      </b-tab>
    </b-tabs>

  </div>

</template>

<script>
import {BTab, BTabs} from "bootstrap-vue"
import Eat from "@/components/NewFood/Eat";
import NewFood from "@/components/NewFood/NewFood";
import NewMeal from "@/components/NewFood/NewMeal";

export default {
name: "AddWrapper",
  components:{
    NewMeal,
    NewFood,
    Eat,
    BTab,
    BTabs
  },
  data: ()=> {
    return {
      tabIndex: 0,
      showModal: false,
      message: ""
    }
  }
}
</script>

<style scoped>

</style>