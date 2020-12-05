package main

import (
	"io/ioutil"
	"os"
	"bufio"
	"fmt"
	"net/rpc"
)

type Mensaje struct {
	Nickname string
	Mensaje  string
	Bytes    []byte
}

func main() {
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	var msg Mensaje
	var result string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Escribe tu nickname: ")
	scanner.Scan()
	msg.Nickname = scanner.Text()
	err = c.Call("Server.RegistraIngresos", msg, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	var op int64
	for {
		fmt.Println("1.- Enviar mensaje")
		fmt.Println("2.- Enviar archivo")
		fmt.Println("3.- Mostrar Mensajes")
		fmt.Println("4.- Salir")
		fmt.Scanln(&op)

		switch op {
		case 1:
			fmt.Print("Escriba mensaje: ")
			scanner.Scan()
			msg.Mensaje = scanner.Text()

			err = c.Call("Server.EnviarMensaje", msg, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
			break
		case 2:
			var sn string
			fmt.Print("¿Quiere espesificar la ruta de archivo? s/n: ")
			fmt.Scanln(&sn)
			switch sn {
				case "s":
					fmt.Print("Escriba nombre de archivo: ")
					scanner.Scan()
					msg.Mensaje = scanner.Text()
					fmt.Print("Ingrese path del archivo: ")
					scanner.Scan()
					path := scanner.Text()
					msg.Bytes, err = ioutil.ReadFile(path + msg.Mensaje)
					if err != nil { 
						fmt.Print(err) 
					}
					
					err = c.Call("Server.EnviarArchivo", msg, &result)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(result)
					}
				break
				case "n":
					fmt.Print("Escriba nombre de archivo: ")
					scanner.Scan()
					msg.Mensaje = scanner.Text()
					msg.Bytes, err = ioutil.ReadFile(msg.Mensaje)
					if err != nil { 
						fmt.Print(err) 
					}
					
					err = c.Call("Server.EnviarArchivo", msg, &result)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(result)
					}
				break
			}
			break
		case 3:
			err = c.Call("Server.MostrarMensajes", msg, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
			break
		case 4:
			err = c.Call("Server.QuitarUsuario", msg, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
			return
		}
		fmt.Println()
	}
}