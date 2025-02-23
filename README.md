An exploratory project while learning Go.

Nothing special here, but I might use this as a handy no-fuss file system utility in the future.

One hurdle so far is cases where I'd use default arguments in other languages.

I've not yet decided which tradeoff to choose between required, but non-interesting
arguments, some kind of configuration struct, a builder pattern, or just a bunch of
named functions for combinations of arguments.

I'm leaning toward a combination of a configuration struct and a simple non-configurable
version of the function for the "default" case, but some convention exploration is needed, because I haven't read enough Go.

Until then, we will have no append option in our WriteFile functions :^)