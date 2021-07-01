<template>
  <div id="app">
    <h1>Feminicidio en el Perú con kmeans</h1>
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
        <div>Feminicidio en el Perú con KMeans</div>
        <b-form-group
          id="input-group-1"
          label="Ingrese Datos a predecir"
          label-for="input-1"
          description="K representa el numero de clusters a agrupar"
        >
          <b-form-input
            id="input-1"
            v-model="mes"
            placeholder="0 - Enero, 12 - Diciembre"
            required
          ></b-form-input>
          <b-form-input
            id="input-1"
            v-model="edadVictim"
            placeholder="ej: 35"
            required
          ></b-form-input>
          <b-form-input
            id="input-1"
            v-model="hijosVictim"
            placeholder="ej: 3"
            required
          ></b-form-input>
          <b-form-input
            id="input-1"
            v-model="embaVictim"
            placeholder="0 - No Embarazada, 1 - Embarazada"
            required
          ></b-form-input>
          <b-form-input
            id="input-1"
            v-model="edadAgresor"
            placeholder="ej: 38"
            required
          ></b-form-input>
          <b-form-input
            id="input-1"
            v-model="drogasAlcoh"
            placeholder="0 - No drogas o alcohol en feminicidio, 1 - Si"
            required
          ></b-form-input>
          <b-form-input
            id="input-1"
            v-model="trabajaAgresor"
            placeholder="0 - No trabaja, 1 - Trabaja"
            required
          ></b-form-input>
          <b-form-input
            id="input-1"
            v-model="medidasTomadas"
            placeholder="ej: 4"
            required
          ></b-form-input>
          <b-form-input
            id="input-1"
            v-model="situacionAgresor"
            placeholder="ej: 1"
            required
          ></b-form-input>
        </b-form-group>
        <b-button class="mr-1" @click="algoritmo"
          >Predecir a que grupo pertenece</b-button
        >
        <b-button class="mr-1" @click="reset">Reset</b-button>
        <h2>Tus datos pertencen al grupo {{this.clusterPrediction}}</h2>
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
    mes: null,
    edadVictim: null,
    hijosVictim: null,
    embaVictim: null,
    edadAgresor: null,
    drogasAlcoh: null,
    trabajaAgresor: null,
    medidasTomadas: null,
    situacionAgresor: null,
    grupos: [],
    centroids: [],
    clusterPrediction: 0,
  }),
  mounted() {
    axios.get("http://localhost:9000/listar").then((response) => {
      this.datos = response.data;
      //console.log(this.datos);
    });
    axios.get("http://localhost:9000/centroids").then((response) => {
      this.centroids = response.data;
      //console.log(this.centroids);
    });
    axios.get("http://localhost:9000/grupos").then((response) => {
      this.grupos = response.data;
      //console.log(this.grupos);
    });
    axios.get("http://localhost:9000/casosCentroids").then((response) => {
      this.clusters = response.data;
      //console.log(this.grupos);
    });
  },
  methods: {
    algoritmo() {
      var distances = [];
      for (var i = 0; i < this.centroids.length; i++) {
        var sum =
          Math.pow(this.mes - this.centroids[i].mes, 2) +
          Math.pow(this.edadVictim - this.centroids[i].victimaEdad, 2) +
          Math.pow(this.hijosVictim - this.centroids[i].numeroHijosVictima, 2) +
          Math.pow(this.embaVictim - this.centroids[i].embarazoVictima, 2) +
          Math.pow(this.edadAgresor - this.centroids[i].edadAgresor, 2) +
          Math.pow(this.drogasAlcoh - this.centroids[i].Alcohol, 2) +
          Math.pow(this.trabajaAgresor - this.centroids[i].trabajaAgresor, 2) +
          Math.pow(this.medidasTomadas - this.centroids[i].medidasTomadas, 2) +
          Math.pow(this.situacionAgresor - this.centroids[i].A_Situacion, 2);
        var distance = Math.sqrt(sum);
        distances.push(distance);
      }

      var rpta = distances[0];
      for (var j = 0; j < distances.length; j++) {
        if (distances[j] < rpta) {
          rpta = distances[j];
          this.clusterPrediction = j;
        }
      }
    },
    reset() {
      this.mes = null;
      this.edadVictim = null;
      this.hijosVictim = null;
      this.embaVictim = null;
      this.edadAgresor = null;
      this.drogasAlcoh = null;
      this.trabajaAgresor = null;
      this.medidasTomadas = null;
      this.situacionAgresor = null;
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
  color: #2c3e50;
  margin: 60px;
  width: 100%;
}
.buton {
  margin: 10;
}
b-tab {
  width: 100%;
}
</style>
