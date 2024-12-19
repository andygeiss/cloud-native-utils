<p align="center">
<img src="https://github.com/andygeiss/cloud-native-utils/blob/main/logo.png?raw=true" />
</p>

# Cloud Native Utils

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a4850cd982334d38b743a2e1f54ac62a)](https://app.codacy.com/gh/andygeiss/cloud-native-utils?utm_source=github.com&utm_medium=referral&utm_content=andygeiss/cloud-native-utils&utm_campaign=Badge_Grade)
[![License](https://img.shields.io/github/license/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/cloud-native-utils)](https://goreportcard.com/report/github.com/andygeiss/cloud-native-utils)

A collection of high-performance, modular utilities for enhancing testing,
transactional consistency, efficiency, security and stability in cloud-native
Go applications.

## **Module Features**

- [**`assert`**](assert/): Provides tools for testing, including utility functions
to assert value equality and simplify debugging during development.
- [**`consistency`**](consistency/): Implements transactional log management with
`Event` and `EventType` abstractions, and supports file-based persistence using
`JsonFileLogger` and `GobFileLogger` for reliable data storage.
- [**`efficiency`**](efficiency/): Offers utilities for generating read-only
channels, merging and splitting streams, concurrent processing of channel items,
and partitioning key-value stores using shards for scalability and performance.
- [**`security`**](security/): Includes encryption and decryption with AES-GCM,
secure key generation, HMAC hashing, bcrypt-based password handling, and a
preconfigured secure HTTP server with TLS for robust application security.
- [**`service`**](service/): Enhances service orchestration by grouping related
functionality and wrapping functions to support context-aware execution in
cloud-native environments.
- [**`stability`**](stability/): Ensures service robustness with mechanisms like
circuit breakers, retries for transient failures, throttling for rate limiting,
debounce for execution control, and timeouts for enforcing execution limits.
