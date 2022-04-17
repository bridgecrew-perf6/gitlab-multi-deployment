# Gitlab Multi Deployment
## The idea
Deploy a service including it's dependencies (e.g. services that dependent on an updated API). All services will be deployed using kubectl and tagged via the Gitlab API.
Works out of the box for Gitlab based repositories, other repository providers need adaptions.

## How to use?
Configure a prod deployment using the `deployment.json` file. A deployment is only executed when all services and their defined pipelines are found, otherwise nothing is deployed.

### Example
```
[
    {
        "repository": "service-a",
        "pipeline": 265929
    },
    {
        "repository": "service-b",
        "pipeline": 27543
    },
    {
        "repository": "service-c",
        "pipeline": 34959
    }
]
```

To run, trigger the second manual pipeline step of the latest commit.