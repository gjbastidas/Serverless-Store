# Backend

Based on Golang AWS SDK and AWS Terraform provider to build a serverless backend.

## Diagram

TODO

## Code organization

- [**environments**](./environments/): contains the backend infrastructure code to deploy each environment.
- [**modules**](./modules/): contains terraform modules
- [**store_apis**](./store_apis/): contains the backend store logic based on Golang

## Prerequisites

This app has been tested with:

- Terraform =1.5.2
- Go =1.19
- GNU Make =3.81
