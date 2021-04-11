#!/bin/sh

for i in $(ls po/*.pot)
do
    p=$(basename ${i})
    base=${p%.pot}
    mkdir -p po/${LANG}
    msginit --no-translator --no-wrap --input=po/${base}.pot --locale=${LANG} --output=po/${LANG}/${base}.po
done
