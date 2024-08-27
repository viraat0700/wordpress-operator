Wordpress Go-Operator

Pre-Requisities :

Dependencies to run Operators on System:

i.Minikube to create single node Kubernetes Cluster on your Local Machine. ✓
ii.Install GoLang on your Local machine. ✓
iii.Install Docker Desktop / Docker CLI on your Local Machine. ✓
iv.Install Operator SDK
v.Install Kubectl CLI to interact with Kubernetes CLI. ✓

Steps to install all the dependencies in a Linux based System:

i.Minikube  For Minikube to work properly we need to install Docker and kubectl before hand.

Step 1: Update Your System
sudo apt update
sudo apt upgrade -y


Step 2: Install Required Packages
sudo apt install -y curl apt-transport-https


Step 3: Install Docker
sudo apt install -y docker.io
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
If facing any issue you may need to follow the documentation by Docker website. I am Pasting the Link for your reference.
// https://docs.docker.com/engine/install/


Step 4: Install Kubectl

curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"

chmod +x kubectl

sudo mv kubectl /usr/local/bin/


Step 5: Install Minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64


chmod +x minikube-linux-amd64
sudo mv minikube-linux-amd64 /usr/local/bin/minikube

Step 6: Start Minikube
minikube start --driver=docker

Step 7: Verify Installation
minikube status

ii.Install GoLang :
To install GoLang on Ubuntu, follow these steps:

Step 1: Update your system
sudo apt update
sudo apt upgrade

Step 2: Download GoLang
wget <https://go.dev/dl/go1.19.4.linux-amd64.tar.gz>

Step 3: Install GoLang
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.19.4.linux-amd64.tar.gz

Step 4: Set up GoLang environment
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile

Step 5: Verify the installation
go version
go version go1.19.4 linux/amd64


iii. Install Operator-sdk :
Install Prerequisites:

sudo apt update
sudo apt install -y curl git gcc make

Install Go:
wget <https://go.dev/dl/go1.20.6.linux-amd64.tar.gz>
sudo tar -C /usr/local -xzf go1.20.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

echo "export PATH=\\$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile

Download and Install Operator SDK:-

 export ARCH=$(case $(uname -m) in x86_64) echo amd64 ;; aa
 rch64) echo arm64 ;; ppc64le) echo ppc64le ;; s390x) echo 
s390x ;; esac)
 export OS=$(uname | awk '{print tolower($0)}')
 export OPERATOR_SDK_DL_URL=https://github.com/operator-fra
 mework/operator-sdk/releases/download/v1.27.0
 curl -LO ${OPERATOR_SDK_DL_URL}/operator-sdk_${OS}_${ARCH}
gpg --keyserver hkps://keys.openpgp.org --recv-key 052996E
2A20B5C7E
curl -LO ${OPERATOR_SDK_DL_URL}/checksums.txt
curl -LO ${OPERATOR_SDK_DL_URL}/checksums.txt.asc
gpg --verify checksums.txt.asc
grep operator-sdk_${OS}_${ARCH} checksums.txt | sha256sum -c



chmod +x operator-sdk_${OS}_${ARCH} && sudo mv operator-sd
 k_${OS}\_${ARCH} /usr/local/bin/operator-sdk

Verify Installation:
operator-sdk version

Now that every dependency is installed you can start with your Operator
development.

Note: Sometimes you may get error with operator-sdk and go. Then you have to
check with the version of these two because not all versions are compatible.


Steps to create Go based Operator :

1.Initialize a new Go operator project:
operator-sdk init --domain=my.domain --repo=github.com/myusername/myoperator

The command operator-sdk init --domain=my.domain --repo=github.com/myusername/my
operator sets up the basic structure for your Go-based operator project.

i.operator-sdk init: This initializes a new operator project. It creates all the
necessary files and folders to start building your operator.

ii.--domain=my.domain: This sets the domain for your custom resources. Think of
it like the unique "address" that identifies the resources managed by your operator
(e.g.,myapp.my.domain ). It's like a namespace for your custom Kubernetes resources.

iii.--repo=github.com/myusername/my-operator: This defines the Go module path for your project. It tells Go where the project's code will live (e.g., in a GitHub repository). This is important because Go needs this path to properly handle dependencies and organize your code.

2.Create API and Controller:

The command operator-sdk create api --group=app --version=v1 --kind=MyApp --resource --controller helps you add a new custom resource CRD and its controller to your operator project.

i.operator-sdk create api: This command generates the code and files for a new custom resource and its controller in your operator project.

ii.--group=app: This sets the API group for your custom resource. Think of it as a category or grouping for related resources. In this case, all resources in this group will belong to the app category.

iii.--version=v1 This specifies the version of the API for your custom resource. Versioning helps you manage changes and updates to your resource definitions over time.

iv.--kind=MyApp: This defines the name of your custom resource. In this case,
you're creating a resource called MyApp . This is the type of resource your operator will manage (eg., like a "Deployment" or "Service" in Kubernetes).

v.--resource: This flag tells the operator to generate the CustomResourceDefinition CRD code. The CRD defines what your custom resource looks like in Kubernetes.

vi.--controller: This flag generates a controller, which is the code that watches and manages the lifecycle of your custom resource. The controller ensures that the actual state of your resource matches the desired state defined by the user.


To Deploy on the cluster
Build and push your image to the location specified by `IMG`:

```sh
 make docker-build docker-push IMG=<some-registry>/wordpress-oper
```

NOTE:This image ought to be published in the personal registry y
And it is required to have access to pull the image from the wor
Make sure you have the proper permission to the registry if the
Install the CRDs into the cluster:

```sh
make install
```

Deploy the Manager to the cluster with the image specified by `I

```sh
make deploy IMG=<some-registry>/wordpress-operator:tag
```

NOTE: If you encounter RBAC errors, you may need to grant yourse
privileges or be logged in as admin.
Create instances of your solution
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

NOTE: Ensure that the samples has default values to test it out
To Uninstall
Delete the instances (CRs) from the cluster:

```sh
kubectl delete -k config/samples/
```

Delete the APIs(CRDs) from the cluster:

```sh
make uninstall
```

UnDeploy the controller from the cluster:

```sh
make undeploy
```