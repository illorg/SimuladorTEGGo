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

var fichas_ataque int = 0
var fichas_defensa int = 0
var simulaciones int = 1000000
var vict_ataque int = 0
var vict_defensa int = 0
var proceso_porct int = 0

//muestra = []
//var victorias_ant int = 0
//var simulaciones_ant int = 0

func simula(sim int) {
	for simulacion := 1; simulacion <= sim; simulacion++ { // Ciclo for de simulacciones : por defect 10K
		if Jugar(fichas_ataque, fichas_defensa) { // llama funcion jugar, envia cant fichas, devuelve ganador
			vict_ataque += 1
		} else {
			vict_defensa += 1
		}

		//muestreo(simulaciones)
	}
}
func Jugar(fatq int, fdef int) bool {
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

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Ingrese Fichas Ataque")
	scanner.Scan()
	fichas_ataque, _ = strconv.Atoi(scanner.Text())
	//fichas_ataque = 12
	fmt.Println("Ingrese Fichas Defensa")
	scanner.Scan()
	fichas_defensa, _ = strconv.Atoi(scanner.Text())
	//fichas_defensa = 6 //temporal . sacar!!

	fmt.Println(fichas_ataque)
	fmt.Println(fichas_defensa)
	comienzo := time.Now()
	for i := 0; i < 4; i++ {
		go simula(simulaciones / 4.0)
	}
	for (vict_ataque + vict_defensa) < (simulaciones - 20) {
		if (float64(vict_ataque+vict_defensa) / float64(simulaciones) * 100) >= float64(proceso_porct+1) { // porcentaje de calculo simulaciones
			proceso_porct = int(float64(vict_ataque+vict_defensa) / float64(simulaciones) * 100)
			//muestreo(simulacion) # llamo a muestreo
			fmt.Println("Simulando: %", proceso_porct)
			time.Sleep(3000 * time.Millisecond)
		}
	}

	//fmt.Scanln()
	porcent_vict := float32(float64(vict_ataque) / float64(simulaciones) * 100.0)
	porcent_derr := float32(float64(vict_defensa) / float64(simulaciones) * 100.0)
	fmt.Println("Tiempo consumido en el cálculo : ")
	tiempofinal := time.Now()
	fmt.Println("Victoria ataque: ", vict_ataque, " %", porcent_vict)
	fmt.Println("victoria defensa: ", vict_defensa, " %", porcent_derr)
	fmt.Println("el calculo tardó: ", tiempofinal.Sub(comienzo))

	// base de datos
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
	dbdata := db
	fmt.Println(dbdata.Driver())
	insertarRegistros, err := db.Prepare("INSERT INTO registro (fecha, simulaciones, fichas_ataque, fichas_defensa,	porct_victoria_ataque, porct_victoria_defensa) VALUES(?,?,?,?,?,?)")
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
