package main

import (
	"io"
	"log"
	"net"
	"sync"
	"strings"
	"github.com/erny03/is105sem03/mycrypt"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8110")
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

					msg := string(dekryptertMelding)
					switch {
  				        	case strings.HasPrefix(msg, "ping"):
							kryptertPong := mycrypt.Krypter([]rune(string("pong")), mycrypt.ALF_SEM03, 4)
							_, err = c.Write([]byte(string(kryptertPong)))

						case strings.HasPrefix(msg, "Kjevik"):
							kryptertKjevik := mycrypt.Krypter([]rune("Kjevik;SN39040;18.03.2022 01:50;42.8"), mycrypt.ALF_SEM03, 4)
							_, err = c.Write([]byte(string(kryptertKjevik)))

						default:
						_, err = c.Write(buf[:n])
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
