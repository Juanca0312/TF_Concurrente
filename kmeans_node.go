package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
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
		casos_centroids[i] = c_menor
		centroids_count[c_menor] = centroids_count[c_menor] + 1
	}
	//fmt.Println("Casos centroids: ", casos_centroids)
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

	string_array += strconv.Itoa(currentIt + 1)
	string_array += "\n!"
	//fmt.Println(string_array)
	return string_array

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
	fmt.Println("Centroids nuevos: ", centroids)
}

func enviar(data string) {
	indice := rand.Intn(len(direcciones))
	fmt.Printf("Enviando data a %s\n", direcciones[indice])
	hostremoto := fmt.Sprintf("%s:%d", direcciones[indice], puerto_procesoHP)
	conn, _ := net.Dial("tcp", hostremoto) //comunicación hacia un nodo. pide protocolo y destino
	defer conn.Close()
	fmt.Fprintf(conn, "%s!", data)
}

func manejadorKmeans(conn net.Conn) {
	defer conn.Close()
	//leer el dato enviado
	bufferIn := bufio.NewReader(conn)
	data, _ := bufferIn.ReadString('!')
	//fmt.Print("Se recibió string de data\n", data)

	//hacer k means, enviar al sig.

	centroids = []Caso{}
	casos_centroids = []int{}
	centroids_count = []int{}
	casos = []Caso{}

	convertStringToArrays(data)

	asignCentroid()
	newCentroids()
	//fmt.Println("Centroids: ", centroids)
	//fmt.Println("Centroids count: ", centroids_count)
	for j := 0; j < len(centroids_count); j++ {
		centroids_count[j] = 0
	}

	for i := 0; i < len(casos_centroids); i++ {
		casos_centroids[i] = 0
	}

	enviarString := convertArrayToString()
	//fmt.Println("Enviar string: ", enviarString)
	if currentIt < 6 {
		enviar(enviarString)
	}
}

var remotehost string
var chCont chan int
var n, min int
var direcciones []string //ips nodos de la red

const (
	puerto_registro  = 8000
	puerto_notifica  = 8001
	puerto_procesoHP = 8002
)

var direccion_nodo string //host del nodo actual

//lógica de servicios

func ManejadorNotificacion(conn net.Conn) {
	defer conn.Close()
	//leer msj enviado
	bufferIn := bufio.NewReader(conn)
	msgIP, _ := bufferIn.ReadString('\n')
	msgIP = strings.TrimSpace(msgIP)
	//agregamos ip de nuevo nodo a bitácora del nodo
	direcciones = append(direcciones, msgIP)
	fmt.Println(direcciones)
}

func AtenderNotificaciones() {
	//modo escucha
	hostlocal := fmt.Sprintf("%s:%d", direccion_nodo, puerto_notifica)
	ln, _ := net.Listen("tcp", hostlocal)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go ManejadorNotificacion(conn)
	}
}

func RegistrarCliente(ip_remoto string) {
	//solicitud a un nodo de la red
	hostremoto := fmt.Sprintf("%s:%d", ip_remoto, puerto_registro) //IP:puerto
	//llamada de conexion al host remoto
	conn, _ := net.Dial("tcp", hostremoto)
	defer conn.Close()
	//enviar ip del nuevo nodo a host rmoto
	fmt.Fprintf(conn, "%s\n", direccion_nodo)
	//recibe la bitácora del host remoto
	bufferIn := bufio.NewReader(conn)
	msgDirecciones, _ := bufferIn.ReadString('\n')

	//decodificar msg recibido
	var auxDirecciones []string
	json.Unmarshal([]byte(msgDirecciones), &auxDirecciones)
	//agregar como parted de la bitácora del nuevo nodo
	direcciones = append(auxDirecciones, ip_remoto) //agregar ip remota a la bitácora
	fmt.Println(direcciones)                        //imprimir la bitácora del nuevo nodo

}

func Notificar(direccion, ip string) {
	//formato del host remoto
	hostremoto := fmt.Sprintf("%s:%d", direccion_nodo, puerto_notifica)
	//estab. conexion con host remoto
	conn, _ := net.Dial("tcp", hostremoto)
	defer conn.Close()
	//enviar el msj ip a nodo remoto
	fmt.Fprintf(conn, "%s\n", ip)
}

func NotificarTodos(ip string) {
	for _, direcc := range direcciones {
		Notificar(direcc, ip)
	}
}

func ManejadorRegistro(conn net.Conn) {
	defer conn.Close()
	//leer ip que llega como msj de la solicitud
	bufferIn := bufio.NewReader(conn)

	msgIP, _ := bufferIn.ReadString('\n')
	msgIP = strings.TrimSpace(msgIP)
	//notificar a todos las ips de la bitácora
	//codificar el msj en formato json
	bytesDirecc, _ := json.Marshal(direcciones)
	//enviar msj de rpta al nuevo nodo con la bitácora actual
	fmt.Fprintf(conn, "%s\n", string(bytesDirecc))

	//enviar a los ips
	NotificarTodos(msgIP)

	//act. su bitácora
	direcciones = append(direcciones, msgIP)
	fmt.Println(direcciones)
}

func AtenderRegistroCliente() {
	//rol de escucha
	hostlocal := fmt.Sprintf("%s:%d", direccion_nodo, puerto_registro)

	//modo escucha
	ln, _ := net.Listen("tcp", hostlocal)
	defer ln.Close()
	//atencion de solicitudes
	for {
		conn, _ := ln.Accept()
		go ManejadorRegistro(conn)
	}
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

func AtenderProcesoHP() {
	//modo escucha
	hostlocal := fmt.Sprintf("%s:%d", direccion_nodo, puerto_procesoHP)
	ln, _ := net.Listen("tcp", hostlocal)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		//manejador kmeans?
		go manejadorKmeans(conn)
	}
}

func main() {
	direccion_nodo = localAddress()
	fmt.Printf("IP: %s\n", direccion_nodo)

	//Rol servidor
	go AtenderRegistroCliente()
	go AtenderProcesoHP()

	//Rol cliente
	//enviar solicitud para unirse a red
	bufferIn := bufio.NewReader(os.Stdin) //ingreso por consola
	fmt.Print("Ingrese IP del host para solicitud: ")
	ip_remoto, _ := bufferIn.ReadString('\n')
	ip_remoto = strings.TrimSpace(ip_remoto)
	if ip_remoto != "" { //solo para nuevos nodos
		RegistrarCliente(ip_remoto)
	}

	//rol de servidor
	AtenderNotificaciones()

	/* bufferIn := bufio.NewReader(os.Stdin)
	fmt.Printf("Ingrese el puerto local: ")
	puerto, _ := bufferIn.ReadString('\n')
	puerto = strings.TrimSpace(puerto)               //elimina espacion de la cadena
	localhost := fmt.Sprintf("localhost:%s", puerto) //IP:Puerto

	//remotehost
	fmt.Print("Ingrese el puerto remoto: ")
	puerto, _ = bufferIn.ReadString('\n')
	puerto = strings.TrimSpace(puerto) //elimina espacion de la cadena
	remotehost = fmt.Sprintf("localhost:%s", puerto) */

	//rol de nodo escucha
	/* ln, _ := net.Listen("tcp", localhost)
	defer ln.Close()
	for {
		//manejador de conexiones
		conn, _ := ln.Accept()
		go manejadorKmeans(conn)
	} */
}
