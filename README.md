# Accessing an FTP server using Skupper

[![main](https://github.com/skupperproject/skupper-example-ftp/actions/workflows/main.yaml/badge.svg)](https://github.com/skupperproject/skupper-example-ftp/actions/workflows/main.yaml)

#### Securely connect to an FTP server on a remote Kubernetes cluster


This example is part of a [suite of examples][examples] showing the
different ways you can use [Skupper][website] to connect services
across cloud providers, data centers, and edge sites.

[website]: https://skupper.io/
[examples]: https://skupper.io/examples/index.html


#### Contents

* [Overview](#overview)
* [Prerequisites](#prerequisites)
* [Step 1: Configure separate console sessions](#step-1-configure-separate-console-sessions)
* [Step 2: Access your clusters](#step-2-access-your-clusters)
* [Step 3: Set up your namespaces](#step-3-set-up-your-namespaces)
* [Step 4: Install Skupper in your namespaces](#step-4-install-skupper-in-your-namespaces)
* [Step 5: Check the status of your namespaces](#step-5-check-the-status-of-your-namespaces)
* [Step 6: Link your namespaces](#step-6-link-your-namespaces)
* [Step 7: Deploy the FTP server](#step-7-deploy-the-ftp-server)
* [Step 8: Expose the FTP server](#step-8-expose-the-ftp-server)
* [Step 9: Run the FTP client](#step-9-run-the-ftp-client)
* [Accessing the web console](#accessing-the-web-console)
* [Cleaning up](#cleaning-up)

## Overview

This example shows you how you can use Skupper to connect an FTP
client on one Kubernetes cluster to an FTP server on another.

## Prerequisites


* The `kubectl` command-line tool, version 1.15 or later
  ([installation guide][install-kubectl])

* The `skupper` command-line tool, the latest version ([installation
  guide][install-skupper])

* Access to at least one Kubernetes cluster, from any provider you
  choose

[install-kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[install-skupper]: https://skupper.io/install/index.html


## Step 1: Configure separate console sessions

Skupper is designed for use with multiple namespaces, typically on
different clusters.  The `skupper` command uses your
[kubeconfig][kubeconfig] and current context to select the
namespace where it operates.

[kubeconfig]: https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/

Your kubeconfig is stored in a file in your home directory.  The
`skupper` and `kubectl` commands use the `KUBECONFIG` environment
variable to locate it.

A single kubeconfig supports only one active context per user.
Since you will be using multiple contexts at once in this
exercise, you need to create distinct kubeconfigs.

Start a console session for each of your namespaces.  Set the
`KUBECONFIG` environment variable to a different path in each
session.

_**Console for west:**_

~~~ shell
export KUBECONFIG=~/.kube/config-west
~~~

_**Console for east:**_

~~~ shell
export KUBECONFIG=~/.kube/config-east
~~~

## Step 2: Access your clusters

The methods for accessing your clusters vary by Kubernetes
provider. Find the instructions for your chosen providers and use
them to authenticate and configure access for each console
session.  See the following links for more information:

* [Minikube](https://skupper.io/start/minikube.html)
* [Amazon Elastic Kubernetes Service (EKS)](https://skupper.io/start/eks.html)
* [Azure Kubernetes Service (AKS)](https://skupper.io/start/aks.html)
* [Google Kubernetes Engine (GKE)](https://skupper.io/start/gke.html)
* [IBM Kubernetes Service](https://skupper.io/start/ibmks.html)
* [OpenShift](https://skupper.io/start/openshift.html)
* [More providers](https://kubernetes.io/partners/#kcsp)

## Step 3: Set up your namespaces

Use `kubectl create namespace` to create the namespaces you wish
to use (or use existing namespaces).  Use `kubectl config
set-context` to set the current namespace for each session.

_**Console for west:**_

~~~ shell
kubectl create namespace west
kubectl config set-context --current --namespace west
~~~

_Sample output:_

~~~ console
$ kubectl create namespace west
namespace/west created

$ kubectl config set-context --current --namespace west
Context "minikube" modified.
~~~

_**Console for east:**_

~~~ shell
kubectl create namespace east
kubectl config set-context --current --namespace east
~~~

_Sample output:_

~~~ console
$ kubectl create namespace east
namespace/east created

$ kubectl config set-context --current --namespace east
Context "minikube" modified.
~~~

## Step 4: Install Skupper in your namespaces

The `skupper init` command installs the Skupper router and service
controller in the current namespace.  Run the `skupper init` command
in each namespace.

**Note:** If you are using Minikube, [you need to start `minikube
tunnel`][minikube-tunnel] before you install Skupper.

[minikube-tunnel]: https://skupper.io/start/minikube.html#running-minikube-tunnel

_**Console for west:**_

~~~ shell
skupper init
~~~

_Sample output:_

~~~ console
$ skupper init
Waiting for LoadBalancer IP or hostname...
Skupper is now installed in namespace 'west'.  Use 'skupper status' to get more information.
~~~

_**Console for east:**_

~~~ shell
skupper init
~~~

_Sample output:_

~~~ console
$ skupper init
Waiting for LoadBalancer IP or hostname...
Skupper is now installed in namespace 'east'.  Use 'skupper status' to get more information.
~~~

## Step 5: Check the status of your namespaces

Use `skupper status` in each console to check that Skupper is
installed.

_**Console for west:**_

~~~ shell
skupper status
~~~

_Sample output:_

~~~ console
$ skupper status
Skupper is enabled for namespace "west" in interior mode. It is connected to 1 other site. It has 1 exposed service.
The site console url is: <console-url>
The credentials for internal console-auth mode are held in secret: 'skupper-console-users'
~~~

_**Console for east:**_

~~~ shell
skupper status
~~~

_Sample output:_

~~~ console
$ skupper status
Skupper is enabled for namespace "east" in interior mode. It is connected to 1 other site. It has 1 exposed service.
The site console url is: <console-url>
The credentials for internal console-auth mode are held in secret: 'skupper-console-users'
~~~

As you move through the steps below, you can use `skupper status` at
any time to check your progress.

## Step 6: Link your namespaces

Creating a link requires use of two `skupper` commands in
conjunction, `skupper token create` and `skupper link create`.

The `skupper token create` command generates a secret token that
signifies permission to create a link.  The token also carries the
link details.  Then, in a remote namespace, The `skupper link
create` command uses the token to create a link to the namespace
that generated it.

**Note:** The link token is truly a *secret*.  Anyone who has the
token can link to your namespace.  Make sure that only those you
trust have access to it.

First, use `skupper token create` in one namespace to generate the
token.  Then, use `skupper link create` in the other to create a
link.

_**Console for west:**_

~~~ shell
skupper token create ~/secret.token
~~~

_Sample output:_

~~~ console
$ skupper token create ~/secret.token
Token written to ~/secret.token
~~~

_**Console for east:**_

~~~ shell
skupper link create ~/secret.token
~~~

_Sample output:_

~~~ console
$ skupper link create ~/secret.token
Site configured to link to https://10.105.193.154:8081/ed9c37f6-d78a-11ec-a8c7-04421a4c5042 (name=link1)
Check the status of the link using 'skupper link status'.
~~~

If your console sessions are on different machines, you may need
to use `sftp` or a similar tool to transfer the token securely.
By default, tokens expire after a single use or 15 minutes after
creation.

## Step 7: Deploy the FTP server

In the east namespace, use `kubectl apply` to deploy the FTP
server.

_**Console for east:**_

~~~ shell
kubectl apply -f ftp-server
~~~

_Sample output:_

~~~ console
$ kubectl apply -f ftp-server
deployment.apps/ftp-server created
~~~

## Step 8: Expose the FTP server

In the east namespace, use `skupper expose` to expose the FTP
server on all linked sites.

_**Console for east:**_

~~~ shell
skupper expose deployment/ftp-server --port 21100 --port 2121 --target-port 21100 --target-port 21
~~~

_Sample output:_

~~~ console
$ skupper expose deployment/ftp-server --port 21100 --port 2121 --target-port 21100 --target-port 21
deployment ftp-server exposed as ftp-server
~~~

## Step 9: Run the FTP client

In the west namespace, use `kubectl run` and the `curl` image to
perform FTP put and get operations.

_**Console for west:**_

~~~ shell
kubectl run ftp-put --attach --rm --image=docker.io/curlimages/curl --restart=Never -- -T /etc/os-release ftp://example:example@ftp-server:2121
kubectl run ftp-get --attach --rm --image=docker.io/curlimages/curl --restart=Never -- ftp://example:example@ftp-server:2121/os-release
~~~

_Sample output:_

~~~ console
$ kubectl run ftp-put --attach --rm --image=docker.io/curlimages/curl --restart=Never -- -T /etc/os-release ftp://example:example@ftp-server:2121
pod "ftp-put" deleted

$ kubectl run ftp-get --attach --rm --image=docker.io/curlimages/curl --restart=Never -- ftp://example:example@ftp-server:2121/os-release
NAME="Alpine Linux"
ID=alpine
VERSION_ID=3.15.4
PRETTY_NAME="Alpine Linux v3.15"
HOME_URL="https://alpinelinux.org/"
BUG_REPORT_URL="https://bugs.alpinelinux.org/"
pod "ftp-get" deleted
~~~

## Accessing the web console

Skupper includes a web console you can use to view the application
network.  To access it, use `skupper status` to look up the URL of
the web console.  Then use `kubectl get
secret/skupper-console-users` to look up the console admin
password.

**Note:** The `<console-url>` and `<password>` fields in the
following commands are placeholders.  The actual values are
specific to your environment.

_**Console for west:**_

~~~ shell
skupper status
kubectl get secret/skupper-console-users -o jsonpath={.data.admin} | base64 -d
~~~

_Sample output:_

~~~ console
$ skupper status
Skupper is enabled for namespace "west" in interior mode. It is connected to 1 other site. It has 1 exposed service.
The site console url is: <console-url>
The credentials for internal console-auth mode are held in secret: 'skupper-console-users'

$ kubectl get secret/skupper-console-users -o jsonpath={.data.admin} | base64 -d
<password>
~~~

Navigate to `<console-url>` in your browser.  When prompted, log
in as user `admin` and enter the password.

## Cleaning up

To remove Skupper and the other resources from this exercise, use
the following commands.

_**Console for west:**_

~~~ shell
skupper delete
~~~

_**Console for east:**_

~~~ shell
skupper delete
kubectl delete deployment/ftp-server
~~~

## Next steps


Check out the other [examples][examples] on the Skupper website.
