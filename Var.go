package testcase

// TODO: update Ts to [T] when Go2 released

// Var is a testCase helper structure, that allows easy way to access testCase runtime variables.
// In the future it will be updated to use Go2 type parameters.
//
// Var allows creating testCase variables in a modular way.
// By modular, imagine that you can have commonly used values initialized and then access it from the testCase runtime spec.
// This approach allows an easy dependency injection maintenance at project level for your testing suite.
// It also allows you to have parallel testCase execution where you don't expect side effect from your subject.
//   e.g.: HTTP JSON API testCase and GraphQL testCase both use the business rule instances.
//   Or multiple business rules use the same storage dependency.
//
// The last use-case it allows is to define dependencies for your testCase subject before actually assigning values to it.
// Then you can focus on building up the testing spec and assign values to the variables at the right testing subcontext. With variables, it is easy to forget to assign a value to a variable or forgot to clean up the value of the previous run and then scratch the head during debugging.
// If you forgot to set a value to the variable in testcase, it warns you that this value is not yet defined to the current testing scope.
type Var struct /* [T] */ {
	// Name is the testCase spec variable group from where the cached value can be accessed later on.
	// Name is Mandatory when you create a variable, else the empty string will be used as the variable group.
	Name string
	// Init is an optional constructor definition that will be used when Var is bonded to a *Spec without constructor function passed to the Let function.
	// The goal of this field to initialize a variable that can be reused across different testing suites by bounding the Var to a given testing suite.
	//
	// Please use #Get if you wish to access a testCase runtime across cached variable value.
	// The value returned by this is not subject to any #Before and #Around hook that might mutate the variable value during the testCase runtime.
	// Init function doesn't cache the value in the testCase runtime spec but literally just meant to initialize a value for the Var in a given testCase case.
	// Please use it with caution.
	Init letBlock /* [T] */
}

// Get returns the current cached value of the given Variable
// When Go2 released, it will replace type casting
func (v Var) Get(t *T) (T interface{}) {
	if !t.vars.knows(v.Name) && v.Init != nil {
		t.vars.let(v.Name, v.Init)
	}

	return t.I(v.Name).(interface{}) // cast to T
}

// Set sets a value to a given variable during testCase runtime
func (v Var) Set(t *T, value interface{}) {
	t.Let(v.Name, value)
}

// Let allow you to set the variable value to a given spec
func (v Var) Let(s *Spec, blk letBlock) Var {
	if blk == nil && v.Init != nil {
		return s.Let(v.Name, v.Init)
	}
	return s.Let(v.Name, blk)
}

// LetValue set the value of the variable to a given block
func (v Var) LetValue(s *Spec, value interface{}) Var {
	return s.LetValue(v.Name, value)
}

// EagerLoading allows the variable to be loaded before the action and assertion block is reached.
// This can be useful when you want to have a variable that cause side effect on your system.
// Like it should be present in some sort of attached resource/storage.
//
// For example you may persist the value in a storage as part of the initialization block,
// and then when the testCase/then block is reached, the entity is already present in the resource.
func (v Var) EagerLoading(s *Spec) Var {
	s.Before(func(t *T) { _ = v.Get(t) })
	return v
}
