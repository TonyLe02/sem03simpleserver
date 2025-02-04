package main

import (
	"io"
	"log"
	"net"
	"sync"
	"github.com/TonyLe02/is105sem03/mycrypt"
	"github.com/TonyLe02/funtemps/conv"
	"fmt"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					switch msg := string(dekryptertMelding); msg {
						  				        
					case "ping":
						_, err = c.Write([]byte("pong"))

					case "Kjevik;SN39040;18.03.2022 01:50;6":

						var celsius float64 = 6
						fahrenheit := conv.CelsiusToFarhenheit(celsius)
						response := fmt.Sprintf("Kjevik;SN39040;18.03.2022 01:50;%.2f", fahrenheit)
						_, err = c.Write([]byte(response))

					default:
						_, err = c.Write([]byte(msg))
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
