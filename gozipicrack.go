/*
BUG:
Parar cuando encuentre el pass.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/yeka/zip"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage gozipcrack.go -f file.zip -d Dictionary")
		os.Exit(1)
	}
	// Choose zip file to crack.
	zipfile := flag.String("f", " ", "Zip file to crack")
	// Choose dictionary.
	dictfile := flag.String("d", "/usr/share/wordlists/rockyou.txt", "Dictionary")
	flag.Parse()
	z := *zipfile
	df := *dictfile
	r, err := zip.OpenReader(z)
	if err != nil {
		fmt.Println("Zip file doesnt exist.")
		return
	}
	defer r.Close() // <-- Con defer le decimos que si termina main, que cierre el archivo zip
	// Replace with your dictionary.
	lines, err := readLines(df)
	if err != nil {
		fmt.Println("Dictionary doesnt exist.")
		return
	}
	wg := &sync.WaitGroup{}
	c := 0
	for _, line := range lines {
		if c != 1 {
			wg.Add(1) // <-- Sólo añades el hilo si c != 1, si no, creas hilos que nunca acaban
			func(line string) {
				defer wg.Done()
				for _, f := range r.File {
					if f.IsEncrypted() {
						f.SetPassword(line)
					}
					p, err := f.Open()
					if err != nil {
						log.Fatal(err)
					}
					buf, err := ioutil.ReadAll(p)
					if err == nil { // <-- El break realmente no hace falta, dado que no tiene ningún código dentro
						fmt.Printf("Password found: %v", line)
						fmt.Printf("\nSize of %v: %v byte(s)\n", f.Name, len(buf))
						c = 1
					}
					p.Close()
				}
			}(line)
		}
	}
	wg.Wait()
}
