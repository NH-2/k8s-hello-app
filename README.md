# Test different scenarios with the hello-app
This document is a guide for exploring the ingress, network policies and Autoscaling functionalities of Kubernetes.


## Set up the Kubernetes Cluster
For a start, we'll be setting up a local "playground" cluster using [minikube](https://minikube.sigs.k8s.io). 


### Start a minikube cluster
Follow the installation instructions for your system [here](https://minikube.sigs.k8s.io/docs/start/). 

Once installed, start cluseter with cni enabled and configured for cilium :
```bash
minikube start --cni=cilium
```

### Have a domain point to the cluster
To demonstrate ingress capabilities, a domain name is needed. For this, an entry can be added to the hosts file (e.g `/etc/hosts` on linux).
```bash 
sudo echo "$(minikube ip) nabilhouidi.io"   >> /etc/hosts
```

### install cilium
Next, We need to install Cilium. Cilium will be responsible for handling the networking inside the cluster and making sure out NetWorkingPolicies are implemented and respected.

For that, first  install `cilium-cli`, the helper cli for cilium (our CNI implementation).

follow the installation instructions [here](https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/#install-the-cilium-cli).

Then, you can install cilium on the cluster 
```bash
cilium install --version 1.15.5
```

### enable relevant Addons
The final step in the setup is to enable the relevant addons for minikube. You must enable ingress and the metrics server (used by the HorizontalPodAutoscaler)
```bash
minikube addons enable ingress
minikube addons enable metrics-server   
```

## Deploy the app
Now that the kubernetes cluster is ready. The application can be deployed to it. the file manifests are under the `./kubernetes` directory.

```bash
 k apply -f kuberenetes       
 ```

then, we can inspect the resources we deployed:
```bash
kubectl get all -n hello-application
```
which print an output similar to : 
```
NAME                             READY   STATUS   RESTARTS   AGE
pod/hello-app-55d4dc48bf-kmk8k   0/1     Error    1          41h

NAME                        TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/hello-app-service   LoadBalancer   10.109.251.191   <pending>     80:31114/TCP   23h

NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello-app   0/1     1            0           46h

NAME                                   DESIRED   CURRENT   READY   AGE
replicaset.apps/hello-app-55d4dc48bf   1         1         0       46h

NAME                                                REFERENCE              TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
horizontalpodautoscaler.autoscaling/hello-app-hpa   Deployment/hello-app   0%/50%    1         5         1          46h
```

## Testing Scenarios
### Testing scaling
To test the ability of our deployment setup to scale, an artificial load will be directed to the server, spiking the CPU usage of the single replica present by default.
As per the scaling policy defined in `kubernetes/hpa.yaml`, when the CPU load surpasses 50%, the HPA kicks in and new replicas are created to handle the new load.
1. Ensure the app is 1 replica and ready
```bash
kubectl get pods -n hello-application
kubectl get hpa -n hello-application
```
2. Run the load testing command
```bash
while sleep 0.0000001; do wget -q -O- http://nabilhouidi.io/hello-app; done
```
3. Monitor cpu usage
```bash
kubectl get hpa -n hello-application -w
```
4. Watch deployment count. In a seperate terminal
```bash
kubectl get pods -n hello-application -w
```
5. `CTRL+C` out of the load generating loop and close the extra terminals

### Testing the network policy
The network policies defined in `kubernetes/networkpolicy.yaml` denies all access to the service except through the defined ingress.
To test this, we can try reaching the pod through different methods that should work by default: through the nodeport, and from another pod in the cluster.

** Reaching from nodeport**
In minikube, the `hello-app-service` service created of type LoadBanacer is listening on a NodePort as well.
Let's try reaching our application through that service.
1. get the node ip and NodePort
```bash
kubectl get svc -n hello-application
minikube ip
curl $(minikube ip):<node port>
```
Notice we get no response because the networkingPolicy blocked the network flow.


** Reaching pod through other namespace**
By default, kubernetes allows all pods to reach all other pods through either the pod's ip or the service. 
To test this, Let's spin up a new pod in the default namespace, ssh into it, and try sending a request. 
```bash
kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh
```
Now that a prompt inside the new busybox container is presented, let's send a request using the service qualified domain (`{servicename}.{namespaceName}`:
```sh
 wget -O- hello-app-service.hello-application
```
Notice how no response is returned because the NetworkingPolicy blocked the network flow 

 
