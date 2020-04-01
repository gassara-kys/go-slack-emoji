#!/bin/bash -e

# move files to "_mv_XXX" directory per 500 files.

function init() {
  rm -rf _mv_*
}

function createDir() {
  DIR="_mv_$1"
  mkdir -p $DIR
}

init
COUNT=0
for file in `\find . -maxdepth 1 -type f`; do
  if [ `basename $file` = ".gitignore" ] \
      || [ `basename $file` = "_mv.sh" ]; then
    continue
  fi
  if [ $(($COUNT % 100)) -eq 0 ]; then
		createDir $COUNT
	fi
  mv $file "$DIR/$file"
  COUNT=$(( COUNT + 1 ))
done
