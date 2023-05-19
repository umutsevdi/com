#!/bin/bash
css=./content/static/css/style.css
htmls=`ls ./content/*.html | grep -v -- "-min"`
# sed 's/<!--.*-->//g'
cat content/static/css/style.css | sed 's/\/\*.*\*\///g'  | tr -d "\n" | sed 's/ \{2,\}//g' > content/static/css/style-min.css

for i in $htmls; do
    cat $i | tr "\n" " " | sed 's/ \{3,\}/ /g' > `sed "s/\(.*\)\./\1-min\./g" <<< $i`
done
