<template>
  <div id="app">
    <b-tabs content-class="mt-3">
      <b-tab title="Dataset" active>
        <p v-for="(item, index) in grupos" v-bind:key="index">
          El cluster {{ index }} hay {{ item }} elementos
        </p>
        <table>
          <tr>
            <th>Mes</th>
            <th>Edad Victima</th>
            <th># Hijos Victima</th>
            <th>Embarazos Victima</th>
            <th>Edad Agresor</th>
            <th>Drogas/Alcohol</th>
            <th>Agresor Trabaja</th>
            <th>Medidas Tomadas</th>
            <th>Situacion Agresor</th>
          </tr>
          <tr
            v-for="(dato, index) in datos"
            v-bind:key="index"
            :style="{
              background:
                clusters[index] == 1
                  ? '#7F5FDD'
                  : clusters[index] == 2
                  ? '#FFD764'
                  : clusters[index] == 3
                  ? '#E0584F'
                  : clusters[index] == 4
                  ? '#98B755'
                  : 'white',
            }"
          >
            <td>{{ dato.mes }}</td>
            <td>{{ dato.victimaEdad }}</td>
            <td>{{ dato.numeroHijosVictima }}</td>
            <td>{{ dato.embarazoVictima }}</td>
            <td>{{ dato.edadAgresor }}</td>
            <td>{{ dato.Alcohol }}</td>
            <td>{{ dato.trabajaAgresor }}</td>
            <td>{{ dato.medidasTomadas }}</td>
            <td>{{ dato.A_Situacion }}</td>
          </tr>
        </table>
      </b-tab>
      <b-tab title="Algoritmo">
        <div>Kmeans Algorithm</div>
        <b-form-group
          id="input-group-1"
          label="Ingrese K"
          label-for="input-1"
          description="K representa el numero de clusters a agrupar"
        >
          <b-form-input
            id="input-1"
            v-model="cluster"
            placeholder="ej: 3"
            required
          ></b-form-input>
        </b-form-group>
        <b-button class="mr-1" @click="algoritmo">Empezar KMeans</b-button>
        <b-button class="mr-1" @click="reset">Reset</b-button>
        <p v-for="(item, index) in clusters" v-bind:key="index">
          El dato {{ index }} pertenece al cluster {{ item }}
        </p>
      </b-tab>
    </b-tabs>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data: () => ({
    datos: [],
    clusters: [],
    cluster: "",
    grupos: [],
  }),
  mounted() {
    axios.get("http://localhost:9000/listar").then((response) => {
      this.datos = response.data;
      //console.log(this.datos);
    });
  },
  methods: {
    algoritmo() {
      console.log(this.cluster);
      axios
        .get(`http://localhost:9000/funcion?k=${this.cluster}`)
        .then((response) => {
          this.clusters = response.data;
          //console.log(this.clusters);
          axios.get("http://localhost:9000/grupos").then((response) => {
            this.grupos = response.data;
            console.log(this.grupos);
          });
        });
    },
    reset() {
      this.clusters = [];
      this.grupos = [];
    },
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  align-content: center;
  justify-content: center;
  display: flex;
  color: #2c3e50;
  margin: 60px;
}
.buton {
  margin: 10;
}
</style>
