#!/bin/bash

# URL for the API request
URL="https://power-api.loe.lviv.ua/api/pw_accidents?page=1&otg.id=28&city.id=693&street.id=12783"
#URL="https://power-api.loe.lviv.ua/api/pw_accidents?page=1&otg.id=28&city.id=693&street.id=13968"

# Fetch the data using curl and parse it using jq
RESULT=$(curl -s "$URL" | jq -c '."hydra:member"[] | select(.buildingNames | test("\\b21\\b")) | {
    Accident_ID: .id,
    OTG_Code: .code_otg,
    City: .city.name,
    Street: .street.name,
    Date_Event: .dateEvent,
    Date_Planned_Resolution: .datePlanIn,
    Comment: .koment,
    Buildings_Affected: .buildingNames
}')

if [[ -z $RESULT ]]; then
    echo "No power outages found."
else
    echo "$RESULT" | jq .
fi
