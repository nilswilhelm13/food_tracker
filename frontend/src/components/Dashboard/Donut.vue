<template>

  <div v-if="intake.energy > 0" id="chart">
    <apexchart type="donut" :options="chartOptions" :series="series"></apexchart>
    <div id="energy">
      <div>Energy</div>
      <div>{{ intake.energy }} kcal</div>
    </div>
  </div>

</template>

<script>
import VueApexCharts from 'vue-apexcharts'
import {mapActions} from "vuex";

export default {
  name: 'Donut',
  props: {
    msg: String
  },
  created() {
    this.fillIntake()
  },
  computed: {
    intake() {
      return this.$store.state.intake
    },
    series() {
      return this.$store.state.donutSeries
    }
  },
  methods: {
    ...mapActions([
      "fillIntake",
    ]),

  },
  data: function () {
    return {
      // intake: {},
      // series: [],
      chartOptions: {
        colors: ['#f09263', '#59bdc6', '#e9e795', '#5a2a27'],
        chart: {
          width: '100%',
          type: 'donut',
        },
        labels: ["Carbs", "Protein", "Fat"],
        theme: {
          monochrome: {
            enabled: false
          }
        },
        plotOptions: {
          pie: {
            dataLabels: {
              style: {
                fontSize: '30px',
                fontFamily: 'Helvetica, Arial, sans-serif',
                fontWeight: 'bold',
                colors: undefined
              },
              offset: -5
            }
          }
        },
        title: {},
        // dataLabels: {
        //   formatter(val, opts) {
        //     const name = opts.w.globals.labels[opts.seriesIndex]
        //     return [name, val.toFixed(1) + "%"]
        //   }
        // },
        dataLabels: {
          enabled: true, formatter: function (val, opt) {
            return opt.w.config.series[opt.seriesIndex].toFixed() + " " + opt.w.globals.labels[opt.seriesIndex]
          }
        },
        legend: {
          show: false
        }
      },

    }
  },
  components: {
    'apexchart': VueApexCharts
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

h3 {
  margin: 40px 0 0;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}

#energy {
  position: absolute;
  top: 50%;
  left: 50%;

  transform: translate(-50%, -50%);

}

@media (min-width: 1000px) {
  #energy {
    font-size: 2rem;

  }
}


</style>
