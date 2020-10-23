# CloudFoundry Basic Auth Route Service

This is a Proof-of-Concept CloudFoundry app that implements a
[route-service](https://docs.cloudfoundry.org/services/route-services.html) to
add HTTP basic authentication to an application.

This uses a single pre-configured username and password. These are configured
by setting the `AUTH_USERNAME` and `AUTH_PASSWORD` environment variables, which
can be set with the `cf set-env` command.

If your CF deployment has a self-signed SSL certificate, set the
`SKIP_SSL_VALIDATION` environment variable to avoid SSL errors when proxying to
the backend.

## Deployment script
You can deploy this route service to GOV.UK PaaS using the deployment in `deploy.sh`. It will not work outside of GOV.UK PaaS, because of some assumptions made about domain names.

To run the script, see the example below:
```shell script
AUTH_USERNAME=username \
AUTH_PASSWORD=password \
ROUTE_SERVICE_APP_NAME=name_of_this_app \
ROUTE_SERVICE_NAME=name_to_give_the_route_service \
PROTECTED_APP_NAME=name_of_the_app_to_protect_with_basic_auth \
./deploy.sh
```
