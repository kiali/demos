# Travels Demo - Version 2
A Microservices demo based on Istio Service Mesh tool. 

This goal of this demo is to demostrate [Istio](https:/istio.io/) cappabilities observed and managed by [Kiali](https://kiali.io) tool.

## Platform Install

This demo has been tested using [Minikube](https://istio.io/latest/docs/setup/platform-setup/minikube/) and [Istio 1.6 Demo Profile](https://istio.io/latest/docs/setup/install/istioctl/#install-a-different-profile)

## Quick Start

Create `travel-agency`, `travel-portal` and `travel-control` namespaces. Add `istio-injection` label and deploy demo app. 

```yaml
kubectl create namespace travel-agency
kubectl create namespace travel-portal
kubectl create namespace travel-control

kubectl label namespace travel-agency istio-injection=enabled
kubectl label namespace travel-portal istio-injection=enabled
kubectl label namespace travel-control istio-injection=enabled

kubectl apply -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_agency.yaml) -n travel-agency
kubectl apply -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_portal.yaml) -n travel-portal
kubectl apply -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_control.yaml) -n travel-control
```

Open Kiali dashboard:

```bash
istioctl dashboard kiali
```

Expose `travel-control` service to your localhost machine:

```bash
kubectl port-forward svc/control 8080:8080 -n travel-control
```

Open [Travels Dashboard](http://localhost:8080).

Undeploy the example:
```yaml
kubectl delete -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_agency.yaml) -n travel-agency
kubectl delete -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_portal.yaml) -n travel-portal
kubectl delete -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_control.yaml) -n travel-control

kubectl delete namespace travel-agency
kubectl delete namespace travel-portal
kubectl delete namespace travel-control
```

## Travels Demo Design

This demo creates two groups of services to simulate a travel portal scenario.

### Travel Portal Namespace

In a first namespace called **travel-portal** there will be deployed several travel shops, where typically users access to search and book flights, hotels, cars or insurances.

There will be several shops to simulate that business of every portal may be different. There will be different characteristics of types of travels, channel (web, mobile), or target (new or existings customers).

These workloads may generate different type of traffic to imitate different real scenarios.

All the portals use a service called *travels* deployed in the **travel-agency** namespace.    
  
### Travel Agency Namespace

A second namespace called **travel-agency** will host a set of services created to provide quotes for travels.

A main *travels* service will be the business entry point for the travel agency. It receives a destination city and a user as parameters and it calculates all elements that compose a travel budget: airfares, lodging, car reservation and travel insurances.

There are several services that calculates a separate price and the travels service is responsible to aggregate them in a single response.

Additionally, some users like *registered* users can have access to special discounts, managed as well by an external service.

Service relations between namespaces can be described in the following diagram:

![Design](doc/Preliminary-Design.png)

## Travels Dashboard

In the Travels Demo Version 2 it has been introduced a **business dashboard** with two features:

- Allow to change the settings of every travel shop simulator. Allowing to change the characteristics of the traffic (ratio, device, user, type of travel).
- Providing a **business** view of the total requests generated from the **travel-portal** namespace to the **travel-agency** services, organized with business criteria as grouped per shop, per type of traffic and per city. 

![Travels Dashboard](doc/Travels-Dashboard.png)

## Kiali Dashboard

Travels Demo Version 2 has introduced a **database** in the example.

A typical flow consists of following steps:

1. A portal queries the travels service for available destinations.
2. Travels service queries the available hotels and return to the portal shop.
3. A user selects a destination and a type of travel, which may include a *flight* and/or a *car*, *hotel* and *insurance*.  
4. *Cars*, *Hotels* and *Flights* may have available discounts depending of user type. 

All traffic and relation between services and workloads can be visualized using Kiali:

![Kiali Dashboard](doc/Kiali-Travel-Graph.png)

## Feedback

This demo is a pet project but if you think is useful to simulate some scenario, test some use case or you miss some feature, please, feel free to provide us feedback.

In any way: from a comment or even a change in the repo.

Thanks !


