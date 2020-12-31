# gokit-todo

![GitHub Workflow Status](https://github.com/cage1016/gokit-todo/workflows/ci/badge.svg)
[!["GitHub Discussions"](https://img.shields.io/badge/%20GitHub-%20Discussions-gray.svg?longCache=true&logo=github&colorB=purple)](https://github.com/cage1016/gokit-todo/discussions)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/cage1016/gokit-todo)
[![codecov](https://codecov.io/gh/cage1016/gokit-todo/branch/master/graph/badge.svg)](https://codecov.io/gh/cage1016/gokit-todo)
[![Go Report Card](https://goreportcard.com/badge/cage1016/gokit-todo)](https://goreportcard.com/report/cage1016/gokit-todo)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)

| Service  | Description           |
| -------- | --------------------- |
| todo     | todoMVC backend API   |
| [gokit-todo-frontend](https://github.com/cage1016/gokit-todo-frontend)| todoMVC client        |

![gokit-todo](image.gif)

## Features

- **[Kubernetes](https://kubernetes.io)/[GKE](https://cloud.google.com/kubernetes-engine/):**
  The app is designed to run on Kubernetes (both locally on "Docker for
  Desktop", as well as on the cloud with GKE).
- **[gRPC](https://grpc.io):** Microservices use a high volume of gRPC calls to
  communicate to each other.
- **[Istio](https://istio.io):** Application works on Istio service mesh.
- **[Skaffold](https://skaffold.dev):** Application
  is deployed to Kubernetes with a single command using Skaffold.
- **[go-kit/kit](https://github.com/go-kit/kit):** Go kit is a programming toolkit for building microservices (or elegant monoliths) in Go. We solve common problems in distributed systems and application architecture so you can focus on delivering business value.
- **[todomvc](https://github.com/tastejs/todomvc):** Helping you select an MV* framework

## Motivation

I use [go-kit/kit](https://github.com/go-kit/kit)(microservices toolkit) to build business project and run on Kubernetes and Istio. I also to give few talks about How i use gokit
 - [GoPherCon 2020 TW: 如何透過 Go-kit 快速搭建微服務架構應用程式實戰 ｜ KaiChu](https://kaichu.io/posts/gokit-engineering-operation/)
 - [GDG Cloud Taipei meetup #50 - Build go kit microservices at kubernete…](https://www2.slideshare.net/cagechung/gdg-cloud-taipei-meetup-50-build-go-kit-microservices-at-kubernetes-with-ease-206252668)
- [GDG Devfest 2019 - Build go kit microservices at kubernetes with ease](https://www2.slideshare.net/cagechung/gdg-devfest-2019-build-go-kit-microservices-at-kubernetes-with-ease)
- [cage1016/gokit-workshop](https://github.com/cage1016/gokit-workshop)

Go kit microservices toolkit include many microservices components itself (auth, circuitbreaker, ratelimit etcs.). We build microservices application and deploy to Kubernetes. We could drop those microservices components out with Service Mesh soluation (Istion Envoy proxy) and keep single microservice core business logic clear without any infra codes.

This demo project is [todomvc/gokit-todo-frontend](https://github.com/cage1016/gokit-todo-frontend) backend API implemented by gokit microservice tookit with best practice by myself.

## Goals
- **Project**: quick start a new microservice by toolchain [cage1016/gk](https://github.com/cage1016/gk/tree/feature/gokitconsulk8sistio). Team member could follow same guideline to develop microservice quickly.
  - project layout
- **Testing**: how to write `unit` test with table test `go-sqlmock` and `gomock`. `integration` test, `e2e` test with `TestMain` and `docker-compose`
- **DevOps**: setup CI/CD workflow to increase devops lifecycle
  - skaffold
  - github action (CI)
  - Kubernetes
  - Istio
## Install

this demo support `nginx-ingress` or `istio`

### k8s + Istio 

1. Prepare a Kubernetes cluster
2. Install Istio (**1.6.11**)
    ```sh
    istioctl install
    kubectl label namespace default istio-injection=enabled
    ```
3. Install `gokit-todo` & `frontend`
    ```sh
    kubectl apply -f https://raw.githubusercontent.com/cage1016/gokit-todo/master/deployments/k8s-istio.yaml
    ```
4. Set up `GATEWAY_HTTP_URL`
    ```sh
    export INGRESS_HTTP_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
    export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    export GATEWAY_HTTP_URL=$INGRESS_HOST:$INGRESS_HTTP_PORT
    echo $GATEWAY_HTTP_URL
    ```
5. Visit `$GATEWAY_HTTP_URL` to access todomvc with gokit-todo backend API
6. Delete `gokit-todo` & `frontend`
    ```sh
    kubectl delete -f https://raw.githubusercontent.com/cage1016/gokit-todo/master/deployments/k8s-istio.yaml
    ```
7. Uninstall Istio
    ```sh
    istioctl manifest generate | kubectl delete -f -
    kubectl delete namespace istio-system
    kubectl label namespace default istio-injection=enabled --overwrite
    ```

### k8s + nginx-ingress

1. Prepare a Kubernetes cluster
2. Install `nginx-ingress` CRD by helm3
    ```sh
    helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
    helm repo update
    helm install ingress-nginx ingress-nginx/ingress-nginx
    ```
3. Install `gokit-todo` & `frontend`
    ```sh
    kubectl apply -f https://raw.githubusercontent.com/cage1016/gokit-todo/master/deployments/k8s-nginx-ingress.yaml
    ```
4. Set up `INGRESS_HTTP_URL`
    ```sh
    export INGRESS_HTTP_PORT=$(kubectl get ingress frontend-ingress -o jsonpath='{.spec.rules.*.http.paths.*.backend.servicePort}')
    export INGRESS_HOST=$(kubectl get ingress frontend-ingress -o jsonpath='{.status.loadBalancer.ingress.*.hostname}')
    export INGRESS_HTTP_URL=$INGRESS_HOST:$INGRESS_HTTP_PORT
    echo $INGRESS_HTTP_URL
    ```
5. Visit `$INGRESS_HTTP_URL` to access todomvc with gokit-todo backend API
6. Delete `gokit-todo` & `frontend`
    ```sh
    kubectl delete -f https://raw.githubusercontent.com/cage1016/gokit-todo/master/deployments/k8s-nginx-ingress.yaml
    ```
7. Uninstall nginx ingress
    ```sh
    helm uninstall ingress-nginx
    ```  

## Testing

1. `Makefile`
    ```sh
    $ make
      down                           docker-compose down
      generate                       Regenerates GRPC proto and gomock
      help                           this help
      mod                            tidy go mod
      run                            docker-compose stop & up
      stop                           docker-compose stop
      test                           test: run unit test
      test-e2e                       test-e2e: run e2e test
      test-integration               test-integration: run integration test
      up                             docker-compose up
    ```
2. unit test
    ```sh
    make test
    ```
3. unit integration test by docker-compose
    ```sh
    make test-integration
    ```
4. unit e2e test by docker-compose
    ```sh
    make test-e2e
    ```
## License

Copyright © 2020 [kaichu.io](https://kaichu.io/).<br />
This project is [MIT](https://github.com/cage1016/gokit-todo/blob/master/LICENSE) licensed.