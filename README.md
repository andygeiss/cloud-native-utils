<p align="center">
<img src="https://github.com/andygeiss/cloud-native-utils/blob/main/logo.png?raw=true" />
</p>

# Cloud Native Utils

[![License](https://img.shields.io/github/license/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/cloud-native-utils)](https://goreportcard.com/report/github.com/andygeiss/cloud-native-utils)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b4e3a9c4859b47f1bc43613970ec8d12)](https://app.codacy.com/gh/andygeiss/cloud-native-utils/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/b4e3a9c4859b47f1bc43613970ec8d12)](https://app.codacy.com/gh/andygeiss/cloud-native-utils/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

A collection of high-performance, modular utilities for enhancing testing,
transactional consistency, efficiency, security, and stability in cloud-native
Go applications.

## **Module Features**

- [**`assert`**](assert/): Provides tools for testing, including utility functions
  to assert value equality and simplify debugging during development.
- [**`consistency`**](consistency/): Implements transactional log management with
  `Event` and `EventType` abstractions, and supports file-based persistence using
  `JsonFileLogger` for reliable data storage.
- [**`efficiency`**](efficiency/): Offers utilities for generating read-only
  channels, merging and splitting streams, concurrent processing of channel items,
  and partitioning key-value stores using shards for scalability and performance.
- [**`extensibility`**](extensibility/): Dynamically load external Go plugins using
  `LoadPlugin`. Just provide a symbol name (e.g., a function) to integrate new
  features on-the-fly—no rebuilds or redeploys required.
- [**`security`**](security/): Includes encryption and decryption with AES-GCM,
  secure key generation, HMAC hashing, bcrypt-based password handling, and a
  preconfigured secure HTTP(S) server with liveness and readiness probes for
  robust application security.
- [**`service`**](service/): Enhances service orchestration by grouping related
  functionality, wrapping functions to support context-aware execution and add
  lifecycle-oriented functionality like signal handling in cloud-native
  environments.
- [**`stability`**](stability/): Ensures service robustness with mechanisms like
  circuit breakers, retries for transient failures, throttling for rate limiting,
  debounce for execution control, and timeouts for enforcing execution limits.
- [**`templating`**](templating/): Provides an `Engine` for managing templates
  stored in an embedded filesystem. Use `Parse` to load multiple templates (via
  glob patterns), and `Render` to execute them with custom data.
