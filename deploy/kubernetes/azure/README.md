# Auto Node Repair on Azure

The auto node repair on Azure scales worker nodes with Virtual Machine Scale Sets, Availability Sets, Azure Kubernetes Service (AKS) and Azure Container Service (ACS). It will run as a Kubernetes deployment in your cluster. This README will go over some of the necessary steps required to get the auto node repair up and running.

## Kubernetes Version

Auto node repair for Azure supports four VM types with Azure cloud provider:

- **vmss**: For kubernetes cluster running on VMSS instances. Azure cloud provider's `vmType` parameter must be configured as 'vmss'. It requires Kubernetes with Azure VMSS support (supported for [acs-engine](https://github.com/Azure/acs-engine) VMSS clusters)
- **standard**: For kubernetes cluster running on VMAS instances. Azure cloud provider's `vmType` parameter must be configured as 'standard' or left as empty string. It only supports Kubernetes cluster deployed via [acs-engine](https://github.com/Azure/acs-engine).
- **aks**: Managed Container Service([AKS](https://docs.microsoft.com/en-us/azure/aks/))
- **acs**: Container service([ACS](https://docs.microsoft.com/en-us/azure/container-service/kubernetes/))

Only **vmss** vmType supports scaling to zero nodes.

## Permissions

Get azure credentials by running the following command

```sh
# replace <subscription-id> with yours.
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/<subscription-id>" --output json
```

## Deployment manifests

### VMSS deployment

Pre-requirements:

- Get credentials from above `permissions` step.
- Get the scale set name which is used for node auto repair.
- Encode each data with base64.

Fill the values of auto-node-repair-azure secret in [auto-node-repair-vmss.yaml](auto-node-repair-vmss.yaml), including

- ClientID: `<base64-encoded-client-id>`
- ClientSecret: `<base64-encoded-client-secret>`
- ResourceGroup: `<base64-encoded-resource-group>`
- SubscriptionID: `<base64-encode-subscription-id>`
- TenantID: `<base64-encoded-tenant-id>`

> Note that all data above should be encoded with base64.

And fill the node groups in container command by `--nodes`, e.g.

```yaml
        - --nodes=1:10:vmss1
```

or multiple node groups:

```yaml
        - --nodes=1:10:vmss1
        - --nodes=1:10:vmss2
```

Then deploy auto-node-repair by running

```sh
kubectl create -f auto-node-repair-vmss.yaml
```

To run a CA pod in master node - CA deployment should tolerate the master `taint` and `nodeSelector` should be used to schedule the pods in master node.

```sh
kubectl create -f auto-node-repair-vmss-master.yaml
```

To run a CA pod with Azure managed service identity (MSI), use [auto-node-repair-vmss-msi.yaml](auto-node-repair-vmss-msi.yaml) instead:

```sh
kubectl create -f auto-node-repair-vmss-msi.yaml
```

### Standard deployment

Pre-requirements:

- Get credentials from above `permissions` step.
- Get the initial Azure deployment name from azure portal. If you have multiple deployments (e.g. have run `acs-engine scale` command), make sure to get the first one.
- Get a node pool name for nodes scaling from acs-engine deployment manifests
- Encode each data with base64.

Fill the values of auto-node-repair-azure secret in [auto-node-repair-standard-master.yaml](auto-node-repair-standard-master.yaml), including

- ClientID: `<base64-encoded-client-id>`
- ClientSecret: `<base64-encoded-client-secret>`
- ResourceGroup: `<base64-encoded-resource-group>`
- SubscriptionID: `<base64-encode-subscription-id>`
- TenantID: `<base64-encoded-tenant-id>`
- Deployment: `<base64-encoded-azure-initial-deploy-name>`

> Note that all data above should be encoded with base64.

And fill the node groups in container command by `--nodes`, e.g.

```yaml
        - --nodes=1:10:agentpool1
```

or multiple node groups:

```yaml
        - --nodes=1:10:agentpool1
        - --nodes=1:10:agentpool2
```

Create Azure deploy parameters secret `auto-node-repair-azure-deploy-parameters` by running

```sh
kubectl -n kube-system create secret generic auto-node-repair-azure-deploy-parameters --from-file=deploy-parameters=./_output/<your-output-path>/azuredeploy.parameters.json
```

Then deploy auto-node-repair by running

```sh
kubectl create -f auto-node-repair-standard-master.yaml
```

To run a CA pod with Azure managed service identity (MSI), use [auto-node-repair-standard-msi.yaml](auto-node-repair-standard-msi.yaml) instead:

```sh
kubectl create -f auto-node-repair-standard-msi.yaml
```

**WARNING**: Auto node repair depends on user provided deployment parameters to provision new nodes. It should be redeployed with new parameters after upgrading Kubernetes cluster (e.g. upgraded by `acs-engine upgrade` command), or else new nodes will be provisioned with old version.

### AKS or ACS deployment

Pre-requirements:

- Get credentials from above `permissions` step.
- Get the cluster name using the following:

  for AKS:

  ```sh
  az aks list
  ```

  for ACS:

  ```sh
  az acs list
  ```

- Get a node pool name by extracting the value of the label **agentpool**

  ```sh
  kubectl get nodes --show-labels
  ```

- In case of AKS we need additional information in the form of node resource group.
  Use the value of the label by name **kubernetes.azure.com/cluster** as the node resource group.

- Encode each data with base64.

Fill the values of auto-node-repair-azure secret in [auto-node-repair-containerservice](auto-node-repair-containerservice.yaml), including

- ClientID: `<base64-encoded-client-id>`
- ClientSecret: `<base64-encoded-client-secret>`
- ResourceGroup: `<base64-encoded-resource-group>` (Note: Please use lower case)
- SubscriptionID: `<base64-encode-subscription-id>`
- TenantID: `<base64-encoded-tenant-id>`
- ClusterName: `<base64-encoded-clustername>`
- NodeResourceGroup: `<base64-encoded-node-resource-group>` (Note: AKS only parameter. Please use the value of kubernetes.azure.com/cluster label verbatim (case sensitive))

> Note that all data above should be encoded with base64.


And fill the node groups in container command by `--nodes`, with the range of nodes (minimum to be set as 3 which is the default cluster size) and node pool name obtained from pre-requirements steps above, e.g.

```yaml
        - --nodes=3:10:nodepool1
```

The vmType param determines the kind of service we are interacting with.

For AKS fill the following base64 encoded value:

```sh
$ echo AKS | base64
QUtTCg==
```

and for ACS fill the following base64 encoded value:

```sh
$echo ACS | base64
QUNTCg==
```

Then deploy auto-node-repair by running

```sh
kubectl create -f auto-node-repair-containerservice.yaml
```