# Modeling Services in Kubernetes

Generally speaking there are two ways to model [Services](http://kubernetes.io/v1.0/docs/user-guide/services.html) in Kubernetes:

1. Service providers (i.e. a database) define their own services
2. Service consumers (i.e. an application) define the services they consume

With the first option, service consumers must constantly track the service providers they depend on.  When a service provider changes, all consumers must be notified and restarted.  This results in tight coupling between service consumers and providers.

With the second option, consumers defines the services they require with label selectors.  Once pods are launched that fulfill the label selectors, the consumer can begin accessing the service.  This facilitates looser coupling between service consumers and providers.

While the first option is more common in traditional orchestration systems, the second is a more natural fit for Kubernetes.  As a result, #2 is how we model services in Helm Classic.  See the documentation on [using labels](using_labels.md) for more details.

## Example: Ponycorn, a Redis-backed Application

We are authoring a chart for an application called `ponycorn` which relies on a Redis backend.  In addition to creating manifests for the replication controllers and pods required by the application, we also include manifests to define the services we consume.

To model the Redis backend in the Ponycorn chart, we include:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: ponycorn-redis
  labels:
    heritage: helm
spec:
  ports:
    - port: 6379
      name: redis
      protocol: TCP
  selector:
    provider: redis
```

For service discovery, the `ponycorn` application uses the `PONYCORN_REDIS_SERVICE_HOST` and `PONYCORN_REDIS_SERVICE_PORT` environment variables.  These environment variable names remain constant and point to a static Kubernetes Service VIP.  As pods are launched into the cluster that _fulfill_ the `provider: redis` label selector, the Redis service for `ponycorn` comes online.

Following this approach requires a shared vocabulary among package contributors, but that is a problem that can be addressed (if not solved) by a combination of good tooling and documentation.
