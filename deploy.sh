#!/bin/bash

gcloud functions deploy practicecalendar-backend \
--memory=128MB \
--entry-point=F \
--env-vars-file=.env.yaml \
--region=europe-west1 \
--runtime=go111 \
--trigger-http \
--verbosity debug
