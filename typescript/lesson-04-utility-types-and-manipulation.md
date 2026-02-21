# Lesson 04: Utility Types & Type Manipulation

Master TypeScript's built-in utility types and advanced type manipulation techniques.

## Built-in Utility Types

### Partial<T>

Makes all properties optional.

```typescript
interface User {
  name: string;
  email: string;
  age: number;
}

type PartialUser = Partial<User>;
// { name?: string; email?: string; age?: number; }

// Common use: update functions
function updateUser(id: number, updates: Partial<User>) {
  // Can pass any subset of properties
}

updateUser(1, { age: 30 });  // ✓ Only updating age
```

**Implementation**:
```typescript
type Partial<T> = {
  [P in keyof T]?: T[P];
};
```

### Required<T>

Makes all properties required (opposite of Partial).

```typescript
interface Config {
  apiUrl?: string;
  timeout?: number;
}

type RequiredConfig = Required<Config>;
// { apiUrl: string; timeout: number; }

// Ensure all fields are set
function validateConfig(config: Required<Config>) {
  // config.apiUrl and config.timeout are guaranteed
}
```

**Implementation**:
```typescript
type Required<T> = {
  [P in keyof T]-?: T[P];  // -? removes optionality
};
```

### Readonly<T>

Makes all properties readonly.

```typescript
interface Todo {
  title: string;
  completed: boolean;
}

const todo: Readonly<Todo> = {
  title: 'Learn TypeScript',
  completed: false
};

todo.completed = true;  // ❌ Error: readonly property
```

**Implementation**:
```typescript
type Readonly<T> = {
  readonly [P in keyof T]: T[P];
};
```

### Pick<T, K>

Selects subset of properties.

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  password: string;
}

type PublicUser = Pick<User, 'id' | 'name' | 'email'>;
// { id: number; name: string; email: string; }

// Return safe user data
function getPublicUser(user: User): PublicUser {
  return {
    id: user.id,
    name: user.name,
    email: user.email
  };
}
```

**Implementation**:
```typescript
type Pick<T, K extends keyof T> = {
  [P in K]: T[P];
};
```

### Omit<T, K>

Removes properties (opposite of Pick).

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  password: string;
}

type UserWithoutPassword = Omit<User, 'password'>;
// { id: number; name: string; email: string; }

type UserIdOnly = Omit<User, 'name' | 'email' | 'password'>;
// { id: number; }
```

**Implementation**:
```typescript
type Omit<T, K extends keyof any> = Pick<T, Exclude<keyof T, K>>;
```

### Record<K, T>

Creates object type with keys K and values T.

```typescript
type Role = 'admin' | 'user' | 'guest';

const permissions: Record<Role, string[]> = {
  admin: ['read', 'write', 'delete'],
  user: ['read', 'write'],
  guest: ['read']
};

// Ensures all roles are defined
// permissions.superadmin = [...];  // ❌ Error
```

**Common pattern**: Dictionary/map type
```typescript
const userCache: Record<string, User> = {};
userCache['123'] = { id: 123, name: 'Alice' };
```

**Implementation**:
```typescript
type Record<K extends keyof any, T> = {
  [P in K]: T;
};
```

### Exclude<T, U>

Removes types from union.

```typescript
type AllTypes = 'a' | 'b' | 'c' | 'd';
type Excluded = Exclude<AllTypes, 'a' | 'c'>;
// 'b' | 'd'

// Real-world: Remove null/undefined
type NonNullable<T> = Exclude<T, null | undefined>;

type Result = string | null | undefined;
type SafeResult = NonNullable<Result>;  // string
```

**Implementation**:
```typescript
type Exclude<T, U> = T extends U ? never : T;
// Distributes over union
```

### Extract<T, U>

Keeps only types assignable to U.

```typescript
type AllTypes = 'a' | 'b' | 'c' | 'd';
type Extracted = Extract<AllTypes, 'a' | 'c' | 'e'>;
// 'a' | 'c'

// Real-world: Extract specific types
type Mixed = string | number | (() => void);
type Functions = Extract<Mixed, Function>;  // () => void
```

**Implementation**:
```typescript
type Extract<T, U> = T extends U ? T : never;
```

### ReturnType<T>

Extracts function return type.

```typescript
function getUser() {
  return { id: 1, name: 'Alice' };
}

type User = ReturnType<typeof getUser>;
// { id: number; name: string; }

// Works with generic functions
function createPair<T>(value: T) {
  return { value, timestamp: Date.now() };
}

type Pair = ReturnType<typeof createPair<string>>;
// { value: string; timestamp: number; }
```

**Implementation**:
```typescript
type ReturnType<T extends (...args: any) => any> =
  T extends (...args: any) => infer R ? R : any;
```

### Parameters<T>

Extracts function parameter types as tuple.

```typescript
function createUser(name: string, age: number, active: boolean) {
  // ...
}

type CreateUserParams = Parameters<typeof createUser>;
// [name: string, age: number, active: boolean]

// Use to match function signature
function logAndCreate(...args: CreateUserParams) {
  console.log('Creating user with:', args);
  return createUser(...args);
}
```

**Implementation**:
```typescript
type Parameters<T extends (...args: any) => any> =
  T extends (...args: infer P) => any ? P : never;
```

### Awaited<T>

Unwraps Promise type (TypeScript 4.5+).

```typescript
type AsyncValue = Promise<string>;
type Value = Awaited<AsyncValue>;  // string

// Nested Promises
type Nested = Promise<Promise<number>>;
type Unwrapped = Awaited<Nested>;  // number

// Real-world: Infer async function return
async function fetchUser() {
  return { id: 1, name: 'Alice' };
}

type User = Awaited<ReturnType<typeof fetchUser>>;
// { id: number; name: string; }
```

## Mapped Types

Create new types by transforming properties.

### Basic Mapped Type

```typescript
type Flags<T> = {
  [K in keyof T]: boolean;
};

interface Features {
  darkMode: string;
  notifications: number;
}

type FeatureFlags = Flags<Features>;
// { darkMode: boolean; notifications: boolean; }
```

### Adding Modifiers

```typescript
// Make everything optional and readonly
type Frozen<T> = {
  readonly [K in keyof T]?: T[T];
};

// Remove readonly (using -)
type Mutable<T> = {
  -readonly [K in keyof T]: T[K];
};

// Remove optional (using -)
type Concrete<T> = {
  [K in keyof T]-?: T[K];
};
```

### Key Remapping (TypeScript 4.1+)

```typescript
// Prefix all keys with 'get'
type Getters<T> = {
  [K in keyof T as `get${Capitalize<string & K>}`]: () => T[K];
};

interface Person {
  name: string;
  age: number;
}

type PersonGetters = Getters<Person>;
// {
//   getName: () => string;
//   getAge: () => number;
// }
```

### Filtering Keys

```typescript
// Remove specific keys
type OmitType<T, V> = {
  [K in keyof T as T[K] extends V ? never : K]: T[K];
};

interface Data {
  id: number;
  name: string;
  active: boolean;
  count: number;
}

type NonNumbers = OmitType<Data, number>;
// { name: string; active: boolean; }
```

## Conditional Types

Types that depend on conditions.

### Basic Conditional

```typescript
type IsString<T> = T extends string ? true : false;

type A = IsString<string>;   // true
type B = IsString<number>;   // false
```

### With Unions (Distributive)

```typescript
type ToArray<T> = T extends any ? T[] : never;

type Result = ToArray<string | number>;
// string[] | number[] (distributes over union)
```

### Infer Keyword

Extract types from within other types.

```typescript
// Extract array element type
type ElementType<T> = T extends (infer E)[] ? E : never;

type Str = ElementType<string[]>;  // string
type Num = ElementType<number[]>;  // number

// Extract function return type (ReturnType implementation)
type GetReturn<T> = T extends (...args: any[]) => infer R ? R : never;

// Extract Promise value
type UnwrapPromise<T> = T extends Promise<infer U> ? U : T;

type A = UnwrapPromise<Promise<string>>;  // string
type B = UnwrapPromise<number>;           // number
```

### Advanced: DeepPartial

```typescript
type DeepPartial<T> = {
  [K in keyof T]?: T[K] extends object
    ? DeepPartial<T[K]>
    : T[K];
};

interface Nested {
  user: {
    profile: {
      name: string;
      age: number;
    };
  };
}

type PartialNested = DeepPartial<Nested>;
// {
//   user?: {
//     profile?: {
//       name?: string;
//       age?: number;
//     };
//   };
// }
```

## Template Literal Types

String manipulation at type level (TypeScript 4.1+).

### Basic Template Literals

```typescript
type Greeting = `Hello ${string}`;

const g1: Greeting = 'Hello World';  // ✓
const g2: Greeting = 'Hi World';     // ❌ Error

// With unions
type Color = 'red' | 'blue' | 'green';
type HexColor = `#${string}`;
type ColorPalette = Color | HexColor;
```

### String Manipulation Utilities

```typescript
type Uppercase<S extends string> = ...;  // Built-in
type Lowercase<S extends string> = ...;  // Built-in
type Capitalize<S extends string> = ...; // Built-in
type Uncapitalize<S extends string> = ...;  // Built-in

type Loud = Uppercase<'hello'>;      // 'HELLO'
type Quiet = Lowercase<'WORLD'>;     // 'world'
type Title = Capitalize<'typescript'>;  // 'Typescript'
```

### API Route Builder

```typescript
type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';
type Route = '/users' | '/posts' | '/comments';
type APIRoute = `${HTTPMethod} ${Route}`;

// 'GET /users' | 'GET /posts' | ... | 'DELETE /comments' (12 combinations)

function handleRoute(route: APIRoute) {
  // Type-safe route handling
}

handleRoute('GET /users');      // ✓
handleRoute('PATCH /users');    // ❌ Error
```

### Event Names

```typescript
type EventName<T extends string> = `on${Capitalize<T>}`;

type Events = 'click' | 'focus' | 'blur';
type EventHandlers = EventName<Events>;
// 'onClick' | 'onFocus' | 'onBlur'

type Handlers = {
  [E in Events as EventName<E>]: (event: Event) => void;
};
// {
//   onClick: (event: Event) => void;
//   onFocus: (event: Event) => void;
//   onBlur: (event: Event) => void;
// }
```

## Hands-On Exercise 1: Build Custom Utilities

Create these utility types:

```typescript
// 1. PickByType<T, V> - Pick properties of specific type
interface Data {
  id: number;
  name: string;
  active: boolean;
  count: number;
}

type Numbers = PickByType<Data, number>;
// { id: number; count: number; }

// 2. FunctionKeys<T> - Extract function property names
interface API {
  baseUrl: string;
  timeout: number;
  get: () => void;
  post: () => void;
}

type Methods = FunctionKeys<API>;
// 'get' | 'post'
```

<details>
<summary>Solution</summary>

```typescript
// 1. PickByType
type PickByType<T, V> = {
  [K in keyof T as T[K] extends V ? K : never]: T[K];
};

// 2. FunctionKeys
type FunctionKeys<T> = {
  [K in keyof T]: T[K] extends Function ? K : never;
}[keyof T];

// Alternative:
type FunctionKeys<T> = keyof {
  [K in keyof T as T[K] extends Function ? K : never]: T[K];
};
```

</details>

## Hands-On Exercise 2: Event Emitter Types

Build type-safe event emitter:

```typescript
interface Events {
  click: { x: number; y: number };
  submit: { value: string };
  error: { message: string };
}

// Create:
// - EventMap that converts to { 'click': (data: { x, y }) => void, ... }
// - EventNames type union
```

<details>
<summary>Solution</summary>

```typescript
type EventMap<T> = {
  [K in keyof T]: (data: T[K]) => void;
};

type EventNames<T> = keyof T;

// Usage:
class EventEmitter<T> {
  private handlers: Partial<EventMap<T>> = {};

  on<K extends EventNames<T>>(event: K, handler: EventMap<T>[K]) {
    this.handlers[event] = handler;
  }

  emit<K extends EventNames<T>>(event: K, data: T[K]) {
    this.handlers[event]?.(data);
  }
}

const emitter = new EventEmitter<Events>();

emitter.on('click', (data) => {
  console.log(data.x, data.y);  // Type-safe: { x, y }
});

emitter.emit('click', { x: 10, y: 20 });  // ✓
emitter.emit('click', { x: 10 });         // ❌ Error: missing y
```

</details>

## Interview Questions

### Q1: Difference between Pick and Omit?

<details>
<summary>Answer</summary>

- **Pick<T, K>**: Include only specified keys
- **Omit<T, K>**: Exclude specified keys

```typescript
interface User {
  id: number;
  name: string;
  email: string;
}

type A = Pick<User, 'id' | 'name'>;    // { id, name }
type B = Omit<User, 'email'>;          // { id, name }
// Same result, different approach
```

**When to use**:
- Pick: Few properties to include
- Omit: Few properties to exclude

</details>

### Q2: How does `infer` work?

<details>
<summary>Answer</summary>

`infer` declares a type variable within a conditional type.

```typescript
// Extract return type
type ReturnType<T> =
  T extends (...args: any[]) => infer R  // Declare R
    ? R                                   // Use R
    : never;

// Extract array element
type ElementType<T> = T extends (infer E)[] ? E : never;
```

Only valid in `extends` clause of conditional types.

</details>

### Q3: What's distributive conditional type?

<details>
<summary>Answer</summary>

When conditional type is applied to union, it distributes over each member.

```typescript
type ToArray<T> = T extends any ? T[] : never;

type Result = ToArray<string | number>;
// = ToArray<string> | ToArray<number>
// = string[] | number[]

// Non-distributive (using tuple):
type ToArrayNonDist<T> = [T] extends [any] ? T[] : never;

type Result2 = ToArrayNonDist<string | number>;
// = (string | number)[]
```

</details>

### Q4: How to make all properties of nested objects optional?

<details>
<summary>Answer</summary>

```typescript
type DeepPartial<T> = {
  [K in keyof T]?: T[K] extends object
    ? DeepPartial<T[K]>
    : T[K];
};

// Better version (handles arrays):
type DeepPartial<T> = T extends object
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : T;
```

</details>

## Key Takeaways

1. **Utility Types**: Master Partial, Pick, Omit, Record, ReturnType, Parameters
2. **Mapped Types**: Transform properties with `[K in keyof T]`
3. **Modifiers**: Use `?` `-?` `readonly` `-readonly`
4. **Conditional Types**: `T extends U ? X : Y`
5. **Infer**: Extract types with `infer` in conditional types
6. **Template Literals**: Build string types, combine with unions
7. **Key Remapping**: Filter/transform keys with `as`

## Next Steps

In [Lesson 05: Generics Deep Dive](lesson-05-generics-deep-dive.md), you'll learn:
- Generic constraints
- Default type parameters
- Variance in generics
- Generic performance considerations
