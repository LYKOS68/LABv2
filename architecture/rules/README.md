---
schema: rap-lab.architecture-rules
schema_version: 1.0.0
status: proposed
versioning: semver
last_updated: 2026-10-25
scope: architecture-and-workflow-governance
normative_terms:
  MUST: required for conformance
  MUST_NOT: prohibited for conformance
  SHOULD: recommended unless an explicit exception is recorded
---

# RAP-LAB Architecture Rules

These rules define testable architecture constraints for RAP-LAB teams and software. They are pattern-level requirements, not claims that the current repository already conforms.

## Rule record contract

Every rule MUST define exactly these fields:

- `ID`: stable identifier; do not reuse retired identifiers.
- `TRIGGER`: condition that activates the rule.
- `PRECONDITIONS`: facts that must hold before execution or review.
- `ALLOWED_OPERATIONS`: actions permitted under this rule.
- `FORBIDDEN_OPERATIONS`: actions that violate this rule.
- `POSTCONDITIONS`: observable state that must hold after completion.
- `TESTS`: reproducible checks with an objective pass/fail result.

A rule conforms only when its postconditions are demonstrably satisfied and its tests pass. Unknown evidence is not a pass. Exceptions MUST identify the rule, scope, reason, owner, expiry/review date, and compensating controls.

## Versioning policy

- **MAJOR**: incompatible change to a rule's required behavior, field contract, or conformance meaning.
- **MINOR**: additive rule or backward-compatible requirement.
- **PATCH**: clarification that does not change conformance meaning.
- Rule IDs remain stable even when the specification version changes.
- A proposal is not an approved baseline until its PR is reviewed and merged.

## P01 — Composition Root

- **ID:** `RAP-ARCH-P01`
- **TRIGGER:** An application, service, CLI, worker, or team workflow is initialized.
- **PRECONDITIONS:**
  - Required components and external dependencies are identified.
  - Startup order and mandatory prerequisites are documented.
- **ALLOWED_OPERATIONS:**
  - Construct top-level components at the composition boundary.
  - Wire concrete implementations to declared interfaces.
  - Validate mandatory configuration and prerequisites.
  - Delegate domain work to the appropriate component.
- **FORBIDDEN_OPERATIONS:**
  - Hide application-wide dependency construction inside unrelated domain functions.
  - Start workers before required dependencies are ready.
  - Continue startup after a mandatory prerequisite fails.
- **POSTCONDITIONS:**
  - The startup path identifies who constructs each top-level dependency.
  - Startup either reaches a documented ready state or returns a classified error and releases resources already acquired.
- **TESTS:**
  - [ ] Startup dependency order is represented in code or a reviewed dependency map.
  - [ ] A test injects a failing mandatory prerequisite and verifies that readiness is not reported.
  - [ ] A startup failure test verifies cleanup of resources initialized earlier.

## P02 — Explicit Dependency

- **ID:** `RAP-ARCH-P02`
- **TRIGGER:** A component needs another component, service, repository, capability, or external system.
- **PRECONDITIONS:**
  - The consumer's required operations are defined.
  - The dependency has a documented interface or contract.
- **ALLOWED_OPERATIONS:**
  - Pass required dependencies through constructors, explicit parameters, or a declared dependency scope.
  - Use narrow interfaces exposing only required methods.
  - Replace dependencies with test doubles at a documented seam.
- **FORBIDDEN_OPERATIONS:**
  - Introduce undocumented global mutable dependencies.
  - Reach into another component's private state to bypass its contract.
  - Use a dependency merely because it is globally reachable when it is not required.
- **POSTCONDITIONS:**
  - A reviewer can identify the consumer's dependencies from its construction or public contract.
  - Tests can substitute required external dependencies without changing domain logic.
- **TESTS:**
  - [ ] Constructor/function signatures or wiring identify required dependencies.
  - [ ] A unit test substitutes at least one external dependency through the declared seam.
  - [ ] Static review finds no new undocumented global dependency path in the changed component.

## P03 — Capability Boundary

- **ID:** `RAP-ARCH-P03`
- **TRIGGER:** A component or actor receives access to resources, files, records, services, tools, teams, or operations.
- **PRECONDITIONS:**
  - Protected resources and allowed operations are identified.
  - The requesting identity and authorization decision point are known.
- **ALLOWED_OPERATIONS:**
  - Grant only the smallest resource and operation scope required.
  - Check authorization at the enforcement boundary.
  - Record denied or security-relevant access decisions without leaking secrets.
- **FORBIDDEN_OPERATIONS:**
  - Grant broad access by default.
  - Treat possession of an identifier, URL, object reference, or message as authorization by itself.
  - Let an untrusted caller enlarge its own capability.
- **POSTCONDITIONS:**
  - Every protected operation has an explicit authorization decision.
  - Access outside the granted scope is denied.
- **TESTS:**
  - [ ] Positive test proves a permitted operation succeeds.
  - [ ] Negative test proves an out-of-scope operation is denied.
  - [ ] Default configuration grants no broader access than the documented policy.
  - [ ] A reviewer can trace the grant to an explicit source and scope.

## P04 — Message Contract

- **ID:** `RAP-ARCH-P04`
- **TRIGGER:** Work crosses a component, process, team, queue, event, API, or agent boundary.
- **PRECONDITIONS:**
  - Sender and receiver roles are identified.
  - Payload, versioning, failure behavior, and ownership are specified.
- **ALLOWED_OPERATIONS:**
  - Send typed or schema-validated messages through a defined boundary.
  - Validate inputs at the receiving boundary.
  - Return an explicit result, acknowledgment, or classified error where required.
  - Add backward-compatible optional fields under the documented versioning policy.
- **FORBIDDEN_OPERATIONS:**
  - Rely on undocumented field names or implicit shared mutable state as a cross-boundary contract.
  - Silently discard malformed or unsupported messages.
  - Assume delivery, ordering, or exactly-once behavior without a defined guarantee.
- **POSTCONDITIONS:**
  - Sender and receiver agree on a versioned payload contract.
  - Invalid messages produce a documented failure outcome.
- **TESTS:**
  - [ ] Valid payload is accepted by the receiver.
  - [ ] Missing, invalid, or unsupported payload is rejected predictably.
  - [ ] Contract compatibility is checked for the supported version range.
  - [ ] Delivery and retry semantics are documented and tested where applicable.

## P05 — Bounded Dispatch

- **ID:** `RAP-ARCH-P05`
- **TRIGGER:** A dispatcher, scheduler, event loop, worker, or agent processes multiple pending units of work.
- **PRECONDITIONS:**
  - A unit of work is defined.
  - The scheduling limit, fairness expectation, and overload response are stated.
- **ALLOWED_OPERATIONS:**
  - Bound dispatch by count, elapsed time, concurrency, or another measurable budget.
  - Yield when the budget is exhausted.
  - Track latency, backlog, and budget exhaustion when operationally relevant.
- **FORBIDDEN_OPERATIONS:**
  - Process an unbounded backlog in a latency-sensitive path without an explicit exception.
  - Allow one work lane to starve another without a documented policy.
  - Report successful completion while work remains silently stranded.
- **POSTCONDITIONS:**
  - Each dispatch invocation stays within the declared budget or returns an explicit budget-exhausted result.
  - Pending work remains queued or is failed explicitly.
- **TESTS:**
  - [ ] A test provides more work than one budget allows and verifies that one dispatch stays bounded.
  - [ ] Remaining work can be processed by a subsequent dispatch.
  - [ ] A fairness/latency test verifies that a competing priority lane is not indefinitely starved.
  - [ ] Budget exhaustion is observable through a result or metric.

## P06 — Explicit Lifecycle

- **ID:** `RAP-ARCH-P06`
- **TRIGGER:** A component starts asynchronous work, opens resources, registers callbacks, or acquires a lock or subscription with a lifetime.
- **PRECONDITIONS:**
  - Resource owner and shutdown responsibility are assigned.
  - Startup and cancellation failure modes are defined.
- **ALLOWED_OPERATIONS:**
  - Use explicit start, stop/cancel, wait/join, unregister, and release steps as appropriate.
  - Propagate cancellation or deadlines through supported mechanisms.
  - Make cleanup safe when initialization only partially succeeded.
- **FORBIDDEN_OPERATIONS:**
  - Ignore failure to start required asynchronous work.
  - Release resources while workers may still access them.
  - Leave callbacks or subscriptions registered beyond their intended lifetime.
  - Assume cancellation means a worker has already stopped.
- **POSTCONDITIONS:**
  - Every started worker has a defined termination path.
  - Shutdown waits for required workers before releasing resources they use.
  - Partial startup failure releases previously acquired resources.
- **TESTS:**
  - [ ] Injected worker-start failure is detected and surfaced.
  - [ ] Shutdown test verifies workers finish before dependent resources are released.
  - [ ] Partial initialization test verifies cleanup of acquired resources.
  - [ ] Callback/subscription registration is paired with removal.

## P07 — Queue Overflow Policy

- **ID:** `RAP-ARCH-P07`
- **TRIGGER:** A queue, channel, inbox, task buffer, or pending-work collection has a capacity or can receive faster than it is consumed.
- **PRECONDITIONS:**
  - Capacity and expected producer/consumer behavior are defined.
  - Message criticality and permitted loss/coalescing behavior are classified.
- **ALLOWED_OPERATIONS:**
  - Reject, backpressure, shed, coalesce, persist, or prioritize work according to an explicit policy.
  - Count or log relevant overflow outcomes without exposing sensitive payloads.
- **FORBIDDEN_OPERATIONS:**
  - Drop work silently.
  - Block indefinitely without a cancellation or timeout policy.
  - Coalesce messages whose semantics require each individual event unless the contract allows it.
- **POSTCONDITIONS:**
  - At capacity, the system follows a deterministic documented behavior.
  - Overflow or shedding is observable to operators/tests when relevant.
- **TESTS:**
  - [ ] Fill the queue to capacity and submit one additional item.
  - [ ] Verify the documented overflow behavior and resulting state.
  - [ ] Verify an overflow counter, error, or event is emitted where specified.
  - [ ] Verify cancellation/timeout can release a blocked producer where blocking is used.

## P08 — Fail Closed

- **ID:** `RAP-ARCH-P08`
- **TRIGGER:** A mandatory prerequisite, validation, authorization, integrity check, or safety gate fails or cannot produce a trustworthy result.
- **PRECONDITIONS:**
  - Mandatory gates are explicitly distinguished from optional/degraded capabilities.
  - Failure outcomes are defined.
- **ALLOWED_OPERATIONS:**
  - Stop the protected operation and return a classified error.
  - Enter a documented degraded mode only when that mode is explicitly authorized and safe.
  - Emit diagnostic information that excludes secrets.
- **FORBIDDEN_OPERATIONS:**
  - Treat unknown, timed-out, malformed, or unavailable authorization as approval.
  - Continue a protected operation after a mandatory gate fails.
  - Convert a failure into success solely to keep startup or a workflow moving.
- **POSTCONDITIONS:**
  - Failed mandatory gates prevent the protected action.
  - Any degraded mode is explicit, bounded, and visible in status/output.
- **TESTS:**
  - [ ] Force each mandatory gate to fail and verify the protected action does not execute.
  - [ ] Timeout and unknown-result cases do not become implicit approval.
  - [ ] Any supported degraded mode has an explicit activation test and documented limits.

## P09 — Ownership Contract

- **ID:** `RAP-ARCH-P09`
- **TRIGGER:** A component creates, receives, borrows, transfers, or disposes a resource or mutable object.
- **PRECONDITIONS:**
  - Creation point, current owner, allowed borrowers, and release/transfer rule are identifiable.
  - Shared access and mutation rules are documented where applicable.
- **ALLOWED_OPERATIONS:**
  - Define single ownership, explicit transfer, reference/lifetime rules, or synchronized shared ownership.
  - Make cleanup idempotent where repeated cleanup is possible.
  - Return a resource to its owner or transfer ownership through an explicit contract.
- **FORBIDDEN_OPERATIONS:**
  - Release a borrowed resource without permission.
  - Access a resource after its documented lifetime ends.
  - Allow multiple components to believe they have exclusive ownership.
  - Mutate shared state outside its synchronization contract.
- **POSTCONDITIONS:**
  - Every resource has one identifiable release responsibility or a documented shared-ownership policy.
  - Ownership transfer and invalidation are observable in the contract.
- **TESTS:**
  - [ ] Resource lifecycle test verifies release occurs exactly as specified.
  - [ ] Failure-path test verifies no acquired resource is abandoned.
  - [ ] Shared mutable state is covered by synchronization or race checks where supported.
  - [ ] Double cleanup or post-release access is rejected or harmless according to the contract.

## P10 — Least Privilege

- **ID:** `RAP-ARCH-P10`
- **TRIGGER:** A role, service, agent, process, token, team, or integration is assigned permissions.
- **PRECONDITIONS:**
  - The role's required tasks and minimum resource scope are documented.
  - Permission grants can be reviewed or inspected.
- **ALLOWED_OPERATIONS:**
  - Grant only permissions required for the assigned task.
  - Separate read, write, administrative, and delegation rights where feasible.
  - Periodically review and revoke unneeded grants.
- **FORBIDDEN_OPERATIONS:**
  - Use wildcard/all-target access as the default.
  - Reuse a high-privilege identity for routine low-privilege tasks without a documented requirement.
  - Grant permissions because they might be useful later.
- **POSTCONDITIONS:**
  - Each grant has a task-related justification and bounded scope.
  - Unneeded permissions are absent or scheduled for removal under an explicit exception.
- **TESTS:**
  - [ ] Permission inventory maps each grant to a required operation.
  - [ ] Negative tests verify roles cannot perform out-of-scope operations.
  - [ ] Broad grants trigger review and cannot enter the approved baseline without an exception record.
  - [ ] Revocation is verified for removed permissions.

## Conformance report

A reviewer evaluating a component MUST record, per applicable rule:

- Rule ID and specification version.
- Evidence paths, test names, or observable outputs.
- Result: `PASS`, `FAIL`, `NOT_APPLICABLE`, or `UNKNOWN`.
- For `NOT_APPLICABLE`, the rationale.
- For `FAIL` or `UNKNOWN`, the missing control/evidence and owner of the follow-up.

A component MUST NOT be declared conformant while any applicable rule is `FAIL` or `UNKNOWN`.

## Scope note

These rules are intended to be reusable across Go services, CLI agents, web components, queues, integrations, and RAP-LAB team workflows. Implementation-specific mappings (for example, Go interfaces, contexts, channels, mutexes, HTTP handlers, or CI checks) are evidence mechanisms, not interchangeable definitions of the rules themselves.
