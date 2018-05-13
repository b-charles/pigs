package ioc

import (
	"fmt"
	"reflect"
)

// Component registration

type ComponentRegistrationError struct {
	name   string
	source interface{}
}

func (self *ComponentRegistrationError) Error() string {
	return fmt.Sprintf("Error during registration of '%s': %v", self.name, self.source)
}

// Component registration - not unique main name

type NotUniqueMainNameError struct {
	name string
}

func (self *NotUniqueMainNameError) Error() string {
	return fmt.Sprintf("A component with a main name '%s' is already registered.", self.name)
}

// Component registration - not unique alias

type NotUniqueAliasError struct {
	alias string
}

func (self *NotUniqueAliasError) Error() string {
	return fmt.Sprintf("Alias specified more than once: '%s'.", self.alias)
}

// Component registration - the factory is not a function

type NotFunctionFactoryError struct {
	kind reflect.Kind
}

func (self *NotFunctionFactoryError) Error() string {
	return fmt.Sprintf("The factory is not a %v, and not a function.", self.kind)
}

// Component registration - invalid input number

type InvalidInputNumberError struct {
	actual, defined int
}

func (self *InvalidInputNumberError) Error() string {
	return fmt.Sprintf("The actual input number (%d) doen't match the definition (%d).", self.actual, self.defined)
}

// Component registration - invalid output number

type InvalidOutputNumberError struct {
	actual int
}

func (self *InvalidOutputNumberError) Error() string {
	return fmt.Sprintf("The output number should be 1, not %d.", self.actual)
}

// Component registration - invalid output kind

type InvalidOutputKindError struct {
	actual reflect.Kind
}

func (self *InvalidOutputKindError) Error() string {
	return fmt.Sprintf("The output type should be a chan, a function, an interface, a map, a pointer or a slice, not a %v.", self.actual)
}

// Component instanciation

type ComponentInstanciationError struct {
	component *Component
	source    interface{}
}

func (self *ComponentInstanciationError) Error() string {
	return fmt.Sprintf("Error during instanciation of '%s': %v", self.component.name, self.source)
}

// Component instanciation - invalid key type of map inject

type InvalidKeyTypeError struct {
	keyType reflect.Type
}

func (self *InvalidKeyTypeError) Error() string {
	return fmt.Sprintf("Unsupported key type for a map injection: %v, only 'string' is valid.", self.keyType)
}

// Component instanciation - no producer found

type NoProducerError struct {
	name string
}

func (self *NoProducerError) Error() string {
	return fmt.Sprintf("No producer found for '%s'.", self.name)
}

// Component instanciation - too many producers found

type TooManyProducersError struct {
	name      string
	producers []*Component
}

func (self *TooManyProducersError) Error() string {
	return fmt.Sprintf("Too many producers found for '%s': %v.", self.name, self.producers)
}

// Component initialisation

type ComponentInitialisationError struct {
	component *Component
	source    interface{}
}

func (self *ComponentInitialisationError) Error() string {
	return fmt.Sprintf("Error during initialisation of '%s': %v", self.component.name, self.source)
}

// Instanciate components

type ComponentInstanciationsError struct {
	name   string
	source interface{}
}

func (self *ComponentInstanciationsError) Error() string {
	return fmt.Sprintf("Error during instanciations of '%s': %v", self.name, self.source)
}
