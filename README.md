# Auto Node Repair
* Kubernetes controller to repair nodes which are `NotReady` by replacing them with new fresh nodes. This is done by manipulating `AutoScalingGroups` to repair the nodes.
* Currently supports only AWS cloud provider.

## Archived component no longer maintained
* This component is not used by [gardener](https://github.com/gardener) anymore and no longer maintained. It was archived in the [gardener-attic](https://github.com/gardener-attic).

### How does it work ?
* Control loop for each Auto Scaling Group configured for a shoot cluster :
	* Identify `Nodes` which are `NotReady` since configurable amount of time (~10 minutes).
	* Create new nodes and wait until they are `Ready`
	* Cordon and drain all `NotReady` nodes.
	* Delete the `NotReady` nodes.
* Apply this approach for each ASG in a shoot cluster one by one.
* For a given ASG, create excess `Nodes` in parallel but cordon, drain and delete `Nodes` one by one.
* If ASG does not have sufficient capacity for excess `Nodes`, first delete the `NotReady` nodes then create new one.

### Make commands

| Command       | Implication                              |
| ------------- |------------------------------------------|
| Make compile  | Build the go code locally                |
| Make release  | Deploy image into Gcloud 				   |

### Steps for development

Use the `deploy/kubernetes/deployment.yaml` to deploy the auto-node-repair into the cluster. Refer to this file for more details.
