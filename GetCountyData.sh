#!/bin/bash

cd /opt/COVID/healthmry

export DT=`date "+%Y-%m-%d"`

wget --output-document ca.${DT}.html https://www.latimes.com/projects/california-coronavirus-cases-tracking-outbreak/reopening-across-counties/

egrep "    <button" ca.${DT}.html > ca.${DT}.json

