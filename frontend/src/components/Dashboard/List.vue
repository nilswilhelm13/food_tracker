<template>
  <div v-if="transactions">
    <div style="font-weight: bold">Eaten today</div>
    <b-list-group>
      <b-list-group-item class="listItem" :key="item.id" v-for="item in transactions">
        <div class="content">
          <div>{{ item.foodName }}</div>
          <div>({{ item.amount }}<span v-if="item.food.isMeal"> portion<span v-if="item.amount !== 1">s</span></span><span v-else> g</span>) <span id="cal">{{ item.nutrition.energy }} kcal</span>
          </div>
        </div>
        <div class="deleteBtn" @click="deleteTransaction(item.id)">
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
  </div>

</template>

<script>
import {BListGroup, BListGroupItem} from "bootstrap-vue"
import {mapActions, mapState} from 'vuex'
// import axios from "axios";

export default {
  name: "List",
  components: {
    BListGroup,
    BListGroupItem
  },
  computed: mapState(['transactions']),
  created() {
    this.getTransactions()
  },

  methods: {
    ...mapActions([
      "getTransactions",
      "deleteTransaction"
    ]),
  }
}
</script>

<style scoped>
#cal {
  color: #34ce57;
}

.listItem {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.deleteBtn:hover {
  border: #c43e3e solid 1px;
}

</style>