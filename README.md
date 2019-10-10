Practice project where Rust, F# and Go meet. Think implementing popular online store backend, using A.n as inspiration.


Should have:
 - Tracing
 - Monitoring
 - Integration testing
 - Local dev environment
 - Independently scalable services(e.g. black friday)
 - Versioning
 - Health checking


Nice to have:
 - Automated deployments of master branch
 - OLAP to analyze demand - predict traffic and most popular products
 - High load shimulation


API Gateway (Go)
 - Authentication
 - Authorization
 - Proxy for UI


Account
 - Orders history
 - Order info


Store front
 - Availability?
 - Expected at?


Cart
 - Check
 - Place order


Warehouse
 - Availability
 - Upcoming shipments


Shipping
  - Ship before date


Billing & returns
  - Bill on shipping
  - Handle failed billing
  - Track


DB
 - Products info: name, price


Email service
 - Order created
 - Order billed
 - Order shipped


Background worker
 - Update release date availability
 - Update warehouse availability



=======================
Tooling?


skaffold:
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/v0.39.0/skaffold-linux-amd64 \
&& chmod +x skaffold \
&& sudo mv skaffold /usr/local/bin

minikube:
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube
mk start --vm-driver kvm2

local registry:
docker run -d -p 5000:5000 --restart=always --name registry registry:2

https://kind.sigs.k8s.io/docs/user/quick-start/

KUBECONFIG="$(kind get kubeconfig-path --name="kind")" skaffold dev
