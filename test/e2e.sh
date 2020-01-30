#!/usr/bin/env bash

failures=0

for test in $(find ./examples/ -maxdepth 1 -mindepth 1 -type f); do
  echo "Testing $test"
  # get output of the template
  templateOut=$(oc process --local -f $test -o yaml | yq -S .items[])
  # convert template to chart
  template2helm convert --template $test --chart /tmp/charts > /dev/null
  # find newly created chart
  chart=$(ls -td /tmp/charts/*/ | head -1)
  echo "Resulting chart: $chart"
  # get output of chart
  chartOut=$(helm template $chart | yq -S .)
  # compare resources produced
  gap=$(diff <(echo "$templateOut") <(echo "$chartOut"))
  if [[ "${gap}x" != "x" ]]; then
    >&2 echo "Test Failed!"
    >&2 echo "${gap}"
    failures=$((failures+1))
  else
    echo "Passed!"
  fi

  echo
  echo
done

if [[ $failures > 0 ]]; then
  echo "$failures Tests Failed"
  exit $failures
fi
