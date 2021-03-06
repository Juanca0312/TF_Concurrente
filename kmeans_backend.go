package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

func resuelveCentroids(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	//serualizar, codificar
	jsonBytes, _ := json.MarshalIndent(centroids, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func resuelveCasosCentroids(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	//serualizar, codificar
	jsonBytes, _ := json.MarshalIndent(casos_centroids, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func manejadorRequest() {
	//definir los endpoints de nestro servicio
	http.HandleFunc("/listar", resuelveListar)
	http.HandleFunc("/centroids", resuelveCentroids)
	http.HandleFunc("/grupos", resuelveListarGrupos)
	http.HandleFunc("/casosCentroids", resuelveCasosCentroids)

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

func convertArrayToString() string { //encodign
	//lista de casos, cada espacio es un column, cada /n es una fila
	var string_array = ""
	for _, row := range casos {
		string_array += fmt.Sprintf("%f", row.Mes)
		string_array += " "
		string_array += fmt.Sprintf("%f", row.V_Edad)
		string_array += " "
		string_array += fmt.Sprintf("%f", row.V_Numero_Hijos)
		string_array += " "
		string_array += fmt.Sprintf("%f", row.V_Embarazo)
		string_array += " "
		string_array += fmt.Sprintf("%f", row.A_Edad)
		string_array += " "
		string_array += fmt.Sprintf("%f", row.Alcohol)
		string_array += " "
		string_array += fmt.Sprintf("%f", row.A_Trabaja)
		string_array += " "
		string_array += fmt.Sprintf("%f", row.Medidas)
		string_array += " "
		string_array += fmt.Sprintf("%f", row.A_Situacion)
		string_array += "\n"
	}
	string_array += "end\n"
	//casos_centroids: recibiendo los centroids asociados por cada item
	for _, item := range casos_centroids {
		string_array += strconv.Itoa(item)
		string_array += " "
	}
	string_array += "\nend\n"

	//centroids_count: cantidad de [12,12,45] significa 12 asignados al centroid 0, 12 asignados al centroid 1, 45 asignados al centroid  2

	for _, item := range centroids_count {
		string_array += strconv.Itoa(item)
		string_array += " "
	}
	string_array += "\nend\n"

	//centroids: k centroids
	for _, item := range centroids {
		string_array += fmt.Sprintf("%f", item.Mes)
		string_array += " "
		string_array += fmt.Sprintf("%f", item.V_Edad)
		string_array += " "
		string_array += fmt.Sprintf("%f", item.V_Numero_Hijos)
		string_array += " "
		string_array += fmt.Sprintf("%f", item.V_Embarazo)
		string_array += " "
		string_array += fmt.Sprintf("%f", item.A_Edad)
		string_array += " "
		string_array += fmt.Sprintf("%f", item.Alcohol)
		string_array += " "
		string_array += fmt.Sprintf("%f", item.A_Trabaja)
		string_array += " "
		string_array += fmt.Sprintf("%f", item.Medidas)
		string_array += " "
		string_array += fmt.Sprintf("%f", item.A_Situacion)
		string_array += "\n"

	}
	string_array += "end\n"
	//current iteration
	string_array += "0"
	string_array += "!"
	//fmt.Println(string_array)
	return string_array

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
		//hallamos la distancia del caso a los k centroids
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

var currentIt int

func convertStringToArrays(string_array string) { //decoding
	//fmt.Println("\n\n\n STRING TO ARRAY \n\n\n")
	var contEnds = 0
	//recibiendo casos
	spliteado := strings.Split(string_array, "\n")
	for i := 0; i < len(spliteado); i++ {
		//fmt.Println(spliteado[i])
		//fmt.Println("\n -------------- \n")

		if spliteado[i] == "end" {
			contEnds++
			continue
		}

		if contEnds == 0 { //primer grupo solo para arreglo casos
			split2 := strings.Split(spliteado[i], " ")
			var caso Caso
			caso.Mes, _ = strconv.ParseFloat(split2[0], 64)
			caso.V_Edad, _ = strconv.ParseFloat(split2[1], 64)
			caso.V_Numero_Hijos, _ = strconv.ParseFloat(split2[2], 64)
			caso.V_Embarazo, _ = strconv.ParseFloat(split2[3], 64)
			caso.A_Edad, _ = strconv.ParseFloat(split2[4], 64)
			caso.Alcohol, _ = strconv.ParseFloat(split2[5], 64)
			caso.A_Trabaja, _ = strconv.ParseFloat(split2[6], 64)
			caso.Medidas, _ = strconv.ParseFloat(split2[7], 64)
			caso.A_Situacion, _ = strconv.ParseFloat(split2[8], 64)

			casos = append(casos, caso)
		}
		if contEnds == 1 { //segundo grupo casos_centroids
			split2 := strings.Split(spliteado[i], " ")
			for j := 0; j < len(split2)-1; j++ {
				inti, _ := strconv.Atoi(split2[j])
				casos_centroids = append(casos_centroids, inti)
			}

		}
		if contEnds == 2 {

			split2 := strings.Split(spliteado[i], " ")
			//print(len(split2))
			for j := 0; j < len(split2)-1; j++ {
				inti, _ := strconv.Atoi(split2[j])
				centroids_count = append(centroids_count, inti)
			}
		}
		if contEnds == 3 {
			if spliteado[i] != "!" {
				//println("LONGITUD DE SPLITEADO: ", spliteado[i])
				split2 := strings.Split(spliteado[i], " ")
				var caso Caso
				caso.Mes, _ = strconv.ParseFloat(split2[0], 64)
				caso.V_Edad, _ = strconv.ParseFloat(split2[1], 64)
				caso.V_Numero_Hijos, _ = strconv.ParseFloat(split2[2], 64)
				caso.V_Embarazo, _ = strconv.ParseFloat(split2[3], 64)
				caso.A_Edad, _ = strconv.ParseFloat(split2[4], 64)
				caso.Alcohol, _ = strconv.ParseFloat(split2[5], 64)
				caso.A_Trabaja, _ = strconv.ParseFloat(split2[6], 64)
				caso.Medidas, _ = strconv.ParseFloat(split2[7], 64)
				caso.A_Situacion, _ = strconv.ParseFloat(split2[8], 64)

				centroids = append(centroids, caso)
			}
		}
		if contEnds == 4 {
			if spliteado[i] != "!" {
				it, _ := strconv.Atoi(spliteado[i])
				currentIt = it
				fmt.Println("Current iteracion: ", currentIt)
			}
		}

	}

	/* fmt.Print(casos)
	fmt.Print("\n\n")
	fmt.Print(casos_centroids)
	fmt.Print("\n\n")
	fmt.Print(centroids_count)
	fmt.Print("\n\n")
	fmt.Print(centroids) */
}

var remotehost string

func enviar(data string) {
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	fmt.Fprintf(conn, "%s!", data)
}

func AtenderProcesoHP() {
	//modo escucha
	hostlocal := fmt.Sprintf("%s:%d", direccion_nodo, 8002)
	ln, _ := net.Listen("tcp", hostlocal)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		//manejador kmeans?

		go manejadorKmeans(conn)
		manejadorRequest()

	}
}

func manejadorKmeans(conn net.Conn) {
	defer conn.Close()
	//leer el dato enviado
	bufferIn := bufio.NewReader(conn)
	data, _ := bufferIn.ReadString('!')
	//fmt.Print("Se recibi?? string de data\n", data)

	//hacer k means, enviar al sig.

	centroids = []Caso{}
	casos_centroids = []int{}
	centroids_count = []int{}
	casos = []Caso{}

	convertStringToArrays(data)

	asignCentroid()
	newCentroids()
	fmt.Println("Centroids finales: ", centroids)
	fmt.Println("Grupos: ", centroids_count)
}

var direccion_nodo string

func main() {
	rand.Seed(time.Now().UnixNano())
	direccion_nodo = localAddress()

	setData()
	selectCentroids(3)
	prueba := convertArrayToString()
	//print(prueba)

	BufferIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el IP remoto: ")
	IP, _ := BufferIn.ReadString('\n')
	IP = strings.TrimSpace(IP)
	remotehost = fmt.Sprintf("%s:%d", IP, 8002)

	enviar(prueba)
	AtenderProcesoHP()

}

func localAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("localAddress: %v\n", err.Error()))
		return "127.0.0.1"
	}
	for _, oiface := range ifaces {
		if strings.HasPrefix(oiface.Name, "Wi-Fi") {
			addrs, err := oiface.Addrs()
			if err != nil {
				log.Print(fmt.Errorf("localAddress: %v\n", err.Error()))
				continue
			}
			for _, dir := range addrs {
				switch d := dir.(type) {
				case *net.IPNet:
					if strings.HasPrefix(d.IP.String(), "192") {
						return d.IP.String()
					}

				}
			}
		}
	}
	return "127.0.0.1"
}
