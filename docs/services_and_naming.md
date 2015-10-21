# Services and Naming

There are two ways that a dependency may manifest itself. First, Package
A may require Package B. This case is accomodated with the
`dependencies` directive in the `Chart.yaml`.

But a possibly more common scenario would be a service-level dependency:
Package A requires _a Postgres service_. It doesn't matter to Package A
whether that's an HA clustered Postgres managed by Governor, or just a
stand-along Postgres server.

Kubernetes provides a huge part of the solution: A `service` can be
declared independent of the actual containers that implement that
service.

So one possible approach would be to provide `service` packages that
declare themselves into the environment, and use label queries to select
a matching backend (or backends).

Following this approach requires a shared vocabulary among package
contributors, but that is a problem that can be addressed (if not
solved) by a combination of good tooling and documentation.

## Example: Memcached Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: memcached-provider
  labels:
    - name: memached-provider
    - from: helm
    - heritage: deis
spec:
  ports:
    - name: client
      port: 11211
      protocol: TCP
  selector:
    - fulfills: memcached-provider
```

Service discovery works like this:

- Anything that needs Memcached support can look for
  `MEMCACHED_PROVIDER_SERVICE_HOST` and
  `MEMCACHED_PROVIER_SERVIER_PORT_CLIENT`.
- Anything that _implements_ a memcached service needs to simply add a
  label to its `metadata: labels` section: `fulfills:
  memcached-provider`.
