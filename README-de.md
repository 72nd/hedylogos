<div align="center" style="border-bottom: none">
  <h1>
    <img src="misc/logo.png" width="120"/>
    <br>
    hedylogos
  </h1>
  <h2>Build interactive audio scenarios using old dial phones or numeric keyboards</h2>
  <p><a href="README.md">English Version</a></p> 
</div>

Mit Hedylogos ist es möglich, interaktive Audioformate zu entwickeln. Die Navigatio zwischen den einzelnen Audiosschnipsel erfolgt über die Eingabe/das Wählen von Nummern, wie das etwa auch bei Menüs von Telefonhotlines bekannt ist. Die Software unterstützt hierfür zwei unterschiedliche Modi für die Nutzenden: Zum einen kann altes Wählscheibentelefon (nach einem kleinen Umbau) dafür verwendet werden oder die Software kann mithilfe eines [Ziffernblocks](https://de.wikipedia.org/wiki/Ziffernblock) gesteuert werden. Die Erstellung der Audioszenarios ist verhältnismäßig einfach und kann sowohl in einem grafischen Editor als auch per Texteditor erfolgen (mehr dazu unten). Es sind sowohl inhaltlich als auch von der Nutzung her viele unterschiedliche Szenarien denkbar, wie interaktive Geschichten oder als Player in einem Museum, welcher je nach Knopfdruck unterschiedliche Informationen abspielt.


## Förderung

<img alt="Logo des Ministeriums für Wissenschaft Forschung und Kultur des Landes Brandenburg" src="misc/mwfk.png" width="200" style="align:left"/>

Das Projekt wurde gefördert mit Mitteln des Ministeriums für Wissenschaft, Forschung und Kultur des Landes Brandenburg.


## Anforderungen

Da Hedylogos in Python geschrieben ist, läuft die Software auf fast jedem Computer. Bitte beachte, dass für die Verwendung eines Wählscheibentelefons ein [Raspberry Pi](https://en.wikipedia.org/wiki/Raspberry_Pi) notwendig ist, da ein normaler Laptop oder Computer nicht über die notwendigen Schnittstellen verfügt. Wenn du den Aufwand scheust, kannst du aber mit dem Ziffernblock-Modus auch einfach einen herkömmlichen Computer verwenden, an welchen ein dezidierten Ziffernblock (wie er etwa hier bei [Galaxus](https://www.galaxus.de/de/s1/product/logilink-id0120-nummernblock-kabellos-tastatur-12817754) zu kaufen ist) und ein Kopfhörer gehängt wird.

Die Installation von Python ist abhängig vom verwendeten Betriebssystem lässt sich aber mit [dieser Anleitung](https://python.land/installing-python) einfach bewerkstelligen. Da Python weit verbreitet sind, sollten sich etwaige Probleme durch das Verwenden einer Suchmaschine einfach lösen lassen.

Weiter solltest du etwas mit der Bedienung der Kommandozeile vertraut sein. Eine Einführung findest du etwa hier für [Windows](https://www.makeuseof.com/tag/a-beginners-guide-to-the-windows-command-line/) und hier für [MacOS](https://www.makeuseof.com/tag/beginners-guide-mac-terminal/).


## Installation

Lade dir zunächst die [neuste Version](https://github.com/72nd/hedylogos/releases/latest) runter. Entpacke das Archiv und navigiere im Terminal in den Ordner.

```
python -m venv .venv
pip install .
```

Danach kann Hedylogos mit dem Befehl `hedylogos` gestartet werden.


## Beispiel

Im Ordner `example` findet sich ein Beispielszenario, welches die wichtigsten Funktionen und Eigenschaften von Hedylogos abdeckt. Um das Beispiel zu starten, kann folgender Befehl verwendet werden.

```
hedylogos run-keyboard example/scenario.json
```


## Ein Szenario erstellen

Es gibt zwei Möglichkeiten, ein Szenario für Hedylogos zu erstellen. Zum einen gibt es einen grafischen Editor, zum anderen können fortgeschrittene Nutzer auch die JSON Datei direkt bearbeiten. In Folge werden beide Wege erläutert.


### Mit dem Editor

Um das Erstellen der Szenarios möglichst einfach zu gestalten, existiert ein [grafischer Editor](https://72nd.github.io/hedylogos/editor/). Mehr zur Bedienung findest du im Editor selbst. Lade am Ende das Szenario als Datei runter und speichere es auf deinem Computer.


### Manuell in der JSON Datei

Wenn dir JSON nicht fremd ist, kannst du das Szenario auch direkt im Editor erstellen. Die Bedeutung der einzelnen Fälle sollten soweit selbsterklärend sein. Sonst lohnt sich ein Blick in den [grafischen Editor](https://72nd.github.io/hedylogos/editor/) und die darin befindlichen Erläuterungen der Felder zu werfen. Die Software bietet auch die Möglichkeit, ein Template der Datei zu generieren.

```
hedylogos init scenario.json
```


## Die Szenariodatei validieren

Bei der Erstellung größerer Szenarios schleichen sich schnell kleine Fehler und Verschreiber ein. Damit diese nicht erst bei der Benutzung auffallen, sind in Hedylogos Validierungsroutinen eingebaut, welche beim Laden einer Szenariodatei automatisiert die allermeisten Probleme erkennen. Nur die Prüfung, ob alle Audiodateien vorhanden sind, muss manuell mit einem Befehl ausgelöst werden.

```
hedylogos check pfad/zum/szenario.json
```

## Das Szenario abspielen

Hedylogos bietet zwei unterschiedliche Modi an. Zum einen mit einem alten Wählscheibentelefon oder mit einer Tastatur bzw. [Ziffernblocks](https://de.wikipedia.org/wiki/Ziffernblock).


### Mit der Tastatur / Ziffernblock

Die Ausführung wird mit diesem Befehl gestartet.

```
hedylogos run-keyboard pfad/zum/szenario.json
```

Die Ausführung kann mit folgenden Tasten gesteuert werden:

- `p` oder `<ENTER>`: Startet das Szenario. Wenn das Szenario bereits wiedergegeben wird, wird die Wiedergabe abgebrochen und startet wieder am Startpunkt.
- `h`: Stopt die Wiedergabe des Szenarios. Simuliert in erster Linie den Moment, wenn der Telefonhörer aufgelegt wird und hat deshalb für den Tastaturmodus keine wirkliche Bedeutung.
- `q`: Beendet die Ausführung des Programms. Damit Besucher diese Aktion nicht auslösen können, empfiehlt sich die verwendung eines dezidierten Nummernblocks.
- `0-9`: Wählen einer Nummern.


### Mit einem Wählscheibentelefon

Um den Input eines Wählscheibentelefons zu verwenden, wird die Bibliothek [RotaryPi](https://pypi.org/project/rotarypi/) verwendet. Voraussetzung hierfür ist die Verwendung eines [Raspberry Pis](https://www.raspberrypi.org/). Mehr über die Belegung der Pins kann in der [Dokumentation von RotaryPi](https://rotarypi.readthedocs.io/en/latest/) erfahren werden.

