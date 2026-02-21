# Lesson 03: Type System Internals

Understanding how TypeScript's type system works under the hood - inference, narrowing, structural typing, and variance.

## Type Inference

TypeScript infers types automatically. Understanding when and how helps write better code.

### Basic Inference

```typescript
let x = 42;           // Type: number (inferred)
let y = 'hello';      // Type: string
let z = [1, 2, 3];    // Type: number[]
let w = { a: 1 };     // Type: { a: number }

// Widening
let a = null;         // Type: any (!)
let b = undefined;    // Type: any (!)
```

### Best Common Type

```typescript
let arr = [1, 2, 'three'];  // Type: (number | string)[]

let mixed = [1, null];      // Type: (number | null)[]

class Animal {}
class Dog extends Animal {}
class Cat extends Animal {}
let animals = [new Dog(), new Cat()];  // Type: (Dog | Cat)[]
// Not Animal[] - TypeScript uses concrete types
```

### Contextual Typing

```typescript
// Type inferred from context
window.addEventListener('click', (e) => {
  // e is inferred as MouseEvent
  console.log(e.clientX);  // ✓ Knows about MouseEvent properties
});

// Array methods
const numbers = [1, 2, 3];
numbers.map(n => n.toFixed(2));  // n is inferred as number

// Return type inferred
function double(x: number) {
  return x * 2;  // Return type: number (inferred)
}
```

## Type Widening

TypeScript "widens" literal types to make them more general.

```typescript
// Widening
let x = 'hello';      // Type: string (not "hello")
let y = 42;           // Type: number (not 42)
let z = true;         // Type: boolean (not true)

// Prevent widening with const
const a = 'hello';    // Type: "hello" (literal)
const b = 42;         // Type: 42 (literal)

// Object widening
let obj = { x: 1 };   // Type: { x: number }
const obj2 = { x: 1 }; // Still: { x: number }
// Note: const doesn't prevent widening for objects!

// Prevent with as const
const obj3 = { x: 1 } as const;  // Type: { readonly x: 1 }
```

### Interview Question: const vs as const

```typescript
const config = {
  apiUrl: 'https://api.example.com',
  timeout: 5000
};
// Type: { apiUrl: string; timeout: number }

const config2 = {
  apiUrl: 'https://api.example.com',
  timeout: 5000
} as const;
// Type: { readonly apiUrl: "https://api.example.com"; readonly timeout: 5000 }

// Usage difference:
config.timeout = 10000;   // ✓ Allowed
config2.timeout = 10000;  // ❌ Error: readonly
```

## Type Narrowing

TypeScript refines types based on control flow analysis.

### typeof Guards

```typescript
function process(value: string | number) {
  if (typeof value === 'string') {
    return value.toUpperCase();  // Type: string
  } else {
    return value.toFixed(2);     // Type: number
  }
}
```

### instanceof Guards

```typescript
class Dog {
  bark() {}
}
class Cat {
  meow() {}
}

function speak(animal: Dog | Cat) {
  if (animal instanceof Dog) {
    animal.bark();  // Type: Dog
  } else {
    animal.meow();  // Type: Cat
  }
}
```

### in Operator

```typescript
interface Fish {
  swim(): void;
}
interface Bird {
  fly(): void;
}

function move(animal: Fish | Bird) {
  if ('swim' in animal) {
    animal.swim();  // Type: Fish
  } else {
    animal.fly();   // Type: Bird
  }
}
```

### Discriminated Unions (Tagged Unions)

```typescript
interface Success {
  kind: 'success';
  data: string;
}
interface Error {
  kind: 'error';
  message: string;
}
type Result = Success | Error;

function handle(result: Result) {
  if (result.kind === 'success') {
    console.log(result.data);  // Type: Success
  } else {
    console.log(result.message);  // Type: Error
  }
}
```

### Truthiness Narrowing

```typescript
function printLength(str: string | null) {
  if (str) {
    console.log(str.length);  // Type: string
  } else {
    console.log('No string');  // Type: null
  }
}

// Caveat: Excludes falsy values
function printLength2(str: string | null) {
  if (str) {
    console.log(str.length);  // Type: string
  }
  // But empty string '' is also falsy!
}

// Better:
function printLength3(str: string | null) {
  if (str !== null) {
    console.log(str.length);  // Type: string (includes '')
  }
}
```

### Custom Type Guards

```typescript
interface User {
  name: string;
  email: string;
}

// Type predicate
function isUser(obj: any): obj is User {
  return (
    typeof obj === 'object' &&
    obj !== null &&
    typeof obj.name === 'string' &&
    typeof obj.email === 'string'
  );
}

function greet(data: unknown) {
  if (isUser(data)) {
    console.log(data.name);  // Type: User
  }
}
```

### Assertion Functions

```typescript
function assert(condition: any, msg?: string): asserts condition {
  if (!condition) {
    throw new Error(msg);
  }
}

function assertIsString(val: any): asserts val is string {
  if (typeof val !== 'string') {
    throw new Error('Not a string');
  }
}

function process(value: string | number) {
  assertIsString(value);
  // Type is now: string (TypeScript knows we threw if not)
  value.toUpperCase();
}
```

## Structural Typing (Duck Typing)

TypeScript uses structural typing, not nominal typing.

### Structural vs Nominal

```typescript
// Structural: Compatibility based on shape
interface Point2D {
  x: number;
  y: number;
}

interface Vector {
  x: number;
  y: number;
}

const point: Point2D = { x: 0, y: 0 };
const vec: Vector = point;  // ✓ Compatible (same structure)

// Nominal typing (e.g., Java, C#):
// Point2D and Vector would be incompatible despite same shape
```

### Excess Property Checking

```typescript
interface Config {
  url: string;
  timeout?: number;
}

// ✓ Allowed (structural)
const config1: Config = {
  url: 'https://api.com',
  timeout: 5000,
  retries: 3  // ❌ Error: Object literal may only specify known properties
};

// ✓ Workaround: Assign to variable first
const temp = {
  url: 'https://api.com',
  timeout: 5000,
  retries: 3
};
const config2: Config = temp;  // ✓ No error (excess property check bypassed)

// ✓ Type assertion
const config3: Config = {
  url: 'https://api.com',
  retries: 3
} as Config;
```

**Why?** Excess property checking only applies to object literals to catch typos.

### Subtyping

```typescript
interface Animal {
  name: string;
}

interface Dog extends Animal {
  breed: string;
}

const dog: Dog = { name: 'Buddy', breed: 'Golden' };
const animal: Animal = dog;  // ✓ Dog is subtype of Animal (more properties OK)

// Reverse doesn't work:
const animal2: Animal = { name: 'Unknown' };
const dog2: Dog = animal2;  // ❌ Error: missing 'breed'
```

## Variance

Understanding how types relate in generic positions.

### Covariance (Most Common)

**Read-only positions** are covariant.

```typescript
interface Animal {
  name: string;
}
interface Dog extends Animal {
  breed: string;
}

// Arrays are covariant in read position
let animals: Animal[] = [];
let dogs: Dog[] = [{ name: 'Buddy', breed: 'Golden' }];

animals = dogs;  // ✓ Dog[] is subtype of Animal[]

// Reading is safe:
const animal: Animal = animals[0];  // ✓ Dog is Animal

// But writing is unsafe (in reality):
// animals.push({ name: 'Cat' });  // Would break dogs array!
// TypeScript allows this - limitation
```

### Contravariance (Function Parameters)

**Write-only positions** are contravariant.

```typescript
type Func<T> = (arg: T) => void;

let animalFunc: Func<Animal> = (animal) => {
  console.log(animal.name);
};

let dogFunc: Func<Dog> = (dog) => {
  console.log(dog.breed);
};

// Contravariance: Can assign Animal handler to Dog position
dogFunc = animalFunc;  // ✓ With strictFunctionTypes: true

// Why? Every Dog is an Animal, so Animal handler can handle Dogs
// Reverse is unsafe:
animalFunc = dogFunc;  // ❌ Error: Cat doesn't have 'breed'
```

### Invariance (Read-Write)

**Mutable properties** are invariant.

```typescript
interface Box<T> {
  value: T;
  set(val: T): void;
  get(): T;
}

let animalBox: Box<Animal>;
let dogBox: Box<Dog>;

animalBox = dogBox;  // ❌ Error: Invariant
dogBox = animalBox;  // ❌ Error: Invariant

// Why? Both read and write:
// If allowed: animalBox.set({ name: 'Cat' }) would corrupt dogBox
```

### Bivariance (Legacy)

```typescript
// strictFunctionTypes: false (old behavior)
type Handler<T> = (arg: T) => void;

let animalHandler: Handler<Animal>;
let dogHandler: Handler<Dog>;

// Both allowed (bivariant - unsafe!)
animalHandler = dogHandler;  // ⚠️ Allowed
dogHandler = animalHandler;  // ⚠️ Allowed

// Always use strictFunctionTypes: true
```

## Type Compatibility

### Function Compatibility

```typescript
type Func1 = (a: number, b: number) => number;
type Func2 = (a: number) => number;

let f1: Func1 = (a, b) => a + b;
let f2: Func2 = (a) => a * 2;

// Can assign fewer parameters
f1 = f2;  // ✓ Allowed (extra params ignored)
f2 = f1;  // ❌ Error (missing required param)

// Array.forEach example:
[1, 2, 3].forEach((n) => console.log(n));  // ✓ Ignores index, array params
```

### Return Type Compatibility

```typescript
type GetAnimal = () => Animal;
type GetDog = () => Dog;

let getAnimal: GetAnimal;
let getDog: GetDog;

// Return type is covariant
getAnimal = getDog;  // ✓ Returning Dog is safe (Dog is Animal)
getDog = getAnimal;  // ❌ Error: Returning Animal might not be Dog
```

## Hands-On Exercise 1: Narrowing Challenge

Fix type errors using narrowing:

```typescript
function process(value: string | number | null) {
  // Goal: Call .toUpperCase() if string, .toFixed() if number
  // Skip if null
}
```

<details>
<summary>Solution</summary>

```typescript
function process(value: string | number | null) {
  if (value === null) {
    return;
  }

  if (typeof value === 'string') {
    console.log(value.toUpperCase());
  } else {
    console.log(value.toFixed(2));
  }
}
```

</details>

## Hands-On Exercise 2: Type Guard

Write a type guard for this shape:

```typescript
interface ApiSuccess {
  status: 'success';
  data: any;
}
interface ApiError {
  status: 'error';
  error: string;
}
type ApiResponse = ApiSuccess | ApiError;

// Write: function isSuccess(response: ApiResponse): response is ApiSuccess
```

<details>
<summary>Solution</summary>

```typescript
function isSuccess(response: ApiResponse): response is ApiSuccess {
  return response.status === 'success';
}

// Usage:
function handle(response: ApiResponse) {
  if (isSuccess(response)) {
    console.log(response.data);  // Type: ApiSuccess
  } else {
    console.log(response.error);  // Type: ApiError
  }
}
```

</details>

## Hands-On Exercise 3: Variance

Explain why this works or doesn't:

```typescript
interface Animal { name: string; }
interface Dog extends Animal { breed: string; }

type AnimalFunc = (animal: Animal) => void;
type DogFunc = (dog: Dog) => void;

let f1: AnimalFunc = (a) => console.log(a.name);
let f2: DogFunc = (d) => console.log(d.breed);

// Will these work? Why?
f1 = f2;  // ?
f2 = f1;  // ?
```

<details>
<summary>Solution</summary>

```typescript
f1 = f2;  // ❌ Error
// AnimalFunc might receive Cat, but f2 expects Dog (needs .breed)

f2 = f1;  // ✓ Allowed (contravariance)
// DogFunc only receives Dogs, and f1 can handle any Animal
// Since Dog extends Animal, f1 is safe to use

// Parameters are contravariant!
```

</details>

## Interview Questions

### Q1: What's type widening?

<details>
<summary>Answer</summary>

TypeScript converts literal types to general types for mutability.

```typescript
let x = 'hello';  // Type: string (widened from "hello")
const y = 'hello';  // Type: "hello" (literal)

// Reason: let can be reassigned
x = 'world';  // Must allow any string

// Prevent with as const:
let z = 'hello' as const;  // Type: "hello"
z = 'world';  // ❌ Error
```

</details>

### Q2: Explain structural vs nominal typing

<details>
<summary>Answer</summary>

**Structural** (TypeScript): Compatibility based on shape/structure.

```typescript
interface A { x: number; }
interface B { x: number; }

const a: A = { x: 1 };
const b: B = a;  // ✓ Same structure
```

**Nominal** (Java, C#, Flow): Compatibility based on explicit declarations.

```java
class A { int x; }
class B { int x; }

A a = new A();
B b = a;  // ❌ Error: Different types despite same structure
```

TypeScript chose structural for JavaScript's duck-typed nature.

</details>

### Q3: How do discriminated unions work?

<details>
<summary>Answer</summary>

Use a common literal property to discriminate union members.

```typescript
type Shape =
  | { kind: 'circle'; radius: number }
  | { kind: 'square'; size: number };

function area(shape: Shape) {
  if (shape.kind === 'circle') {
    return Math.PI * shape.radius ** 2;  // Type: circle
  } else {
    return shape.size ** 2;  // Type: square
  }
}
```

TypeScript narrows based on the discriminant (`kind`).

</details>

### Q4: What's the difference between type predicates and assertion functions?

<details>
<summary>Answer</summary>

**Type Predicate**: Returns boolean, narrows type in if-block.

```typescript
function isString(val: any): val is string {
  return typeof val === 'string';
}

if (isString(x)) {
  x.toUpperCase();  // Type: string
}
```

**Assertion Function**: Throws on false, narrows type after call.

```typescript
function assertString(val: any): asserts val is string {
  if (typeof val !== 'string') throw new Error();
}

assertString(x);
x.toUpperCase();  // Type: string (if we reach here)
```

</details>

### Q5: Explain function parameter contravariance

<details>
<summary>Answer</summary>

Function parameters are contravariant: can assign handler of supertype to position expecting subtype.

```typescript
type Handler<T> = (arg: T) => void;

let animalHandler: Handler<Animal> = (a) => console.log(a.name);
let dogHandler: Handler<Dog>;

dogHandler = animalHandler;  // ✓ Contravariance
```

**Why safe?**
- `dogHandler` receives only Dogs
- `animalHandler` can handle any Animal
- Dog is an Animal, so it works

**Why reverse is unsafe?**
```typescript
animalHandler = dogHandler;  // ❌ With strictFunctionTypes
// animalHandler might receive Cat, but dogHandler expects Dog
```

</details>

## Key Takeaways

1. **Inference**: TypeScript infers types; explicit types for clarity
2. **Widening**: Literals widen to general types; prevent with `as const`
3. **Narrowing**: Use typeof, instanceof, in, discriminated unions, type guards
4. **Structural Typing**: Compatibility by shape, not name
5. **Variance**:
   - Covariance: Return types (Dog → Animal)
   - Contravariance: Parameters (Animal → Dog)
   - Invariance: Mutable (neither direction)
6. **Type Guards**: `is` for predicates, `asserts` for assertions
7. **Discriminated Unions**: Use literal discriminant for safe narrowing

## Next Steps

In [Lesson 04: Utility Types & Type Manipulation](lesson-04-utility-types-and-manipulation.md), you'll learn:
- Built-in utility types (Partial, Pick, Record, etc.)
- Mapped types
- Conditional types
- Template literal types
