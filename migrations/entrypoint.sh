#!/bin/bash

DBSTRING="host=postgres user=intern password=123 dbname=wb sslmode=disable"

goose postgres "$DBSTRING" up