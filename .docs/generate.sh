#!/bin/bash
TIMEFORMAT="Zeit verstrichen: %3lR"
echo "Docs..."
cd ..
# pandoc Allgemein.md -f markdown -t latex --pdf-engine=lualatex -s -o Dokumentation.pdf --include-in-header=.docs/header.tex --lua-filter=.docs/files.lua -V fontsize=12pt -M lang:de
for aufgabe in "Spießgesellen" "Eisbudendilemma" ; do
    echo "$aufgabe"
    cd "$aufgabe"
    echo "Build..."
    go build main.go
    cd beispieldaten
    if [ "$1" == "run" ]; then
        echo "Run..."
        for i in *.txt; do
            filename="../loesungen/${i%.*}.png"
            if [ "$aufgabe" == "Eisbudendilemma" ]; then
                (time (../main "$i" j "$filename" && printf "\n")) > "../loesungen/${i}" 2>&1;
            else
                (time (../main "$i" && printf "\n")) > "../loesungen/${i}" 2>&1;
            fi
        done
    fi
    cd ..
    echo "Docs..."
    pandoc Dokumentation.md -f markdown -t latex --pdf-engine=lualatex -s -o Dokumentation.pdf --include-in-header=../.docs/header.tex --lua-filter=../.docs/files.lua -V fontsize=12pt -M lang:de
    cd ..
done
echo "Zip..."
rm -f "../bwinf39-runde2.zip"
zip -r "../bwinf39-runde2.zip" .
echo "Done"