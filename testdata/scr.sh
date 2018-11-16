#!/bin/bash

files=`ls ./testdata/*.gen.got`
for entry in $files
do
  echo "$entry"
  `cp $entry $entry.expected`
done

