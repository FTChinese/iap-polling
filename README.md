# IAP Polling Service

## Build

This project consists of two apps: `iap-polling` and `migrate-receipt`. The build processes only differs in app name. By default it builds the `iap-polling` app. You can set make's `APP` variable to `migrate` to build the `migrate-receipt` app.

Build `iap-polling:

```
make
```

Build `migrate-receipt`:

```
make APP=migrate
```

## Run

Both apps have a command line argument `-production` used in production environment.
