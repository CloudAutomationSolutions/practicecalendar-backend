#!/bin/bash

gcloud functions deploy test-practicecal-backend \
--memory=128MB \
--entry-point=F \
--env-vars-file=.env.yaml \
--region=europe-west1 \
--runtime=go111 \
--trigger-http