<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Testing Doubles](#testing-doubles)
  - [Dummy](#dummy)
    - [PRO](#pro)
    - [CON](#con)
    - [Example](#example)
    - [Use](#use)
  - [Fake](#fake)
    - [PRO](#pro-1)
    - [CON](#con-1)
    - [Example](#example-1)
    - [Use](#use-1)
  - [Stub](#stub)
    - [PRO](#pro-2)
    - [CON](#con-2)
    - [Example](#example-2)
  - [Spy](#spy)
    - [PRO](#pro-3)
    - [CON](#con-3)
    - [Use](#use-2)
  - [Mock](#mock)
    - [PRO](#pro-4)
    - [CON](#con-4)
    - [Use](#use-3)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Testing Doubles

## Dummy

Dummy objects are values passed around but never actually used.
They meant to fill parameter lists.

### PRO

- they can be anything

### CON

- TBD
  
### Example

```go
func Test(t *testing.T) {
    // ...
    logger.Info("foo") // "foo" is a dummy in this case
    // ...        
}
``` 

### Use
- fill parameter lists

## Fake

A Fake is a working implementation,
but usually take some shortcut which makes them not suitable for production.

Fakes are suppliers that have working implementations, but not the same as the production ones.
Usually, they take some shortcuts and have simplified versions of production code.
The proper fake implementation also compliant with the [contract](/docs/contracts.md) a role interface has, like the production one.
 
An example of this shortcut can be an in-memory implementation of a Repository role interface.
This fake implementation will not engage an actual database,
but will use a simple collection to store data.

This approach allows us to do integration-testing of services without starting up a database and performing time-consuming requests.

Apart from testing, fake implementation can come in handy for prototyping and spikes.
We can quickly implement and run our system with an in-memory store,
deferring decisions about what technology and concrete design should be used.

Fakes can simplify local development when working with complex external systems.

example use cases:  
- payment system that always returns with successful payment and does the callback automatically on request.
- email verification process call verify callback instead of sending an email out.
- [in-memory database for testing](https://martinfowler.com/bliki/InMemoryTestDatabase.html)
 
### PRO

- you can start developing your business rules without the need to choose a technology stack ahead of time before you know your business requirements.
- can support easier local development both in local manual testing and in integration tests.
- allows testing suite optimizations when using real components drastically increases the testing feedback loop time.
- allows taking shortcuts instead of using the concrete external resources when the application runs locally for development purposes.

### CON

- using fake without an [role interface contract](/docs/contracts.md) introduce manual maintenance costs. 
- neglecting to keep fake in sync with production variant will risk violating dev/prod parity in the project's testing suite.

### Example

[Example fake implementation](/docs/testing-double/fake_test.go) for the [example role interface + contract](/docs/testing-double/spec_helper_test.go).

### Use

- test happy path with it
- replace real implementation in tests when the feedback loop with it is too slow
- test business logic with it

## Stub

Stub provide canned answers to calls made during the test,
usually not responding at all to anything outside what's programmed in for the test.

Method Stubbing within the stub allows you to manipulate one or two methods to inject mostly errors to test rainy paths with it.

My suggestion is only to stub a method to fault inject,
and avoid representing a happy path with it whenever possible.   

### PRO

- relatively easy to use
- ideal to inject error with it

### CON

- when stub testing double used for representing a happy path, we need to introduce a manual chore activity
  to the project to ensure the stub content up to date with the production

### Example

- [Stub Object](/docs/testing-double/stub_test.go)
- [Stub One method on a real Object for fault injection](/docs/testing-double/stub_method_test.go) 

## Spy
Spy are stubs that also record some information based on how they were called.

Often used to verify "indirect output" of the tested code,
by asserting the expectations afterward,
without having defined the expectations before the tested code is executed.
It helps in recording information about the indirect object created.

### PRO

- everything true to stub
- can help to debug

### CON

- everything that true to stub
- risk that test will focus on implementation details if misused. 

### Use

- checking retry logic behavioral requirements from an analytical point of view

## Mock

Mock are pre-programmed with expectations which form a specification of the calls they are expected to receive.
They can throw an exception if they receive a call they don't expect 
and are checked during verification to ensure they got all the calls they were expecting.

Mocks shine the most when used for large teams where different parts of the code developed in parallel.
After an initial agreement between members about the interface and high-level behavior between components,
they can start to develop without anything concrete.
This approach cost architecture flexibility by introducing tech debt in the testing suite,
but this debt can be fixed later by cleaning up tests
when all the components are already integrated into the system. 

Another example of using mocks when the project has too much entropy,
and testing with real components would require way too much extra effort.
The extra effort would not stop with only fixing the code,
but most likely involve additional rounds of knowledge sharing as well in the team.
When this is not possible, it is "acceptable" to test implementation details with mocks,
while making sure the behavior requirements of the role interface that being mocked are kept at a minimum.
Ideally, try to avoid using multiple mocks in tests whenever is possible.

### PRO

- allows defining implementation details expectations towards a dependency
- flexible to do many things
- almost every people familiar with using this testing double   

### CON

- same which is true to stub and spy
- introduce technical debt in the project's testing suite
- have a high risk to misuse it, and make your test focus on implementation details.
- if avoiding mocks is not an option, that's possibly feedback about the project software design state.   

### Use

- develop components in parallel with large or distributed teams.
- fault injection
- allows isolated unit testing in projects with high entropy level
