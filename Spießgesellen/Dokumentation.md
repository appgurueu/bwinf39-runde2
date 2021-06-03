# Lösung "Spießgesellen"

Lars Müller, Teilnahme-ID 58886

## Problem

Gegeben:

* Apfel, Banane und Brombeere: Schüsseln 1, 4 und 5
* Banane, Pflaume und Weintraube: Schüsseln 3, 5 und 6
* Apfel, Brombeere und Erdbeere: Schüsseln 1, 2 und 4
* Erdbeere und Pflaume: Schüsseln 2 und 6

Gesucht:

* Weintraube, Brombeere und Apfel: Welche Schüsseln?

## Formalisierung

Mengen-Gleichungen: Jede Frucht ist eine Schüssel zuzuordnen; eine Menge von Früchten gleicht einer Menge von Schüsseln:

Seien die Variablen für die Schüsseln nach den Anfangsbuchstaben der Früchte [A]pfel, [B]anane, [P]flaume, [W]eintraube, [E]rdbeere benannt, mit Ausnahme der Variablen $R$ für die Nummer der die Brombeeren enthaltenden Schüssel.

Gegeben:

* $\{A, B, R\} = \{1, 4, 5\}$ (1)
* $\{B, P, W\} = \{3, 5, 6\}$ (2)
* $\{A, R, E\} = \{1, 2, 4\}$ (3)
* $\{E, P\} = \{2, 6\}$ (4)

Gesucht: Schlüsseln $x, y, z$, für die gilt $\{W, R, A\} = \{x, y, z\}$.

## Lösung

Die Mengen-Gleichungen können über Vereinigung, Schnitt und Differenz so kombiniert werden, dass wir die gewünschte Gleichung erhalten:

1. (3) \- (4): $\{A, R, E\} \setminus \{E, P\} = \{1, 2, 4\} \setminus \{2, 6\} \Leftrightarrow \{A, R\} = \{1, 4\}$ (5)
2. (2) \- (1): $\{B, P, W\} \setminus \{A, B, R\} = \{3, 5, 6\} \setminus \{1, 4, 5\} \Leftrightarrow \{P, W\} = \{3, 6\}$ (6)
3. (6) \- (4): $\{P, W\} \setminus \{E, P\} = \{3, 6\} \setminus \{2, 6\} \Leftrightarrow \{W\} = \{3\}$ (7)
4. (5) \+ (7): $\{A, R\} \cup \{W\} = \{1, 4\} \cup \{3\} \Leftrightarrow \{A, R, W\} = \{1, 3, 4\} \Leftrightarrow \{W, R, A\} = \{1, 3, 4\}$

Donald muss sich also aus den Schalen 1, 3 und 4 bedienen, um Weintraube $W$, Brombeere $R$ und Apfel $A$ zu erhalten.

## Generalisierung

Jedes solche Problem lässt sich analog als Mengen-Gleichungssystem schreiben. Gesucht ist dann eine ableitbare Mengen-Gleichung, wo auf der linken Frucht-Seite die von Donald gewünschten Früchte und auf der rechten Schalen-Seite die diese enthaltenden Schalen stehen.

## Alternativen

Manchmal lässt sich keine Menge von Schalen finden, die exakt die gewünschten Früchte enthält. Es lassen sich allerdings Mengen-Gleichungen finden, bei denen auf der Frucht-Seite manche Früchte fehlen oder andere unerwünscht sind. Drei mögliche Kriterien für die Auswahl einer alternativen Frucht-Menge sind:

1. Insgesamt-Abweichung: Das Fehlen einer Frucht ist genau gleich problematisch wie das Vorhandensein unerwünschter Früchte
2. Minimiere die Anzahl der *fehlenden* Früchte bei den Alternativen. Minimiere dann die Anzahl der *unerwünschten* Früchte ("Lieblingsfrüchte").
3. Minimiere die Anzahl der *unerwünschten* Früchte, dann die der *fehlenden* ("wählerisch").

# Verfahren

Eine Art "Mengen-Gleichungs-Löser" kombiniert über die genannten Mengenoperationen Gleichungen, und stellt weiter sicher, dass das Gleichungssystem währenddessen duplikatfrei bleibt. Dies macht man so lange, bis keine neuen Gleichungen mehr erzeugt werden können.

Nicht nur gemachte Beobachtungen stellen Gleichungen dar. Im gegebenen einfachen Beispiel reichen die gemachten Beobachtungen. In `spiesse1.txt` ist dies allerdings beispielsweise nicht der Fall: Die Grapefruit kommt in Donalds Wünschen vor, wurde von ihm aber noch nicht beobachtet. Dennoch kann er über ein Ausschlussverfahren bestimmen, in welcher Schale die Grapefruit sich befindet. Wir stellen fest, dass wir für solche Fälle die Gleichung "Menge aller Früchte = Menge aller Schalen" unserem Mengen-Gleichungssystem hinzufügen müssen.

# Stufen

Den Kombinationsprozess teilen wir in zwei Stufen auf. Beide Stufen sind im Grunde eine Breadth-First Search über die jeweiligen Mengenoperationen erzeugbaren Mengen - wobei allerdings ein großer, für die Lösung irrelevanter Teil weggelassen wird. Für die Bestimmung einer absolute oberen Grenze für die Anzahl der möglichen (Früchte-)Mengen - und somit auch Gleichungen - lässt sich die Potenzmenge der alle Früchte enthaltenden Mengen heranziehen: Alle Mengengleichungen können nur Teilmengen dieser Menge verwenden. Die Kardinalität der Potenzmenge ist $2^n$ - exponentiell. Tatsächlich liegen die bestimmten Mengengleichungen allerdings weit unter diesem Wert - eine Vielzahl der Elemente der Potenzmenge ist praktisch aus den gegebenen Beobachtungen einfach nicht ableitbar; andere ableitbare Gleichungen werden wie erwähnt weggelassen, oder wie die Teilmengen der Zielmenge in einer Gleichung komprimiert. In der ersten Stufe werden nur die reduzierenden Mengenoperationen Differenz und Schnitt verwendet, in der zweiten nur die additive Vereinigung.

Es gilt: Jeder Term bestehend aus Mengen und den drei Mengenoperationen lässt sich als Vereinigung von Termen aus Schnitten und Differenzen darstellen. Dies folgt aus den Mengen-Gesetzen. Exemplarisch:

$(A \cup B) \cap C = A \cup (B \cap C)$

$A \cap (B \cup C) = A \cup (B \cap C)$

$(A \cup B) \setminus C = (A \setminus B) \cup (A \setminus C)$

$A \setminus (B \cup C) = (A \setminus B) \setminus C$

In der ersten "destruktiven" Stufe behalten wir stets nur die neu entstandenen Gleichungen. Werden für alle Mengen $A$ und $B$ mit $A \neq B$ die Mengen $C = A \ B$, $D = B \ A$ und $E = A \cap B$ gebildet, geht keine Information verloren: $C \cup E = A$ und $D \cup E = B$.

Weiter ziehen wir jeweils Untermengen der Zielmenge "raus": Entsteht eine Gleichung, deren Früchtemenge eine Untermenge der Zielmenge ist, entfernen wir die Gleichung vom Gleichungssystem. Wir verwalten dafür eine Insgesamt-Gleichung, die eine Vereinigung aller solcher Untermengen darstellt und sich immer weiter der Zielmenge annähert. Diese wird in jedem Schritt der ersten Stufe von allen anderen Gleichungen abgezogen.

Die übrig bleibenden Gleichungen werden in der zweiten Stufe miteinander vereinigt, insofern die Früchtemenge der entstehenden Gleichung "näher" an der Zielmenge ist, also neue gewünschte Früchte hinzugekommen sind. Dies geschieht so lange, bis keine neuen Vereinigungen mehr gebildet werden können.

Nach den gegebenen Präferenzen werden dann von allen so entstandenen Gleichungen nach einem Vergleich der Früchtemenge mit der Zielmenge nur die "Besten" ausgewählt.

## Umsetzung

Das Programm akzeptiert bis zu 64 Früchte / Schalen und benötigt entsprechend keine Verschiedenen Anfangsbuchstaben.

### Kompilierung

Go 1.13 oder neuer benötigt. `go build` (erzeugt eine ausführbare Datei namens `Spießgesellen`) oder `go build main.go` (nennt die Datei `main`). Die vorkompilierte ausführbare Datei ist für Linux.

### Verwendung

`go run main.go <pfad> [praeferenzen]` oder `./main <pfad> [praeferenzen]`: Erstes Argument ist der Pfad zur Datei. Das zweite Argument ist optional und gibt Donalds "Präferenzen" an; Standardwert ist `*`. Möglich sind:

* `*`: Bestimme für jedes der folgenden drei Kriterien alle Alternativen
* `+-`: Bestimme Alternativen nach Kriterium 1
* `+`: Bestimme Alternativen nach Kriterium 2
* `-`: Bestimme Alternativen nach Kriterium 3

Beispiel: `./main beispieldaten/spiesse0.txt`

### Ausgabeformat

```
Die Früchte befinden sich in den Schalen: [...]
```

oder

```
Keine exakte Bestimmung der Schalen möglich. Alternativen:

Minimale insgesamte Abweichung ... vom Gewünschten (+-):
=A,B,C,... (erwünschte Früchte); +D,E,F,... (unerwünschte Früchte); -G,H,I,... (fehlende Früchte)  :  [...]

Nur ... fehlende und ... unerwünschte Früchte (+):
...

Nur ... unerwünschte und ... fehlende Früchte (-):
...
```

Je nach Präferenz wird nur ein Teil der Alternativen ausgegeben. Die nach dem Gleichheitszeichen aufgeführten Früchte entsprechen Donalds Wünschen; die hinter dem Pluszeichen sind unerwünscht, und die hinter dem Minuszeichen fehlen Donald.

### Bibliotheken

* `bufio`: `Scanner` zum zeilenweisen Einlesen
* `fmt`: Ausgabe & Formattierung
* `os`: Programmargumente
* `strconv`: Zahl-Text-Konvertierung
* `strings`: Auftrennen von Text

### `BitSet`

> Es gibt höchstens 26 verschiedene Obstsorten.

Entsprechend enthält eine Menge von Früchten, und genauso eine Menge von Schalen, maximal 26 verschiedene Elemente.
Dann kann die Menge aber als "Bit-Set", als "Menge von Bits", dargestellt werden: Jedes Bit steht für eine Frucht bzw. Schale; ist es gesetzt, ist die Frucht oder Schale Teil der Menge, ansonsten nicht. 32 Bits, 4 Bytes, reichen für 26 mögliche Elemente; auf einem 64-Bit System kostet es aber fast keine Laufzeit, auf 64 Bits zu erhöhen.

BitSets sind nicht nur besonders speichereffizient, sondern auch effizient zu verarbeiten: Eine Vereinigung ist über bitweises Oder möglich, ein Schnitt über bitweises And; bitweises Xor stellt die symmetrische Differenz dar, bitweises "A und nicht B" kann für die Differenz $A \ B$ verwendet werden.

#### Kardinalitätsbestimmung

Die Zählung eines Bit-Sets ist die einzige Operation, die nicht mit ein oder zwei simplen bitwise-operators erledigt werden kann.
Ein naiver Ansatz wäre, alle Bits bis zur Anzahl an Früchten durchzugehen und zu zählen. Effizienter ist es allerdings, das Bit-Set in Bytes, 8-Bit-Sets, zu unterteilen, und ein "lookup"-Array zu verwenden: Für alle 256 möglichen 8-Bit-Sets bestimmt man die gesetzten Bits und speichert dies in einem Array. Index ist dann immer ein Bit-Set als 8-Bit-Unsigned-Integer. Für jede der 8-Bit Teilmengen bestimmt man nun über das Array die Anzahl der Elemente und summiert schließlich.

### Gleichungssystem

Das Gleichungssystem wird als "Map" vom BitSet der Früchte zum BitSet der Schalen dargestellt. So können effizient ($O(log n)$) Duplikate verhindert werden.

### Parallelisierbarkeit

Das Verfahren ist schlecht parallelisierbar: Es gibt ein Mengen-Gleichungssystem, auf dem ständig operiert wird; dieses zwischen Threads aufzuteilen wäre mit einem erheblichen Aufwand verbunden. Eine Parallelisierung ist aufgrund der festgestellten Laufzeiten für die Beispiele allerdings nicht unbedingt notwendig.

### Ablauf

Zuerst werden die Programmargumente eingelesen und validiert. Dann wird die Datei eingelesen: In der Reihenfolge des Auftauchens wird jeder Frucht inkrementell eine Nummer zugeordnet, sodass sich die Früchtemenge als Menge von Zahlen von 0 bis Anzahl Früchte (letzteres exklusive) darstellen lässt. Bei allen Operationen müssen entsprechend nur die untersten Anzahl Früchte Bits berücksichtigt werden. Dies wird genutzt, um die Zählung kleiner Mengen weiter zu vereinfachen.

Da die nachgestellten Leerzeichen bei den gegebenen Beispielen habe ich zwar entfernt, dennoch enthält das Programm ein "workaround" (annotiert als "HACK") um diese zu ignorieren: Leere Strings werden nach dem Auftrennen nach Leerzeichen einfach ignoriert.

Es ist möglich, dass nicht alle Früchte benannt werden. In diesem Fall muss aufgefüllt werden, denn es ist möglich, dass am Ende eine Menge gebildet werden muss, in der die unbenannten Früchte vorkommen, ausgehend von der "Insgesamt-Gleichung", in der alle Früchte vorkommen. Für eine nützliche Ausgabe auch in diesem Fall werden die unbekannten Früchte einfach durchnummeriert. Siehe `spiesse8.txt`.

Sobald die Zielmenge als `BitSet` und die Gleichungen als `map[BitSet]BitSet` vorliegen, kann das Verfahren beginnen.

Beide Stufen werden jeweils durch eine `while`-Schleife (in Go: `for` ohne Bedingung) realisiert, die Verlassen wird, sobald keine neuen Einträge der Gleichungs-Map hinzugekommen sind. Nach der ersten Stufe werden alle Gleichungen ohne Schnitt mit der Zielmenge von der Gleichungs-Map entfernt.

Gibt es mehrere Gleichungen mit gleichem Schnitt mit der Zielmenge, können und sollten wir jene mit minimal vielen unerwünschten Früchten bestehen lassen und alle anderen entfernen.

## Quellcode

### `main.go`

```file:go
main.go
```

## Beispiele

Die verstrichene Zeit wurde nicht vom Programm selber, sondern von einem Shellscript gemessen.

### `spiesse0.txt`

Zusätzliches Beispiel (der Aufgabenstellung entnommen):

```file:
beispieldaten/spiesse0.txt
```

```file:
loesungen/spiesse0.txt
```

Dieses Beispiel dient als "sanity-check": Es kann leicht geprüft werden, ob das Ergebnis der per Hand ermittelten Menge der Schalen ${1, 3, 4}$ gleicht.

### `spiesse1.txt`

```file:
loesungen/spiesse1.txt
```

### `spiesse2.txt`

```file:
loesungen/spiesse2.txt
```

### `spiesse3.txt`

```file:
loesungen/spiesse3.txt
```

### `spiesse4.txt`

```file:
loesungen/spiesse4.txt
```

### `spiesse5.txt`

```file:
loesungen/spiesse5.txt
```

### `spiesse6.txt`

```file:
loesungen/spiesse6.txt
```

### `spiesse7.txt`

Zusätzliches Beispiel (`spiesse6.txt` ohne die letzte Beobachtung):

```file:
beispieldaten/spiesse7.txt
```

```file:
loesungen/spiesse7.txt
```

Das so entstandene Beispiel ist eine komplexere Variante von `spiesse3.txt`, um sicherzustellen, dass auch im Fall einer Bestimmung von Alternativen mit mehr Früchten als in `spiesse3.txt` die Laufzeit kurz bleibt.

### `spiesse8.txt`

Zusätzliches Beispiel (`spiesse2.txt` mit Anzahl Früchte 11) - demonstriert die Notwendigkeit des Auffüllens. Als kleiner "stress test" wird die Fruchtzahl außerdem auf 42 gesetzt:

```file:
beispieldaten/spiesse8.txt
```

```file:
loesungen/spiesse8.txt
```

Für alle gegebenen Beispiele ist das Programm ausreichend schnell: In meiner Umgebung sind alle Laufzeiten kürzer als 10 Sekunden.