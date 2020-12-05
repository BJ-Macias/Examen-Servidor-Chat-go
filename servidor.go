package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
)

type Mensaje struct {
	Nickname string
	Mensaje  string
	Bytes    []byte
}

type Server struct {
	usuarios map[string]string
	Mensajes []string
	archivos map[string]string
}

func (this *Server) RegistraIngresos(datos Mensaje, reply *string) error {

	if this.usuarios == nil {
		this.usuarios = make(map[string]string)
	}

	usu, existe := this.usuarios[datos.Nickname]
	if existe {
		return errors.New("ya existe usuario " + usu)
	}
	this.usuarios[datos.Nickname] = "Activo"

	fmt.Println(" Se conecto: " + datos.Nickname)
	*reply = " Bienvenido " + datos.Nickname

	this.Mensajes = append(this.Mensajes, " Se conecto: "+datos.Nickname+"\n")
	//fmt.Println(this.Mensajes)

	return nil
}

func (this *Server) EnviarMensaje(datos Mensaje, reply *string) error {

	if datos.Mensaje == "" {
		return errors.New(" Mensaje vacio")
	}
	fmt.Println(datos.Nickname + ": " + datos.Mensaje)
	this.Mensajes = append(this.Mensajes, " "+datos.Nickname+": "+datos.Mensaje+"\n")

	*reply = " Mensaje enviado"
	return nil
}

func (this *Server) EnviarArchivo(datos Mensaje, reply *string) error {

	if this.archivos == nil {
		this.archivos = make(map[string]string)
	}

	arch, existe := this.archivos[datos.Mensaje]
	if existe {
		return errors.New("ya existe archivo " + arch)
	}
	this.archivos[datos.Mensaje] = string(datos.Bytes)

	err := ioutil.WriteFile("D:/Documentos/Sistemas Distribuidos/ExamenParcial1/ArchivosEnvidos/" + datos.Mensaje, datos.Bytes, 0640)
    if err != nil {
        log.Fatal(err)
    }

	*reply = " Archivo: " + datos.Mensaje + " enviado"

	this.Mensajes = append(this.Mensajes, " "+datos.Nickname+" envio archivo: "+datos.Mensaje+"\n")
	fmt.Println(" " + datos.Nickname + " envio archivo: " + "'" + datos.Mensaje + "'")

	return nil
}

func (this *Server) MostrarMensajes(datos Mensaje, reply *string) error {
	if len(this.Mensajes) < 0 {
		return errors.New("No hay mensajes que mostrar")
	}

	var chat string
	for _, v := range this.Mensajes {
		chat += v
	}

	*reply = "\n" + chat
	return nil
}

func (this *Server) QuitarUsuario(datos Mensaje, reply *string) error {
	if this.usuarios == nil {
		return errors.New("No hay usuarios activos")
	}

	*reply = " Asta la proxima " + datos.Nickname
	fmt.Println(" " + datos.Nickname + " se desconecto")
	this.Mensajes = append(this.Mensajes, " "+datos.Nickname+" se desconecto"+"\n")
	delete(this.usuarios, datos.Nickname)

	return nil

}

func (this *Server) Respaldar(datos Mensaje, reply *string) error {
	if len(this.Mensajes) < 0 {
		return errors.New("No hay mensajes para respaldar")
	}
	var respaldo string
	var input string
	fmt.Print("Ingrese el nombre del respaldo: ")
	fmt.Scanln(&input)
	for _, v := range this.Mensajes {
		respaldo += v
	}
	b := []byte(respaldo)
	err := ioutil.WriteFile(input, b, 0640)
	if err != nil {
		log.Fatal(err)
	}
	*reply = " Achivo: " + input + " creado"
	return nil
}

func server() {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	go server()
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	var msg Mensaje
	var result string

	var op int64
	for {
		fmt.Println("1.- Mostrar mensajes")
		fmt.Println("2.- Respaldar chat")
		fmt.Println("3.- Salir")
		fmt.Scanln(&op)

		switch op {
		case 1:
			err = c.Call("Server.MostrarMensajes", msg, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
			break
		case 2:
			err = c.Call("Server.Respaldar", msg, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
			break
		case 3:
			return
		}
		fmt.Println()
	}
}