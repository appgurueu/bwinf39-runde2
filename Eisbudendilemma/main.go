package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/fogleman/gg"
)

// Liest Umfang und Adressen aus einer Datei
func einlesen(pfad string) (int, []int) {
	file, err := os.Open(pfad)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Umfang einlesen
	scanner.Scan()
	umfang, err := strconv.Atoi(strings.Split(scanner.Text(), " ")[0])
	if err != nil {
		panic(err)
	}
	// Adressen einlesen
	scanner.Scan()
	_addressen := strings.Split(scanner.Text(), " ")
	addressen := make([]int, len(_addressen))
	for index, addresse := range _addressen {
		var err error
		addressen[index], err = strconv.Atoi(addresse)
		if err != nil {
			panic(err)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return umfang, addressen
}

const radius float64 = 300
const rand float64 = 50
const punktradius float64 = 2
func speicherePNG(pfad string, umfang int, adressen, eisbuden []int) {
	zentrum := radius + rand
	dc := gg.NewContext(int(zentrum * 2), int(zentrum * 2))
	dc.DrawCircle(zentrum, zentrum, radius)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	// Bestimmt die Koordinaten eines Kreispunkts
	kreispunkt := func(adresse int, radius float64) (float64, float64) {
		winkel := float64(adresse) * 2 * math.Pi / float64(umfang)
		posX := math.Cos(winkel) * radius + zentrum
		posY := math.Sin(winkel) * radius + zentrum
		return posX, posY
	}
	// Adressen zeichnen
	for _, adresse := range adressen {
		posX, posY := kreispunkt(adresse, radius + punktradius)
		dc.DrawCircle(posX, posY, punktradius)
		dc.SetRGB(1, 0, 0)
		dc.Fill()
		posX, posY = kreispunkt(adresse, radius + rand/4)
		dc.DrawStringAnchored(strconv.Itoa(adresse), posX, posY, 0.5, 0.5)
	}
	// Eisbuden zeichnen
	for _, adresse := range eisbuden {
		posX, posY := kreispunkt(adresse, radius - punktradius)
		dc.DrawCircle(posX, posY, punktradius)
		dc.SetRGB(0, 1, 0)
		dc.Fill()
		posX, posY = kreispunkt(adresse, radius - rand/4)
		dc.DrawStringAnchored(strconv.Itoa(adresse), posX, posY, 0.5, 0.5)
	}
	// Speichern
	dc.SavePNG(pfad)
}

func main() {
	// Argumente
	if len(os.Args) < 2 || len(os.Args) > 4 {
		println("Verwendung: <pfad> [beste-loesung] [png-ausgabe-pfad]")
		return
	}
	ermittle_beste_loesung := true
	if len(os.Args) >= 3 {
		if os.Args[2] != "j" && os.Args[2] != "n" {
			println("Verwendung: <pfad> [beste-loesung] [png-ausgabe-pfad]")
			return
		}
		ermittle_beste_loesung = os.Args[2] == "j"
	}
	// Aufgabenstellung einlesen
	umfang, adressen := einlesen(os.Args[1])
	if len(adressen) <= 3 {
		// Triviale Lösung
		fmt.Println("Lösung:", adressen)
		return
	}
	// Es wird nicht garantiert, dass die Adressen sortiert sind.
	sort.Ints(adressen)
	// Drehe so, dass erste Adresse 0 ist
	drehung := adressen[0]
	for i := range adressen {
		adressen[i] -= drehung
	}

	// Abstand zweier Adressen
	abstand := func(a, b int) int {
		abs := a - b
		if abs < 0 {
			// Betrag
			abs = -abs
		}
		andersrum := umfang - abs
		if andersrum < abs {
			// Weg ist kürzer in entgegengesetzter Richtung
			return andersrum
		}
		return abs
	}

	// Bestimmt den Abstand einer Adresse zur nächsten Eisbude
	minAbstand := func(adresse int, eisbuden []int) int {
		min_abs := math.MaxInt32
		for _, eisbude := range eisbuden {
			abs := abstand(adresse, eisbude)
			if abs < min_abs {
				min_abs = abs
			}
		}
		return min_abs
	}

	// Summe der Abstände
	summeMinAbstaende := func(eisbuden []int) int {
		summe := 0
		for _, adresse := range adressen {
			summe += minAbstand(adresse, eisbuden)
		}
		return summe
	}

	// Bestimmt, ob andere_eisbuden gegen eisbuden in einer Abstimmung gewinnen würde
	istBesser := func(andere_eisbuden, eisbuden []int) bool {
		stimmen := 0
		for _, adresse := range adressen {
			if minAbstand(adresse, andere_eisbuden) < minAbstand(adresse, eisbuden) {
				// Verkürzung, Stimme dafür
				stimmen++
			} else {
				// Stimme dagegen
				stimmen--
			}
		}
		// Mehr Ja- als Nein-Stimmen
		return stimmen > 0
	}

	// So viele Adressen müssen mindestens zwischen zwei Eisbuden liegen
	min_adressen_zwischen_eisbuden := len(adressen) / 3 - 3
	if min_adressen_zwischen_eisbuden < 0 {
		min_adressen_zwischen_eisbuden = 0
	}
	getNaechsteAdresse := func(adresse int) int {
		// Binäre Suche: Bestimmt, hinter der wievielten Häusern die gegebene Adresse liegt
		i := sort.SearchInts(adressen, adresse)
		if i >= len(adressen) {
			// Adresse liegt hinter dem letzten Haus
			return umfang
		}
		i += min_adressen_zwischen_eisbuden
		if i >= len(adressen) {
			return umfang
		}
		if adresse <= adressen[i] {
			return adresse + 1
		}
		return adressen[i]
	}

	// Ermittelt, ob eine Eisbudenverteilung stabil ist
	// Um dies zu beschleunigen, werden zahlreiche "early-returns" genutzt:
	// Es wird versucht, möglichst schnell Standorte zu finden,
	// gegen die die gegebenen Standorte verlieren würden
	istLoesung := func(eisbuden []int) bool {
		// Häufig bessere Standorte: Probiere simple Verschiebungen
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				for k := -1; k < 2; k++ {
					if istBesser([]int{
						(eisbuden[0]+i) % umfang,
						(eisbuden[1]+j) % umfang,
						(eisbuden[2]+k) % umfang,
					}, eisbuden) {
						return false
					}
				}
			}
		}
		// Eisbudenstandorte bei Adressen probieren
		// Häufig besser
		for i := 0; i < len(adressen); i++ {
			for j := i + 1; j < len(adressen); j++ {
				for k := j + 1; k < len(adressen); k++ {
					andere_eisbuden := []int{adressen[i], adressen[j], adressen[k]}
					if istBesser(andere_eisbuden, eisbuden) {
						return false
					}
				}
			}
		}
		// Eisbudenstandorte mit Abstand untereinander probieren
		for i := 0; i < umfang; i++ {
			for j := getNaechsteAdresse(i); j < umfang; j++ {
				for k := getNaechsteAdresse(j); k < umfang; k++ {
					andere_eisbuden := []int{i, j, k}
					if istBesser(andere_eisbuden, eisbuden) {
						return false
					}
				}
			}
		}
		// Keine Standorte gefunden, gegen die die gegebenen Standorte verlieren würden
		return true
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	min_summe := math.MaxInt32
	var beste_loesung []int
	ausgabe := func() {
		// Zurück drehen
		for i := range adressen {
			adressen[i] += drehung
		}
		for i := range beste_loesung {
			beste_loesung[i] += drehung
		}
		if ermittle_beste_loesung {
			fmt.Println("Lösung: Adressen", beste_loesung, "mit Summe der Abstände", min_summe)
		} else {
			fmt.Println("Lösung: Adressen", beste_loesung)
		}
		if len(os.Args) > 3 {
			speicherePNG(os.Args[3], umfang, adressen, beste_loesung)
		}
	}
	// Probiert Eisbudenstandorte aus, bei denen die erste Bude Adresse i hat
	probiere := func(i int) {
		defer wg.Done()
		// Probiere so aus, dass i < j < k gilt
		// und zwischen i und j sowie j und k ausreichend Häuser liegen
		for j := getNaechsteAdresse(i); j < umfang; j++ {
			for k := getNaechsteAdresse(j); k < umfang; k++ {
				eisbuden := []int{i, j, k}
				summe := -1
				if ermittle_beste_loesung {
					summe = summeMinAbstaende(eisbuden)
				}
				if summe < min_summe && istLoesung(eisbuden) {
					// Mutex "sperren":
					// Wird währenddessen eine weitere Lösung gefunden,
					// soll diese warten - die Lösungen
					// sollen hintereinander durchgegangen werden
					mutex.Lock()
					beste_loesung = eisbuden
					min_summe = summe
					if !ermittle_beste_loesung {
						// Erste gefundene Lösung ausgeben
						ausgabe()
						// Fertig
						os.Exit(0)
					}
					// Mutex entsperren
					mutex.Unlock()
				}
			}
		}
	}
	for i := 0; i < umfang; i++ {
		wg.Add(1)
		go probiere(i)
	}
	// Warte auf das Ausprobieren aller infragekommender Standorte
	wg.Wait()
	if min_summe != math.MaxInt32 {
		// Lösung gefunden: min_summe hat nicht mehr den Maximalwert
		ausgabe()
	} else {
		fmt.Println("Keine Lösung möglich")
	}
}
