# maybe
This package provides a sum type inspired by Haskell's `Maybe` monad.

Instances can by associated with a value or be empty:
* The `Unit` function is used to create an instance with a value.
* The `Nothing` function is used to create an instance without an associated value.

The `Maybe` type meshes well with JSON (un)marshalling and with SQL drivers.
