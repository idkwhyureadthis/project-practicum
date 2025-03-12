#!/bin/bash

DBSTRING="postgresql://$DBUSER:$DBPASSWORD@$DBHOST:$DBPORT/$DBNAME?sslmode=$DBSSL"

goose postgres "$DBSTRING" up