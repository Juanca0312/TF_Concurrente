package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var casos []Caso          //lista de casos
var centroids []Caso      //lista de centroids actuales
var centroids_count []int //range(k), en el centroids_count[0] = # de datos asociados al centroid i
var casos_centroids []int //casos_centroids[i] = en el index i, esta asociado al centroid 0, 1, o 2

type Caso struct {
	//V: Victima
	//A: Agresor
	Mes            float64 `json:"mes"`
	V_Edad         float64 `json:"victimaEdad"`
	V_Numero_Hijos float64 `json:"numeroHijosVictima"`
	V_Embarazo     float64 `json:"embarazoVictima"`
	A_Edad         float64 `json:"edadAgresor"`
	Alcohol        float64 `json:alcohol`
	A_Trabaja      float64 `json:"trabajaAgresor"`
	Medidas        float64 `json:"medidasTomadas"`
	A_Situacion    float64 `json:situacionAgresor`
}

func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func resuelveListar(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	//serualizar, codificar
	jsonBytes, _ := json.MarshalIndent(casos, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func resuelveListarGrupos(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	//serualizar, codificar
	jsonBytes, _ := json.MarshalIndent(centroids_count, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func resuelveFuncion(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	//recuperamos parametros de request
	user_k := req.FormValue("k")
	//tipo de contenido de respuesta
	res.Header().Set("Content-Type", "application/json")

	k, _ := strconv.Atoi(user_k)
	//resetar listas para despues
	centroids = []Caso{}
	casos_centroids = []int{}
	centroids_count = []int{}
	for i := 0; i < len(casos); i++ {
		casos_centroids = append(casos_centroids, 0)
	}
	//llamo a la funcion
	kmeans(k)
	jsonBytes, _ := json.MarshalIndent(casos_centroids, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func manejadorRequest() {
	//definir los endpoints de nestro servicio
	http.HandleFunc("/listar", resuelveListar)
	http.HandleFunc("/funcion", resuelveFuncion)
	http.HandleFunc("/grupos", resuelveListarGrupos)

	//establecer el puesto del servicio
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func selectCentroids(k int) {
	for i := 0; i < k; i++ {
		centroid := rand.Intn(len(casos))
		centroids = append(centroids, casos[centroid])
		centroids_count = append(centroids_count, 0)
	}
}

func euclideanDistance(point1, point2 Caso) float64 {
	sum := math.Pow((float64(point2.Mes)-float64(point1.Mes)), 2) +
		math.Pow((float64(point2.V_Edad)-float64(point1.V_Edad)), 2) +
		math.Pow((float64(point2.V_Numero_Hijos)-float64(point1.V_Numero_Hijos)), 2) +
		math.Pow((float64(point2.V_Embarazo)-float64(point1.V_Embarazo)), 2) +
		math.Pow((float64(point2.A_Edad)-float64(point1.A_Edad)), 2) +
		math.Pow((float64(point2.Alcohol)-float64(point1.Alcohol)), 2) +
		math.Pow((float64(point2.A_Trabaja)-float64(point1.A_Trabaja)), 2) +
		math.Pow((float64(point2.Medidas)-float64(point1.Medidas)), 2) +
		math.Pow((float64(point2.A_Situacion)-float64(point1.A_Situacion)), 2)
	//fmt.Print(sum)
	rpta := math.Sqrt(sum)
	return rpta
}

func setData() {
	url := "https://raw.githubusercontent.com/Juanca0312/concurrente-ta2/main/casos_feminicidio.csv"
	data, err := readCSVFromUrl(url)
	if err != nil {
		panic(err)
	}

	for i, row := range data {
		var caso Caso
		if i != 0 {
			caso.Mes, _ = strconv.ParseFloat(row[0], 64)
			caso.V_Edad, _ = strconv.ParseFloat(row[1], 64)
			caso.V_Numero_Hijos, _ = strconv.ParseFloat(row[2], 64)
			caso.V_Embarazo, _ = strconv.ParseFloat(row[3], 64)
			caso.A_Edad, _ = strconv.ParseFloat(row[4], 64)
			caso.Alcohol, _ = strconv.ParseFloat(row[5], 64)
			caso.A_Trabaja, _ = strconv.ParseFloat(row[6], 64)
			caso.Medidas, _ = strconv.ParseFloat(row[7], 64)
			caso.A_Situacion, _ = strconv.ParseFloat(row[8], 64)

			casos = append(casos, caso)
			casos_centroids = append(casos_centroids, 0)
			//fmt.Print(casos[0].A_Edad)

		}

	}
}

var string_array string = ""

func sumPoints(point1, point2 Caso) Caso {
	var newPointSum Caso
	newPointSum.Mes = point1.Mes + point2.Mes
	newPointSum.V_Edad = point1.V_Edad + point2.V_Edad
	newPointSum.V_Numero_Hijos = point1.V_Numero_Hijos + point2.V_Numero_Hijos
	newPointSum.V_Embarazo = point1.V_Embarazo + point2.V_Embarazo
	newPointSum.A_Edad = point1.A_Edad + point2.A_Edad
	newPointSum.Alcohol = point1.Alcohol + point2.Alcohol
	newPointSum.A_Trabaja = point1.A_Trabaja + point2.A_Trabaja
	newPointSum.Medidas = point1.Medidas + point2.Medidas
	newPointSum.A_Situacion = point1.A_Situacion + point2.A_Situacion
	return newPointSum
}

func asignCentroid() {
	for i, caso := range casos {
		//hallamos la distancia del caso a los 3 centroids
		var d_menor float64
		var c_menor int
		for j, centroid := range centroids {
			if j == 0 {
				d_menor = euclideanDistance(caso, centroid)
				c_menor = j
			} else {
				distance := euclideanDistance(caso, centroid)
				if distance < d_menor {
					d_menor = distance
					c_menor = j
				}
			}
		}
		//ya tenemos el centroid mas cercano al caso actual, luego se lo asignamos
		casos_centroids[i] = c_menor
		centroids_count[c_menor] = centroids_count[c_menor] + 1
	}
}

func divPoints(point Caso, divisor int) Caso {

	point.Mes = point.Mes / float64(divisor)
	point.V_Edad = point.V_Edad / float64(divisor)
	point.V_Numero_Hijos = point.V_Numero_Hijos / float64(divisor)
	point.V_Embarazo = point.V_Embarazo / float64(divisor)
	point.A_Edad = point.A_Edad / float64(divisor)
	point.Alcohol = point.Alcohol / float64(divisor)
	point.A_Trabaja = point.A_Trabaja / float64(divisor)
	point.Medidas = point.Medidas / float64(divisor)
	point.A_Situacion = point.A_Situacion / float64(divisor)
	return point
}

func newCentroids() {
	var media_centroids []Caso
	//inicializamos el
	for i := 0; i < len(centroids_count); i++ {
		var newPointSum Caso
		newPointSum.Mes = 0
		newPointSum.V_Edad = 0
		newPointSum.V_Numero_Hijos = 0
		newPointSum.V_Embarazo = 0
		newPointSum.A_Edad = 0
		newPointSum.Alcohol = 0
		newPointSum.A_Trabaja = 0
		newPointSum.Medidas = 0
		newPointSum.A_Situacion = 0
		media_centroids = append(media_centroids, newPointSum)
	}

	//primero hallamos la suma y luego lo dividimos
	//suma
	for i, caso := range casos {
		media_centroids[casos_centroids[i]] = sumPoints(media_centroids[casos_centroids[i]], caso)

	}
	//dividimos
	for i, centroid := range media_centroids {
		centroids[i] = divPoints(centroid, centroids_count[i])
	}

}

func kmeans(k int) {
	//1 Primer paso: Seleccionar K
	//2 Segundo paso: Seleccionar K centroids(en este caso van a formar parte de nuestros datos)
	selectCentroids(k)
	//3 Tercer paso: Agrupar cada Caso al centroid mas cercano
	//4 Cuarto paso: Hallar la media de los casos agrupados y que sean los nuevos centroids
	//repetir 3 y 4 para convergencia
	for i := 0; i < 10; i++ {
		print("\n")
		print("\n")
		print("EPOCA ")
		print(i + 1)
		print("\n")

		//resetear centroids_count
		for j := 0; j < len(centroids_count); j++ {
			centroids_count[j] = 0
		}
		asignCentroid()
		var auxCentroids []Caso = centroids
		fmt.Print(auxCentroids)
		print("\n")
		fmt.Print(auxCentroids[0].Mes)
		print("\n")
		newCentroids()
		fmt.Print(centroids)
		print("\n")
		fmt.Print(auxCentroids[0].Mes)
		print("\n")
		fmt.Print(centroids[0].Mes)

		// if auxCentroids[0].Mes == centroids[0].Mes {
		// 	print("\n CONVERGENCIA")
		// 	break
		// }
	}
	print("\n")
	fmt.Print(centroids_count)
}

func main() {

	setData()
	selectCentroids(3)
	//convertArrayToString()
	//convertStringToArrays()
	//manejadorRequest()
}
