package main

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/jonaslefdal/is105sem03/mycrypt"
	"github.com/jonaslefdal/minyr/yr"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:41391")
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
					dekryptertMelding := mycrypt.Krypter(([]rune(string(buf[:n]))), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					msg := string(dekryptertMelding)
					switch msg {
					case "ping":
						kryptertMelding := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(string(kryptertMelding)))
					default:
						if strings.HasPrefix(msg, "Kjevik") {
							dekryptertMeldingfahr, err := yr.CelsiusToFahrenheitLine(msg)
							if err != nil {
								log.Println(err)
							}
							log.Println("Dekrypter melding i fahr: ", string(dekryptertMeldingfahr))
							kryptertMelding := mycrypt.Krypter([]rune(dekryptertMeldingfahr), mycrypt.ALF_SEM03, 4)
							_, err = c.Write([]byte(string(kryptertMelding)))
						} else {
							kryptertMelding := mycrypt.Krypter([]rune(dekryptertMelding), mycrypt.ALF_SEM03, 4)
							log.Println("Kryptert melding: ", string(kryptertMelding))
							_, err = c.Write([]byte(string(kryptertMelding)))
						}
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
