# Auto Node Repair

- Kubernetes controller to repair nodes which are `NotReady` by replacing them with new fresh nodes. This is done by manipulating `AutoScalingGroups` (AWS) and `VirtualMachineScaleSets` and `AvailabilitySets` (Azure) to repair the nodes.
- Currently supports AWS and Azure (VMSS/VMAS/AKS/ACS) cloud provider.

## How does it work?

- Control loop for each Auto Scaling Group or Virtual Machine Scale Set (VMSS) or Availability Sets (VMAS) configured for a shoot cluster :
  _ Identify `Nodes` which are `NotReady` since configurable amount of time (~10 minutes).
  _ Create new nodes and wait until they are `Ready`
  _ Cordon and drain all `NotReady` nodes.
  _ Delete the `NotReady` nodes.
- Apply this approach for each ASG or VMSS/VMAS in a shoot cluster one by one.
- For a given ASG or VMSS/VMAS, create excess `Nodes` in parallel but cordon, drain and delete `Nodes` one by one.
- If ASG or VMSS/VMAS does not have sufficient capacity for excess `Nodes`, first delete the `NotReady` nodes then create new one.

## Make commands

| Command        | Implication                  |
| -------------- | ---------------------------- |
| `make compile` | Build the go code locally    |
| `make release` | Deploy image into Docker Hub |

## Steps for development

Use the `deploy/kubernetes/deployment-aws.yaml` (AWS) or [Azure examples](deploy/kubernetes/azure) to deploy the auto-node-repair into the cluster. Refer to this file for more details.
