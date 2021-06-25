# Changelog

## 0.9.0

- Changed the data structure. This makes using old version of the database incompatible with this version.
- Events now get assigned an incrementing integer along with their ID.
- Add requests for events which specify an ID which is already in use, are ignored. This makes idempotent addition possible.
- Added Prometheus endpoint at /metrics which is reachable at port 17654 by default.
- Reintroduced a simplified HTTP API which takes json commands, according to the proto file.
- Removed memory mode

## 0.8.0

### **New**

- Added gRPC server instead of HTTP and RESP server, because it's faster, allows for easier client generation, allows easier routing and has a self documenting API. The HTTP server will be replaced with grpcweb to allow browsers to access the database, for dashboards.
- Added flowctl command line interface to easily communicate with a cluster from the terminal.
- Lots of refactoring
- Bugfixes

### **Removed**

- Removed RESP and HTTP API.

## 0.7.0

### **New**

- Replaced BoltDB with BadgerDB for increase performence and additional features. This change makes older version of data completely unusable. Instead of buckets, keys are prefixed by a constant.
- Added TLS support
- Added authentication to the web API and the RESP API

### **Changed**

- Added an options struct to the BadgerEventStore to enable and disable some features.

### **Fixed**

- Fixed how checksums are calculated.
- Fixed the registration of the GetAll command handler.

## 0.6.0-

Not documented
