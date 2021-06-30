package main

import "math"

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
		//casos_centroids[i] = c_menor
		//centroids_count[c_menor] = centroids_count[c_menor] + 1
	}
}

func main() {

}
