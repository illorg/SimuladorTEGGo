package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var fichas_ataque int = 0    // fichas que se jugaran en ataque
var fichas_defensa int = 0   // fichas qie se jugaran en defensa
var simulaciones int = 10000 // cantidad de simulaciones. por defecto 100000
var vict_ataque int = 0      // contador de veces que gana el ataque
var vict_defensa int = 0     // contador de veces que gana la defensa
var proceso_porct int = 0    // porcentaje (nro entero) de avance de la simulacion
var hilosCompletos [4]bool   // bandera que indica la finalizacion de las 4 gorutines

func simula(sim int, id int) { //funcion simula llamada 4 veces por goroutines
	a := 0
	d := 0
	for simulacion := 1; a+d < sim; simulacion++ { // Ciclo for de simulacciones : por defect 10K
		if Jugar(fichas_ataque, fichas_defensa) { // llama funcion jugar, envia cant fichas, devuelve ganador
			a += 1

		} else {
			d += 1
		}
		if id == 1 && (float64(simulacion)/float64(sim)*100) > float64(proceso_porct+1) { // actualiza porcentaje de avance( por ahora solo hilo nro1)
			proceso_porct = int(float64(simulacion) / float64(sim) * 100)
		}

	}
	vict_ataque = vict_ataque + a
	vict_defensa = vict_defensa + d
	fmt.Println("goroutine ID:", id, "realizo operaciones", (a + d))
	hilosCompletos[id] = true

}
func Jugar(fatq int, fdef int) bool { // Funcion jugar se ejecuta hasta que se acaben las fichas
	cant_d_ataque := 0
	cant_d_defensa := 0
	lanzamiento := [2][]int{{0, 0, 0}, {0, 0, 0}}
	comparar := 3

	for (fatq > 1) && (fdef > 0) { // Juega hasta que se acaben fichas
		if fatq > 3 { // defino cant dados de ataque
			cant_d_ataque = 3
		} else if fatq == 3 {
			cant_d_ataque = 2
		} else {
			cant_d_ataque = 1
		}

		if fdef >= 3 { // defino cant dados defen
			cant_d_defensa = 3
		} else {
			cant_d_defensa = fdef
		}
		lanzamiento[0], lanzamiento[1] = (tirar_dados(cant_d_ataque, cant_d_defensa)) // llama funcion tirar dados
		if cant_d_ataque <= cant_d_defensa {                                          // elige cuantos dados se comparan
			comparar = cant_d_ataque
		} else {
			comparar = cant_d_defensa
		}
		for comp := 0; comp < comparar; comp++ {
			if int(lanzamiento[0][comp]) > int(lanzamiento[1][comp]) {
				fdef -= 1
			} else {
				fatq -= 1
			}
		}
	}
	if fdef == 0 {
		return true
	} else {
		return false
	}
}
func tirar_dados(cant_d_ataque int, cant_d_defensa int) ([]int, []int) {
	rand.Seed(time.Now().UnixNano())
	dados_ataque := []int{0, 0, 0}
	dados_defensa := []int{0, 0, 0}
	for x := 0; x <= 3; x++ {
		if (x + 1) <= cant_d_ataque {
			dados_ataque[x] = rand.Intn(6) + 1
		}
		if (x + 1) <= cant_d_defensa {
			dados_defensa[x] = rand.Intn(6) + 1
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(dados_ataque)))
	sort.Sort(sort.Reverse(sort.IntSlice(dados_defensa)))
	return dados_ataque, dados_defensa
}
func main() {
	/// leer argumentos  en el llamado

	if len(os.Args[1:]) > 0 {
		fichas_ataque, _ = strconv.Atoi(os.Args[1])
		fichas_defensa, _ = strconv.Atoi(os.Args[2])
		simulaciones, _ = strconv.Atoi(os.Args[3])
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Ingrese Fichas Ataque")
		scanner.Scan()
		fichas_ataque, _ = strconv.Atoi(scanner.Text())
		fmt.Println("Ingrese Fichas Defensa")
		scanner.Scan()
		fichas_defensa, _ = strconv.Atoi(scanner.Text())
	}

	fmt.Println(fichas_ataque)
	fmt.Println(fichas_defensa)
	comienzo := time.Now()

	/// Creo 4 Goroutines que se reparten el nro total de simulaciones///
	for i := 0; i < 4; i++ {
		go simula(simulaciones/4.0, i)
	}
	/// Bucle para controlar el avance de los Goroutines///
	for !hilosCompletos[0] || !hilosCompletos[1] || !hilosCompletos[2] || !hilosCompletos[3] {

		fmt.Println("Simulando: %", proceso_porct)
		time.Sleep(250 * time.Millisecond)

	}
	//// Reporta resultados en pantalla////
	porcent_vict := float32(float64(vict_ataque) / float64(simulaciones) * 100.0)
	porcent_derr := float32(float64(vict_defensa) / float64(simulaciones) * 100.0)
	fmt.Println("Tiempo consumido en el cálculo : ")
	tiempofinal := time.Now()
	fmt.Println("Victoria ataque: ", vict_ataque, " %", porcent_vict)
	fmt.Println("victoria defensa: ", vict_defensa, " %", porcent_derr)
	fmt.Println("el calculo tardó: ", tiempofinal.Sub(comienzo))

	///////////// BASE DE DATOS //////////////////
	// 1) creo objeto para conectarme
	db, err := sql.Open("mysql", "root:pcshoprg@tcp(179.62.90.59:3306)/SimuladorTeg")

	// manejo si hay error
	if err != nil {
		panic(err)
	}

	// comando Ping returns error, si la base no responde.
	err = db.Ping()

	// manejar si hay error
	if err != nil {
		panic(err)
	}

	fmt.Print("Conexion establecida con Base de datos en MistralHome!!")

	insertarRegistros, err := db.Prepare("INSERT INTO registro (fecha, simulaciones, fichas_ataque, fichas_defensa, porct_victoria_ataque, porct_victoria_defensa) VALUES(?,?,?,?,?,?)")
	// manejar error
	if err != nil {
		panic(err)
	}
	insertarRegistros.Exec(time.Now(), simulaciones, fichas_ataque, fichas_defensa, porcent_vict, porcent_derr)

	// cerrar y liberar base de datos
	if err == nil {
		defer fmt.Println("Base cerrada exitosamente", time.Now())
	}
	defer db.Close()
}
