# Changelog

## 0.7.0

### **New**

- Replaced BoltDB with BadgerDB for increase performence and additional features. This change makes older version of data completely unusable. Instead of buckets, keys are prefixed by a constant.

### **Changed**

- Added an options struct to the BadgerEventStore to enable and disable some features.

### **Fixed**

- Fixed how checksums are calculated.
- Fixed the registration of the GetAll command handler.

## 0.6.0-

Not documented
