# Overview
API for practicecalendar.com
The functions described here are meant to be deployed as GCP Cloud Functions.

# Functionality
The HTTP `POST` and `GET` methods are supported. Each of the implemented functions have authentication with `JWT` enabled. Based on the `sub` **Claim**, the data is fetched from Firestore and it is returned in JSON format to the client calling it, after the `JW Token` has been verified.