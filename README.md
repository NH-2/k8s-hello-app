# Test different scenarios with the hello-app
This document serves as a guide for exploring the ingress, network policies and Autoscaling functionalities of Kubernetes.


## Set up the Kubernetes Cluster
To begin, we'll be setting up a local "playground" cluster using [minikube](https://minikube.sigs.k8s.io). 


### Start a minikube cluster
Follow the installation instructions for your system [here](https://minikube.sigs.k8s.io/docs/start/). 

Once installed, start cluseter with cni enabled and configured for cilium :
```bash
minikube start --cni=cilium
```

### Mapping a Domain to the Cluster
To demonstrate ingress capabilities, a domain name is needed. Fhis can be achieved by adding an entry to the hosts file  (e.g `/etc/hosts` on linux).
```bash 
sudo echo "$(minikube ip) nabilhouidi.io"   >> /etc/hosts
```

### Optional: install cilium-cli

Optionally, install cilium-cli, the utility CLI for Cilium, our chosen CNI implementation, responsible for managing networking within the cluster and enforcing Network Policies. Follow the installation instructions [here](https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/#install-the-cilium-cli).

Verify the CNI installation and configuration using `cilium-cli`:
```bash
cilium status
```
The output would look similar to this:
```
    /¯¯\
 /¯¯\__/¯¯\    Cilium:             OK
 \__/¯¯\__/    Operator:           OK
 /¯¯\__/¯¯\    Envoy DaemonSet:    disabled (using embedded mode)
 \__/¯¯\__/    Hubble Relay:       disabled
    \__/       ClusterMesh:        disabled

Deployment             cilium-operator    Desired: 1, Ready: 1/1, Available: 1/1
DaemonSet              cilium             Desired: 1, Ready: 1/1, Available: 1/1
Containers:            cilium             Running: 1
                       cilium-operator    Running: 1
Cluster Pods:          7/7 managed by Cilium
Helm chart version:    
Image versions         cilium             quay.io/cilium/cilium:v1.12.3@sha256:30de50c4dc0a1e1077e9e7917a54d5cab253058b3f779822aec00f5c817ca826: 1
                       cilium-operator    quay.io/cilium/operator-generic:v1.12.3@sha256:816ec1da586139b595eeb31932c61a7c13b07fb4a0255341c0e0f18608e84eff: 1
```

### Enabling Relevant Add-ons
The final step in the setup is to enable the relevant addons for minikube. Ensure both ingress and the metrics server are enabled:

```bash
minikube addons enable ingress
minikube addons enable metrics-server   
```

## Deploy the App
With the Kubernetes cluster prepared, we can now deploy the application. The Kubernetes manifest files are located in the  `./kubernetes` directory.

```bash
 kubectl apply -f kuberenetes       
 ```

Inspect the deployed resources:
```bash
kubectl get all -n hello-application
```
which print an output similar to : 
```
NAME                             READY   STATUS    RESTARTS   AGE
pod/hello-app-55d4dc48bf-87zmp   1/1     Running   0          21m

NAME                        TYPE           CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
service/hello-app-service   LoadBalancer   10.98.35.20   <pending>     80:32064/TCP   20m

NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello-app   1/1     1            1           21m

NAME                                   DESIRED   CURRENT   READY   AGE
replicaset.apps/hello-app-55d4dc48bf   1         1         1       21m

NAME                                                REFERENCE              TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
horizontalpodautoscaler.autoscaling/hello-app-hpa   Deployment/hello-app   0%/50%    1         5         1          21s
```

## Testing Scenarios
### Testing Scaling
To evaluate the ability of our deployment setup to scale, we'll apply an artificial load to the server, increasing CPU usage beyond the default replicas' capacity.


As per the scaling policy defined in `kubernetes/hpa.yaml`, when the CPU load surpasses 50%, the Horizontal Pod Autoscaler kicks in and new replicas are created to handle the new load.


1. Ensure the app is 1 replica and ready
```bash
kubectl get pods -n hello-application
kubectl get hpa -n hello-application
```
2. Run the load testing command
```bash
while sleep 0.000000000000001; do wget -q -O- http://nabilhouidi.io/hello-app; done
```
3. Monitor cpu usage
```bash
kubectl get hpa -n hello-application -w
```
4. Watch deployment count. In a seperate terminal
```bash
kubectl get pods -n hello-application -w
```
5. Terminate the load generating loop and close additional terminals using `CTRL+C` .


### Testing Network Policies


The network policies defined in `kubernetes/networkpolicy.yaml` denies all access to the service except through the defined ingress.


To test this, we can try reaching the pod through different methods that should work by default: through the nodeport, and from another pod in the cluster.

** Accessing via NodePort **
In minikube, the `hello-app-service` service created of type LoadBanacer is listening on a NodePort as well.
Attempt to access the application through this service:

1. Get the services's Port and the node's IP
```bash
kubectl get svc -n hello-application
minikube ip
```
2. Try sending a request: 
```bash
curl $(minikube ip):<node port>
```


Note the lack of response due to the network policy blocking the network flow.


** Accessing Pod from Another Namespace **

By default, kubernetes allows all pods to reach all other pods through either the pod's ip or the service. Our Network Policies should block this.


To test this, Let's spin up a new pod in the default namespace, ssh into it, and try sending a request. 

```bash
kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh
```
Once inside the busybox container, send a request using the service's qualified domain: (`{servicename}.{namespaceName}`:
```sh
 wget -O- hello-app-service.hello-application
```
Observe the lack of response due to the network policy blocking requests.