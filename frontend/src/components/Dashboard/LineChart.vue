<template>
  <div>
    <apexchart width="100%" type="line" :options="chartOptions" :series="series"></apexchart>
  </div>
</template>

<script>
import VueApexCharts from 'vue-apexcharts'
import axios from "../../store/axios_instance";
import moment from "moment";

export default {
  computed: {},
  created() {
    this.fetchIntakeHistory()
  },
  mounted() {
    this.$forceUpdate();
  },
  components: {
    'apexchart': VueApexCharts
  },
  name: "LineChart",
  methods: {

    fetchIntakeHistory: function () {
      axios
          .get("/history")
          .then(res => {
            this.intake_list = res.data
            console.log(this.intake_list)
            this.populateDays()
          })
    },
    populateDays: function () {
      let days = []
      let values = []


      this.intake_list.reverse().forEach(function (d) {
        let day = moment(d.date).format('ddd');
        days.push(day)
        values.push(Math.round(d.nutrition.energy) || 0)
      })
      this.series = [{data: values}]
      this.chartOptions = {
        ...this.chartOptions, ...{
          xaxis: {
            categories: days
          }
        }
      }
    }
  },
  data: function () {
    return {
      intake_list: [],
      series: [
        {
          name: "Calories",
          data: []
        }
      ],
      chartOptions: {

        chart: {
          width: "100%",
          height: 350,
          type: 'line',
          zoom: {
            enabled: false
          },
          toolbar: {
            show: false
          }
        },
        dataLabels: {
          enabled: false
        },
        stroke: {
          curve: 'straight'
        },
        title: {
          text: 'Calorie History',
          align: 'left'
        },
        grid: {
          row: {
            colors: ['#f3f3f3', 'transparent'], // takes an array which will be repeated on columns
            opacity: 0.5
          },
        },
        xaxis: {
          categories: []
        },
        yaxis: {
          labels: {
            show: false
          }
        },
      },


    }
  },
}
</script>

<style scoped>

</style>