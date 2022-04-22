<template>
  <div id="content">
    <Donut id="donut" msg="Welcome to Your Vue.js App"/>
    <line-chart id="line"></line-chart>
    <List id="list"></List>
  </div>

</template>

<script>
import Donut from './Donut.vue'
import LineChart from "./LineChart";
import List from "./List";

export default {
  name: "Dashboard",
  components: {
    Donut,
    LineChart,
    List
  }
}
</script>

<style scoped>
@media (min-width: 1000px) {
  #donut {
    width: 30%;
  }

  #line {
    width: 60%;
    margin-top: 10%;
  }

  #list {
    width: 30%;
  }

  #content {
    display: flex;
    flex-basis: 50%;
    flex-direction: row;
    flex-wrap: wrap;
    margin-left: 50px;
    margin-right: 50px;
  }
}

</style>