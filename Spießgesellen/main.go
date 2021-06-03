package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// BitSet für bis zu 64 Einträge
type BitSet uint64
const SET_CAPACITY uint8 = 64

// Enthält Element
func (set BitSet) Has(index uint8) bool {
	return (set & (1 << index)) > 0
}

// Füge Element hinzu
func (set BitSet) Add(index uint8) BitSet {
	return set | (1 << index)
}

// Füge Schale hinzu: 1 abziehen
func (set BitSet) AddSchale(schale uint8) BitSet {
	return set.Add(schale - 1)
}

// Schalen-Liste für Ausgabe erstellen: Zu jedem Element wieder 1 hinzuaddieren
func (set BitSet) ToSchalen() []uint8 {
	schalen := make([]uint8, set.Count())
	cursor := 0
	for i := uint8(0); i < frucht_nummer; i++ {
		if set.Has(i) {
			schalen[cursor] = i + 1
			cursor++
		}
	}
	return schalen
}

// Symmetrische Differenz (xor)
func (set BitSet) Difference(other BitSet) BitSet {
	return set ^ other
}

// Differenz / Abzug (and not)
func (set BitSet) Except(other BitSet) BitSet {
	return set &^ other
}

// Schnitt (and)
func (set BitSet) Intersect(other BitSet) BitSet {
	return set & other
}

// Vereinigung (or)
func (set BitSet) Union(other BitSet) BitSet {
	return set | other
}

// Byte-Set counts lookup array
var counts [256]uint8

func (set BitSet) Count() uint8 {
	count := uint8(0)
	for shift := uint8(0); shift < frucht_nummer; shift += 8 {
		count += counts[(set >> shift) & 0xFF]
	}
	return count
}

// Ist set eine Obermenge von other?
func (set BitSet) Super(other BitSet) bool {
	return set.Union(other) == set
}

// Ist set eine Untermenge von other?
func (set BitSet) Sub(other BitSet) bool {
	return other.Super(set)
}

// Frucht-Operationen
// Frucht-Namen: [index] = name
var frucht_name []string
// Frucht-Nummer: Wird zwischenzeitlich (während dem Einlesen) als Cursor verwendet
// Ist danach = die Anzahl der Früchte
var frucht_nummer uint8 = 0
// Frucht-Nummern: [name] = index
var frucht_nummern map[string]uint8 = map[string]uint8{}

// Frucht hinzufügen
func (set BitSet) AddFrucht(frucht string) BitSet {
	index, found := frucht_nummern[frucht]
	// Hat die Frucht noch keine Nummer?
	if !found {
		// Nummer zuweisen
		frucht_name[frucht_nummer] = frucht
		frucht_nummern[frucht] = frucht_nummer
		index = frucht_nummer
		// Nächste Nummer für nächste Frucht
		frucht_nummer++
	}
	return set.Add(index)
}

// Frucht-Liste -> Frucht-Menge
func FromFruechte(fruechte []string) BitSet {
	var set BitSet
	for _, frucht := range fruechte {
		set = set.AddFrucht(frucht)
	}
	return set
}

// Frucht-Menge -> Frucht-Liste für Ausgabe
func (set BitSet) ToFruechte() []string {
	res := make([]string, set.Count())
	c := 0
	for i := uint8(0); i < frucht_nummer; i++ {
		if set.Has(i) {
			res[c] = frucht_name[i]
			c++
		}
	}
	return res
}

func main() {
	// Argumente lesen
	if len(os.Args) < 2 || len(os.Args) > 3 {
		println("Verwendung: <pfad> [praeferenzen]")
		return
	}
	praeferenzen := "*"
	if len(os.Args) == 3 {
		praeferenzen = os.Args[2]
		if praeferenzen != "*" && praeferenzen != "+-" && praeferenzen != "+" && praeferenzen != "-" {
			println("Verwendung: <pfad> [praeferenzen]")
			return
		}
	}
	// Datei öffnen
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	anzahl_fruechte, _ := strconv.Atoi(scanner.Text())
	if anzahl_fruechte > int(SET_CAPACITY) {
		panic("Maximal " + strconv.Itoa(int(SET_CAPACITY)) + " Früchte erlaubt")
	}
	frucht_name = make([]string, anzahl_fruechte)
	scanner.Scan()
	// Zielmenge der gewünschten Früchte erstellen
	var ziel_fruechte BitSet
	for _, frucht := range strings.Split(scanner.Text(), " ") {
		if frucht != "" {
			ziel_fruechte = ziel_fruechte.AddFrucht(frucht)
		}
	}
	scanner.Scan()
	anzahl_beobachtungen, _ := strconv.Atoi(scanner.Text())
	gleichungen := map[BitSet]BitSet{}
	for i := uint8(0); i < uint8(anzahl_beobachtungen); i++ {
		// Schalen-Menge einlesen
		scanner.Scan()
		var schalen BitSet
		for _, schale := range strings.Split(scanner.Text(), " ") {
			if schale == "" {
				// HACK wegen trailing spaces
				continue
			}
			_schale, err := strconv.Atoi(schale)
			if err != nil {
				panic(err)
			}
			schalen = schalen.AddSchale(uint8(_schale))
		}
		// Früchte-Menge einlesen
		scanner.Scan()
		var fruechte BitSet
		for _, frucht := range strings.Split(scanner.Text(), " ") {
			if frucht == "" {
				// HACK wegen trailing spaces
				continue
			}
			fruechte = fruechte.AddFrucht(frucht)
		}
		// Mengen-Gleichung in Map speichern
		gleichungen[fruechte] = schalen
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Initialisiere lookup-Array für effizientes Zählen der gesetzten Bits
	// Erst danach darf BitSet.Count() verwendet werden!
	for i := 0; i < 256; i++ {
		for j := 0; j < 8; j++ {
			counts[i] += uint8(1 & (i >> j))
		}
	}

	// Auffüllen: Alle unbenannten, unbekannten Früchte erhalten einen Namen zugewiesen
	// Diese kommen ausschließlich in der "Insgesamt"-Gleichung vor
	for i := 0; frucht_nummer < uint8(anzahl_fruechte); frucht_nummer++ {
		i++
		name := "Unbekannt_" + strconv.Itoa(i)
		frucht_name[frucht_nummer] = name
	}

	// "Insgesamt"-Gleichung: Menge aller Früchte = Menge aller Schalen
	// Setze unterste anzahl_fruechte Bits auf 1
	insgesamt := ^(^BitSet(0) << BitSet(anzahl_fruechte))
	gleichungen[insgesamt] = insgesamt

	// Teilmengen der Zielmenge
	var teil_fruechte, teil_schalen BitSet
	for {
		neue_gleichungen_erzeugt := false
		neue_gleichungen := map[BitSet]BitSet{}
		neueGleichung := func(fruechte, schalen BitSet) {
			// Neue Gleichung: Schon bestimmte Teilmengen abziehen
			neue_fruechte := fruechte.Except(teil_fruechte)
			// Ermitteln, ob eine neue Gleichung erzeugt wurde
			_, alt := gleichungen[neue_fruechte]
			neue_gleichungen_erzeugt = neue_gleichungen_erzeugt || !alt
			// Gleichung setzen
			neue_gleichungen[neue_fruechte] = schalen.Except(teil_schalen)
		}
		for fruechte, schalen := range gleichungen {
			if (fruechte.Sub(ziel_fruechte)) {
				// Untermenge der Zielmenge: Zu insgesamt-Teilmengen hinzufügen
				teil_fruechte = teil_fruechte.Union(fruechte)
				teil_schalen = teil_schalen.Union(schalen)
				continue
			}
		}
		for fruechte, schalen := range gleichungen {
			for fruechte_2, schalen_2 := range gleichungen {
				// Differenz zweier Mengen
				neueGleichung(fruechte.Except(fruechte_2), schalen.Except(schalen_2))
				if fruechte_2 > fruechte {
					// Schnitt zweier Mengen; Kommutativ, daher die if-Bedingung
					neueGleichung(fruechte.Intersect(fruechte_2), schalen.Intersect(schalen_2))
				}
			}
		}

		gleichungen = neue_gleichungen
		if !neue_gleichungen_erzeugt {
			// Aufhören, wenn keine neue Gleichungen mehr erzeugt werden können
			break
		}
	}

	if teil_fruechte == ziel_fruechte {
		fmt.Println("Die Früchte befinden sich in den Schalen:", teil_schalen.ToSchalen())
		return
	}

	neue_gleichungen := map[BitSet]BitSet{}
	nuetzlicheGleichungen: for fruechte, schalen := range gleichungen {
		intersection := fruechte.Intersect(ziel_fruechte).Except(teil_fruechte)
		if intersection == 0 {
			// Kein Schnitt mit Zielmenge oder alle Schalen schon bekannt
			// Gleichung nutzlos
			continue
		}
		for fruechte_2 := range gleichungen {
			intersection_2 := fruechte_2.Intersect(ziel_fruechte).Except(teil_fruechte)
			if intersection == intersection_2 && fruechte_2.Count() < fruechte.Count() {
				// Gleicher Schnitt, aber Früchte 2 besitzt weniger unerwünschte Früchte
				continue nuetzlicheGleichungen
			}
		}
		neue_gleichungen[fruechte] = schalen
	}

	neue_gleichungen[teil_fruechte] = teil_schalen

	// 2. Stufe
	gleichungen = neue_gleichungen
	for {
		neue_gleichungen_erzeugt := false
		neue_gleichungen := map[BitSet]BitSet{}
		for fruechte, schalen := range gleichungen {
			for fruechte_2, schalen_2 := range gleichungen {
				fruechte_vereinigung := fruechte.Union(fruechte_2)
				abweichung := ziel_fruechte.Except(fruechte_vereinigung).Count()
				// Schnitt ist kommutativ.
				if fruechte_2 < fruechte {
					continue
				}
				// Außerdem soll sich die vereinigte Menge nicht von der Zielmenge "entfernen".
				if abweichung > ziel_fruechte.Except(fruechte).Count() || abweichung > ziel_fruechte.Except(fruechte_2).Count() {
					continue
				}
				schalen_vereinigung := schalen.Union(schalen_2)
				neue_gleichungen[fruechte_vereinigung] = schalen_vereinigung
			}
		}
		for fruechte, schalen := range neue_gleichungen {
			_, alt := gleichungen[fruechte]
			if !alt {
				neue_gleichungen_erzeugt = true
				gleichungen[fruechte] = schalen
			}
		}
		// Aufhören, wenn keine neuen Gleichungen mehr erzeugt wurden
		if !neue_gleichungen_erzeugt {
			break
		}
	}

	minAbweichungInsgesamt := func() (uint8, map[BitSet]BitSet) {
		min_abweichung := uint8(255)
		beste_gleichungen := map[BitSet]BitSet{}
		// Iteriert alle Gleichungen und bestimmt die mit der geringsten Abweichung
		for fruechte, schalen := range gleichungen {
			// Abweichung: Anzahl Elemente der symmetrischen Differenz
			abweichung := fruechte.Difference(ziel_fruechte).Count()
			if abweichung < min_abweichung {
				min_abweichung = abweichung
				beste_gleichungen = map[BitSet]BitSet{}
			}
			if abweichung <= min_abweichung {
				beste_gleichungen[fruechte] = schalen
			}
		}
		return min_abweichung, beste_gleichungen
	}
	// Minimiert Anzahl fehlender oder unerwünschter Früchte je nach Argument
	minAbweichung := func(minimiere_fehlend bool) (uint8, uint8, map[BitSet]BitSet) {
		min_1 := uint8(255)
		min_2 := uint8(255)
		beste_gleichungen := map[BitSet]BitSet{}
		// Gleichungen durchgehen
		for fruechte, schalen := range gleichungen {
			abweichung_1 := fruechte.Except(ziel_fruechte).Count()
			abweichung_2 := ziel_fruechte.Except(fruechte).Count()
			if minimiere_fehlend {
				// Fehlende minimieren: Prioritäten tauschen!
				abweichung_2, abweichung_1 = abweichung_1, abweichung_2
			}
			if abweichung_1 > min_2 {
				continue
			}
			if abweichung_1 < min_2 {
				min_2 = abweichung_1
				min_1 = abweichung_2
				beste_gleichungen = map[BitSet]BitSet{}
			} else if abweichung_2 > min_1 {
				continue
			} else if abweichung_2 < min_1 {
				min_1 = abweichung_2
				beste_gleichungen = map[BitSet]BitSet{}
			}
			beste_gleichungen[fruechte] = schalen
		}
		return min_1, min_2, beste_gleichungen
	}
	// Gibt eine alternative Gleichung aus
	ausgabeAlternative := func(gleichung map[BitSet]BitSet) {
		for fruechte, schalen := range gleichung {
			erwuenscht := fruechte.Intersect(ziel_fruechte)
			unerwuenscht := fruechte.Except(erwuenscht)
			fehlend := ziel_fruechte.Except(erwuenscht)
			teile := []string{}
			// Formattiert einen Teil der Frucht-Menge
			teil := func(prefix string, set BitSet) {
				fruechte := set.ToFruechte()
				if len(fruechte) > 0 {
					teile = append(teile, prefix + strings.Join(fruechte, ","))
				}
			}
			// Erwünschter Teil
			teil("=", erwuenscht)
			// Unerwünschter Teil
			teil("+", unerwuenscht)
			// Fehlender Teil
			teil("-", fehlend)
			fmt.Println(strings.Join(teile, "; "), " : ", schalen.ToSchalen())
		}
	}
	ausgabeAlternativen := func(gesucht string) {
		var min_insgesamt, min_fehlend, min_unerwuenscht uint8
		var gleichungen map[BitSet]BitSet
		// Ausgabe je nach gesuchten Alternativen
		if gesucht == "+-" {
			// Minimale insgesamte Abweichung
			min_insgesamt, gleichungen = minAbweichungInsgesamt()
			fmt.Println("Minimale insgesamte Abweichung", min_insgesamt, "vom Gewünschten (+-):")
		} else if gesucht == "+" {
			// Minimal viele unerwünschte Früchte
			min_fehlend, min_unerwuenscht, gleichungen = minAbweichung(false)
			fmt.Println("Nur", min_fehlend, "fehlende und", min_unerwuenscht, "unerwünschte Früchte (+):")
		} else if gesucht == "-" {
			// Minimal viele fehlende Früchte
			min_unerwuenscht, min_fehlend, gleichungen = minAbweichung(true)
			fmt.Println("Nur", min_unerwuenscht, "unerwünschte und", min_fehlend, "fehlende Früchte (-):")
		}
		ausgabeAlternative(gleichungen)
	}
	if praeferenzen == "*" {
		// Alle nach einem der drei Kriterien optimalen Alternativen
		fmt.Println("Keine exakte Bestimmung der Schalen möglich. Alternativen:")
		fmt.Println()
		ausgabeAlternativen("+-")
		fmt.Println()
		ausgabeAlternativen("+")
		fmt.Println()
		ausgabeAlternativen("-")
	} else {
		ausgabeAlternativen(praeferenzen)
	}
}
