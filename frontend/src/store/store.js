import Vue from "vue";
import Vuex from "vuex";
import axios from "../store/axios_instance";
import {router} from "@/main"
import moment from "moment";
import VuexPersist from 'vuex-persist';


const vuexLocalStorage = new VuexPersist({
    key: 'vuex', // The key to store the state on in the storage provider.
    storage: window.localStorage, // or window.sessionStorage or localForage
    // Function that passes the state and returns the state with only the objects you want to store.
    // reducer: state => state,
    // Function that passes a mutation and lets you decide if it should update the state in localStorage.
    // filter: mutation => (true)
})


Vue.use(Vuex);

export const store = new Vuex.Store({
    plugins: [vuexLocalStorage.plugin],
    state: {
        intake: {
            carbohydrate: 0,
            protein: 0,
            fat: 0,
        },
        intake_list: [],
        transactions: [],
        isAuthorized: false,
        token: null,
        userId: null,
        expirationDate: null,
        scanFood: null,
        donutSeries: [],
        response: null,
        error: null,
        ean: ""
    },
    getters: {
        getIntake: state => {
            return state.intake;
        },
        getintake_list: state => {
            return state.intake_list;
        },
    },
    mutations: {
        setIntakeList: (state, data) => {
            state.intake_list = data;
            console.log(state.intake_list[0].Date)
        },
        setIntake: (state, data) => {
            let intake = data.nutrition
            state.intake = intake;
            state.donutSeries = [intake.carbohydrate, intake.protein, intake.fat]
        },
        setTransactions: (state, data) => {
            state.transactions = data;
            console.log("State")
            console.log(state.transactions)
        },
        authenticate: (state, data) => {
            state.token = data.token
            state.userId = data.userId
            state.expirationDate = data.expirationDate
            state.isAuthorized = true
            localStorage.setItem("token", data.token);
            localStorage.setItem("userId", data.userId);
            localStorage.setItem("expirationDate", data.expirationDate);
            localStorage.setItem("isAuthorized", "true");

        },
        deauthenticate: (state) => {
            state.token = null
            state.userId = null
            state.expirationDate = null
            state.isAuthorized = false
            localStorage.setItem("token", null);
            localStorage.setItem("userId", null);
            localStorage.setItem("expirationDate", null);
            localStorage.setItem("isAuthorized", null);
        },
        setSearchFood: (state, food) => {
            state.scanFood = food
        },
        setResponse: (state, response) => {
            state.response = response
        },
        setError: (state, error) => {
            state.error = error
        },
        setEAN: (state, ean) => {
            state.ean = ean
        },
    },
    //---------------------------------------------
    //------------ACTIONS--------------------------
    //---------------------------------------------
    actions: {
        fetchIntakeHistory: ({commit}) => {
            axios
                .get("/history")
                .then(res => {
                    commit("setIntakeList", res.data);
                }).catch(() => {
                commit('setError', "Could not fetch intake history")
            })
        },
        getIntake: ({commit}) => {
            axios
                .get("/intake")
                .then(res => {
                    commit("setIntake", res.data)
                }).catch(() => {
                commit('setError', "Could not fetch intake")
            })
        },
        getTransactions: ({commit}) => {
            const userId = localStorage.getItem("userId")
            axios
                .get("/transactions/" + userId)
                .then(res => {
                        commit("setTransactions", res.data)
                    }
                ).catch(() => {
                commit('setError', "Could not fetch transactions")
            })
        },
        deleteTransaction: ({commit}, id) => {
            axios
                .delete("/intake/" + id).then(res => {
                    store.dispatch("getTransactions")
                    store.dispatch("getIntake")
                    store.dispatch("")
                    console.log(res)
                }
            ).catch(() => {
                commit('setError', "Could not delete transaction")
            })
            console.log(commit)

        },
        login: ({commit}, loginData) => {
            axios.post("/login", loginData).then(res => {
                console.log(res)

                const expirationDate = moment().add(res.data.expiresIn,
                    'minutes');
                let authData = {
                    token: res.data.token,
                    userId: res.data.userId,
                    expirationDate: expirationDate
                }
                commit('authenticate', authData)
                router.push("/")
            }).catch(() => {
                commit('setError', "Login failed")
            })
        },
        logout: ({commit}) => {
            console.log("logout")
            commit("deauthenticate")
        },
        autoAuth: ({commit}) => {
            let token = localStorage.getItem("token")
            let userId = localStorage.getItem("userId")
            let expirationDate = moment(localStorage.getItem("expirationDate"))
            console.log(expirationDate)
            if (expirationDate > moment() && token !== "" && userId !== "") {
                commit("authenticate", {
                    token: token,
                    userId: userId,
                    expirationDate: expirationDate
                })
            } else {
                commit("deauthenticate")
            }
        },

    },
    fillIntake: ({commit}) => {
        axios
            .get("/intake")
            .then(res => {
                commit("setIntake", res.data)
            }).catch(() => {
            commit('setError', "Could not get intake")
        })

    },
});