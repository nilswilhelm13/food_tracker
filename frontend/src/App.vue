<template>
  <div id="app">
    <Nav id="navigation"></Nav>
    <Modal></Modal>
    <router-view id="content">
      <Login></Login>
      <Register></Register>
      <Dashboard id="dashboard"></Dashboard>
      <AddWrapper></AddWrapper>
      <FoodInput></FoodInput>
    </router-view>


  </div>
</template>

<script>

import Nav from "./components/Nav";
import Dashboard from "./components/Dashboard/Dashboard";
import Login from "./components/Auth/Login";
import Register from "./components/Auth/Register";
import Modal from "@/components/Modal";
import FoodInput from "@/components/NewFood/Eat";
import AddWrapper from "@/components/NewFood/AddWrapper";


export default {
  name: 'App',
  components: {
    AddWrapper,
    FoodInput,
    Modal,
    Register,
    Login,
    Dashboard,
    Nav,
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 0;
  width: 100%;
  padding: 0;
}

#navigation {
  width: 100%;
  position: fixed;
  top: 0;
  z-index: 1;
}

#content {
  margin-top: 54px;
  padding-top: 10px;
}

</style>
