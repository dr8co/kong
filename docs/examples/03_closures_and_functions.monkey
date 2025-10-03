// 03_closures_and_functions.monkey
// Demonstrates closures, nested functions, and higher-order functions in Monkey.

// makeAdder returns a function that adds `x` to its argument.
let makeAdder = fn(x) {
  fn(y) { x + y }
}

let addTwo = makeAdder(2)
let addFive = makeAdder(5)

puts(addTwo(3))    // 5
puts(addFive(10))  // 15

// Higher-order functions: passing and returning functions
let apply = fn(f, v) { f(v) }
let double = fn(x) { x * 2 }

puts(apply(double, 7)) // 14

// Returning functions lets us build small DSL-like helpers
let makeMultiplier = fn(factor) {
  fn(x) { x * factor }
}

let triple = makeMultiplier(3)
puts(triple(6)) // 18

// End of example
