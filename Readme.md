# Multicloud Kubernetes Cluster Setup
**1.1** The pipeline only works if the master node is deployed on Amazon. If the master node is on another cloud, the pipeline needs adjustments. 
**1.2** Create your virtual machines on any cloud platform, but at least the main node on AWS. Assign to all machines static public ip's. 
## To be done on every virtual machine
**2.1** Disable swap 
 ```bash
swapoff -a
```
**2.2** Install docker. For Ubuntu see: https://docs.docker.com/engine/install/ubuntu/ 
INFO: Make sure to use the same versions on every node. I had an issue with Ubuntu 22.04, where the containers did not stabilize. 
Ubuntu 20.04 works fine. Verified functionality with docker and cli version 20.10.13, containerd version 
**2.3** Enable docker
 ```bash
systemctl enable docker --now
```
**2.4** Install kubernetes: kubelet, kubeadm & kubectl: https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/ 
**2.5** Set cgroupdriver as systemd
```bash
cat > /etc/docker/daemon.json
{ "exec-opts": ["native.cgroupdriver=systemd"] }
```
**2.6** Restart services
```bash
 systemctl restart docker
 systemctl enable --now kubelet
```
**2.7** Install iproute if not already installed. On ubuntu:
```bash
apt-get install iproute2
```
**2.8** To enable pod-to-pod connection, network bridges have to be enabled for kubernetes
```bash
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
```
**2.9** Install openvpn or any other vpn, but this tutorial uses openvpn. For ubuntu:
```bash
apt-get install openvpn
```
## To be done on the master node in AWS
**3.1** clone the easy-rsa github repo: 
```bash
git clone https://github.com/OpenVPN/easy-rsa.git
```
**3.2** navigate into the easyrsa3 subfolder 
```bash
cd easy-rsa/easyrsa3
```
Initialize the CA on the master node with: 
```bash
 ./easyrsa init-pki
 ./easyrsa build-ca
```
Now for every vpn client use a different 'EntityName' and generate the private keys and certificates for it with:\\
NOTE: If you enter a pass-phrase, you have to remember it 
```bash
 ./easyrsa gen-req EntityName
 ./easyrsa sign-req client EntityName
```
And a server certificate for the master-node 
```bash
 ./easyrsa gen-req EntityName
 ./easyrsa sign-req server EntityName
```
Now create a ta.key for TLS and dh-params with 
```bash
  ./easyrsa gen-dh
 openvpn --genkey --secret ta.key
```
You will need the created files for the workers later on 
Get the server.conf in our github repo under /vpn/server.conf and specify the right paths to the server certificate, private key, ta.key, dh params and
ca cert in line 78-85 and line 244.

**3.3** Now you can start the vpn server on the master node under sudo. Specify the right path to the server.conf. 
You will be prompted to enter a passphrase for the private key of the server, if you have specified one. 
```bash
openvpn --config /path/to/server.conf --daemon
```
**3.4** Pull the kubernetes api images with 
```bash
kubeadm config images pull
```
**3.5** Before initializing the cluster, you have to specify the node-ip to be used in the cluster. This has to be the vpn ip. 
For ubuntu this is done in the file under the path "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf" 
Check with ifconfig, the ip of the vpn interface tun0. For the master node it should be 10.8.0.1. Add the 
--node-ip=10.8.0.1 flag to ExecStart=... --node-ip=10.8.0.1. 
NOTE: Use the last ExecStart= in the config not the empty one. 

For Amazon Linux 2 you can find the file under the path /etc/sysconfig/kubelet and write to the file the following (for master otherwise change ip) 
KUBELET_EXTRA_ARGS='--node-ip 10.8.0.1' 
**3.6** If you did not change the subnet in the vpn server.conf, then you can initialize the cluster with the following:
```bash
kubeadm init --control-plane-endpoint "PUBLIC-IP:6443" --apiserver-advertise-address=10.8.0.1 --pod-network-cidr=192.168.0.0/16
``` 
The control-plane-endpoint flag with the public ip of the master node is only needed if the pipeline is used, otherwise you can skip this 
option, but you have to manually deploy the application. 
Remember the kubadm join command, this used used on the other nodes to join the cluster.

**3.7** Copy the admin conf and modify the access rights to be able to access the cluster without the sudo user:
```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
``` 
**3.8** The relay-pod is deploying the enclave and needs a cluster-admin role to complete this task. This is done by the following command:
 ```bash
kubectl create clusterrolebinding serviceaccounts-cluster-admin 
  --clusterrole=cluster-admin 
  --group=system:serviceaccounts
``` 
**3.9** Now the last step on the master node is to apply Calico-networking to the cluster. You can find the yaml file in our repo under /deployments/networking/calico.yaml 
Specify the right path before entering
```bash
kubectl apply -f /path/to/calico.yaml
``` 
This file is edited to work in a VPN setting. It automaticly selects the tun0 interface on every node. If the vpn network interface should be 
another interface, then you have to edit the calico.yaml file.
## To be done on the worker nodes
**4.1** The ta.key has to be the same on every node. Distribute it to every vpn client with the client certificate, CA certificate, 
the client private key and the client.conf. The client.conf can be found in our github repo under /vpn/client.conf. You have to  
edit the client.conf in line 90-92 to specify the right paths to the files. Also in line 111 the path to the ta.key has to be specified. 
**4.2** edit the hosts file with the sudo user 
```bash
nano /etc/hosts
``` 
apply to the file: 
public-ip master-node 
where public-ip is the public ip of the master. 
**4.3** Now you can open the vpn connection to the master node with: 
```bash
openvpn --config /path/to/client.conf --daemon
``` 
**4.4** Before joining the cluster, you have to edit the node ip like in step 3.5. After editing join the cluster with the kubeadm join ... command.

## To be done on the master node

**5.1** Verify, that all nodes and all pods are up with
```bash
kubectl get nodes -o wide 
kubectl get pods -o wide --all-namespaces
``` 
**5.2** The next step is to label the nodes as the workers. 
For every node except the master-node, where nodename equals the name of the node listed in the output of the first command in 5.1, execute the following command:
```bash
kubectl label node nodename node-role.kubernetes.io/worker=worker
``` 
**5.3** To deploy only the enclaves on azure use the following command to additionaly label the azure node or nodes:
```bash
kubectl label node azurenode disktype=ssd
``` 
all other nodes except the master-node should be labeled with: 
```bash
kubectl label node azurenode disktype=db
``` 
**5.4** Now your cluster is all set!
## To be done on AWS and Github if the pipeline is used. 
NOTE: Notice that, if the pipeline is not used, the image names of the dockerfiles and the environment variables have to be manually set to the right values. 

**6.1** The amazon nodes have to pull the images from ECR and need IAM-Roles for that. Create a new IAM Role on aws by searching for IAM>Roles>Create new Role 
Select AWS-Service and as use case EC2 
Click on next and add the following policies: AmazonEC2FullAccess, IAMFullAccess and AmazonEC2ContainerRegistryFullAccess. 
Attach the IAM Role to every instance by navigating to your instances right clicking on it, then select security and change IAM-Role. Select your created IAM-Role. 

**6.2** Now again navigate to IAM and create a new User. This user is needed for the pipeline to access the cloud environment and push the images. Add all the following policies
to the user: AWSAgentlessDiscoveryService, AWSApplicationDiscoveryAgentAccess, AmazonEC2ContainerRegistryFullAccess, AWSMigrationHubDiscoveryAccess,  AWSDiscoveryContinuousExportFirehosePolicy,
AWSCloudMapFullAccess, AWSApplicationDiscoveryServiceFullAccess, AmazonElasticFileSystemFullAccess, AmazonEC2FullAccess, IAMFullAccess, AutoScalingFullAccess, ElasticLoadBalancingFullAccess,
ApplicationAutoScalingForAmazonAppStreamAccess, AmazonAPIGatewayPushToCloudWatchLogs, CloudWatchLogsFullAccess, AmazonECS_FullAccess, AmazonRoute53FullAccess, AWSCloudFormationFullAccess. 
After completion you will get your access key id and the secret. Don't close this window, because otherwise, you will have to create a new user. The secret is only shown once.

**6.3** Navigate to your fork or github repo and add under settings>secrets the access key id with the variable name AWS_ACCESS_KEY_ID and the secret with the variablename AWS_SECRET_ACCESS_KEY. 
Here you can also specify the password of the database, by adding the variablename DB_SECRET and specifying a password. 

**6.4** Now the last step is to add the kubeconfig to the repo as the secret KUBE_CONFIG_DATA. To get the contents of the kubeconfig, navigate the master-node and execute the following command: 
```bash
cat $HOME/.kube/config
``` 

**6.5** Now everything is set and you can execute the pipeline to deploy the application! 

 
