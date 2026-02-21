# Lesson 06: Module Systems

Deep dive into ESM vs CommonJS, module resolution, and import/export patterns.

## ESM vs CommonJS

### CommonJS (Node.js Traditional)

```javascript
// math.js - Exporting
exports.add = (a, b) => a + b;
exports.subtract = (a, b) => a - b;

// Or
module.exports = {
  add: (a, b) => a + b,
  subtract: (a, b) => a - b
};

// Or default export
module.exports = class Calculator {
  add(a, b) { return a + b; }
};

// app.js - Importing
const math = require('./math');
const { add } = require('./math');
const Calculator = require('./calculator');
```

**Characteristics**:
- Synchronous loading
- Dynamic (can `require()` conditionally)
- Runtime resolution
- `this` equals `exports`
- File extension optional

### ESM (ECMAScript Modules)

```typescript
// math.ts - Exporting
export const add = (a: number, b: number) => a + b;
export const subtract = (a: number, b: number) => a - b;

// Or
export { add, subtract };

// Default export
export default class Calculator {
  add(a: number, b: number) { return a + b; }
}

// app.ts - Importing
import * as math from './math.js';  // Note: .js extension!
import { add } from './math.js';
import Calculator from './calculator.js';
```

**Characteristics**:
- Static analysis possible
- Async loading
- Compile-time resolution
- `this` is `undefined` at top level
- File extension required (in Node16+)

### Side-by-Side Comparison

| Feature | CommonJS | ESM |
|---------|----------|-----|
| Syntax | `require()` / `module.exports` | `import` / `export` |
| Loading | Synchronous | Async |
| When resolved | Runtime | Parse time |
| Tree-shaking | ❌ No | ✓ Yes |
| Dynamic imports | Always | Via `import()` |
| Top-level await | ❌ No | ✓ Yes (Node 14.8+) |
| File extension | Optional | Required (Node16+) |
| `this` | `exports` | `undefined` |

## TypeScript Compilation

### Output Module Format

```json
// tsconfig.json
{
  "compilerOptions": {
    "module": "CommonJS"  // or "ES2020", "ESNext", "NodeNext"
  }
}
```

```typescript
// Source (TypeScript)
import { foo } from './bar';
export const baz = 42;

// Output with "module": "CommonJS"
const bar_1 = require("./bar");
exports.baz = 42;

// Output with "module": "ES2020"
import { foo } from './bar';
export const baz = 42;
```

### Interop Flags

```json
{
  "compilerOptions": {
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true
  }
}
```

**Without esModuleInterop**:
```typescript
import * as express from 'express';  // Must use namespace import
const app = express();
```

**With esModuleInterop**:
```typescript
import express from 'express';  // Can use default import
const app = express();
```

## Module Resolution

### Resolution Strategies

```json
{
  "compilerOptions": {
    "moduleResolution": "node" | "node16" | "nodenext" | "bundler"
  }
}
```

| Strategy | Use Case | Behavior |
|----------|----------|----------|
| `node` | Legacy Node | Classic Node resolution |
| `node16` / `nodenext` | Modern Node (16+) | Respects package.json `exports`, requires extensions |
| `bundler` | Webpack/Vite | Like node16 but relaxed (no extensions needed) |

### Classic Resolution (node)

```
import { foo } from './bar';

Looks for:
1. ./bar.ts
2. ./bar.tsx
3. ./bar.d.ts
4. ./bar/index.ts
5. ./bar/index.tsx
6. ./bar/index.d.ts
```

### Node16/NodeNext Resolution

```typescript
// package.json
{ "type": "module" }

// MUST include extension:
import { foo } from './bar.js';  // ✓ (looks for bar.ts)
import { foo } from './bar';     // ❌ Error

// Respects package.json exports:
import { foo } from 'my-lib';    // Uses "exports" field
```

**Critical**: Use `.js` extension even for `.ts` files.

```typescript
// file: src/utils.ts
export const helper = () => {};

// file: src/app.ts
import { helper } from './utils.js';  // ✓ Correct
// TypeScript finds utils.ts, emits utils.js
```

### Package Exports Resolution

```json
// node_modules/my-lib/package.json
{
  "exports": {
    ".": {
      "types": "./dist/index.d.ts",
      "import": "./dist/index.mjs",
      "require": "./dist/index.cjs"
    },
    "./utils": "./dist/utils.js"
  }
}
```

```typescript
import lib from 'my-lib';          // Uses exports["."]
import { util } from 'my-lib/utils';  // Uses exports["./utils"]
import { foo } from 'my-lib/internal';  // ❌ Not in exports
```

## Import/Export Patterns

### Named Exports

```typescript
// Inline export
export const PI = 3.14;
export function add(a: number, b: number) {
  return a + b;
}

// Export list
const PI = 3.14;
function add(a: number, b: number) {
  return a + b;
}
export { PI, add };

// Re-export
export { foo, bar } from './other';
export * from './other';  // Re-export all
```

### Default Export

```typescript
// Only one per module
export default class Calculator {
  // ...
}

// Or
class Calculator {
  // ...
}
export default Calculator;

// Or inline value
export default { version: '1.0.0' };
```

### Import Patterns

```typescript
// Named imports
import { foo, bar } from './module';

// Aliasing
import { foo as f, bar as b } from './module';

// Namespace import
import * as utils from './utils';
utils.foo();

// Default import
import Calculator from './calculator';

// Mixed
import Calculator, { add, subtract } from './math';

// Side-effects only
import './polyfills';
```

### Re-exporting

```typescript
// Re-export named
export { foo, bar } from './other';

// Re-export all
export * from './other';

// Re-export with rename
export { foo as newFoo } from './other';

// Re-export default as named
export { default as Calculator } from './calculator';
```

## Dynamic Imports

```typescript
// Static import - always loaded
import { large } from './large-module';

// Dynamic import - loaded on demand
async function loadFeature() {
  const { large } = await import('./large-module');
  large();
}

// Conditional loading
if (condition) {
  const module = await import('./conditional');
}

// Type-safe dynamic import
type MathModule = typeof import('./math');

async function getMath(): Promise<MathModule> {
  return import('./math');
}
```

### Use Cases for Dynamic Imports

1. **Code splitting** (reduce initial bundle)
2. **Conditional loading** (feature flags)
3. **Lazy loading** (on user interaction)
4. **Dynamic module paths**

```typescript
// Route-based loading
async function loadRoute(route: string) {
  const module = await import(`./routes/${route}`);
  return module.default;
}

// Feature flags
if (featureFlags.newFeature) {
  const { NewFeature } = await import('./new-feature');
  // Use NewFeature
}
```

## Top-Level Await

Only works in ESM (Node 14.8+, "module": "ES2022"+).

```typescript
// ❌ CommonJS
const data = await fetch('/api/data');  // Error: await in top-level

// ✓ ESM (package.json: { "type": "module" })
const response = await fetch('/api/data');
const data = await response.json();

export const config = data;
```

**Use carefully**: Blocks module execution, affects load time.

## Import Assertions (Type Attributes)

For importing non-JavaScript files (Node 17+).

```typescript
// JSON
import data from './data.json' assert { type: 'json' };

// CSS (in bundlers)
import styles from './styles.css' assert { type: 'css' };
```

**New syntax (TypeScript 5.3+)**:
```typescript
import data from './data.json' with { type: 'json' };
```

## CommonJS Interop

### Importing CommonJS in ESM

```typescript
// CommonJS module: const foo = { bar: 42 }; module.exports = foo;

// ESM import
import foo from './cjs-module';  // ✓ Works with esModuleInterop
import * as foo from './cjs-module';  // ✓ Always works

console.log(foo.bar);
```

### Importing ESM in CommonJS

```javascript
// ❌ Can't use static import in CommonJS
import { foo } from './esm-module';  // Error

// ✓ Use dynamic import
async function load() {
  const { foo } = await import('./esm-module.mjs');
  console.log(foo);
}
```

## Module Augmentation

Extend existing modules with new declarations.

```typescript
// node_modules/express/index.d.ts
declare namespace Express {
  interface Request {
    user?: User;
  }
}

// Your code: extend Express.Request
declare module 'express' {
  interface Request {
    customProp: string;
  }
}

// Now available:
app.get('/', (req, res) => {
  console.log(req.customProp);  // ✓ TypeScript knows about it
});
```

### Global Augmentation

```typescript
// Extend global namespace
declare global {
  interface Window {
    myLib: MyLibrary;
  }
}

// Must export something to be a module
export {};
```

## Hands-On Exercise 1: Module Migration

Convert this CommonJS to ESM:

```javascript
// math.js
const add = (a, b) => a + b;
const subtract = (a, b) => a - b;
exports.add = add;
exports.subtract = subtract;

// app.js
const { add } = require('./math');
console.log(add(1, 2));
```

<details>
<summary>Solution</summary>

```typescript
// math.ts
export const add = (a: number, b: number) => a + b;
export const subtract = (a: number, b: number) => a - b;

// app.ts
import { add } from './math.js';  // Note: .js extension
console.log(add(1, 2));

// package.json
{
  "type": "module"
}

// tsconfig.json
{
  "compilerOptions": {
    "module": "NodeNext",
    "moduleResolution": "NodeNext"
  }
}
```

</details>

## Hands-On Exercise 2: Dynamic Import

Implement lazy loading for a feature:

```typescript
// Goal: Load analytics only when user opts in

interface Analytics {
  track(event: string): void;
}

async function initAnalytics(enabled: boolean): Promise<Analytics | null> {
  // Implement
}
```

<details>
<summary>Solution</summary>

```typescript
// analytics.ts
export interface Analytics {
  track(event: string): void;
  pageView(url: string): void;
}

export function createAnalytics(): Analytics {
  return {
    track(event: string) {
      console.log('Track:', event);
    },
    pageView(url: string) {
      console.log('Page view:', url);
    }
  };
}

// app.ts
async function initAnalytics(enabled: boolean): Promise<Analytics | null> {
  if (!enabled) {
    return null;
  }

  const { createAnalytics } = await import('./analytics.js');
  return createAnalytics();
}

// Usage
const analytics = await initAnalytics(true);
analytics?.track('app_start');
```

</details>

## Interview Questions

### Q1: ESM vs CommonJS - key differences?

<details>
<summary>Answer</summary>

| Aspect | CommonJS | ESM |
|--------|----------|-----|
| **When evaluated** | Runtime | Parse time |
| **Loading** | Synchronous | Async |
| **Tree-shaking** | No | Yes |
| **Dynamic imports** | Built-in | Via `import()` |
| **Extensions** | Optional | Required (Node16+) |
| **Top-level await** | No | Yes |

**Why ESM?**
- Better for bundlers (tree-shaking)
- Standard across environments
- Static analysis enables optimization

</details>

### Q2: Why require .js extension in imports with NodeNext?

<details>
<summary>Answer</summary>

```typescript
import { foo } from './bar.js';  // Why .js for .ts file?
```

**Reason**: TypeScript emits code that Node will run. Node resolves imports in the *emitted* .js files.

```typescript
// Source: src/app.ts
import { foo } from './utils.js';

// Emitted: dist/app.js
import { foo } from './utils.js';  // Node looks for utils.js
```

TypeScript finds the source (.ts) but emits the extension you specified (.js).

</details>

### Q3: What's esModuleInterop?

<details>
<summary>Answer</summary>

Enables default imports from CommonJS modules.

```typescript
// CommonJS module: module.exports = function() {}

// Without esModuleInterop
import * as express from 'express';
const app = express();  // ❌ Error: not callable

// With esModuleInterop: true
import express from 'express';
const app = express();  // ✓ Works
```

**What it does**: Adds runtime helpers to make CJS default exports compatible with ESM default imports.

</details>

### Q4: When to use dynamic imports?

<details>
<summary>Answer</summary>

**Use cases**:
1. **Code splitting**: Reduce initial bundle
   ```typescript
   const module = await import('./large-feature');
   ```

2. **Conditional loading**: Feature flags
   ```typescript
   if (featureEnabled) {
     await import('./new-feature');
   }
   ```

3. **Lazy loading**: On-demand
   ```typescript
   button.onclick = async () => {
     const { modal } = await import('./modal');
     modal.show();
   };
   ```

4. **Dynamic paths**:
   ```typescript
   const locale = await import(`./i18n/${lang}.js`);
   ```

</details>

### Q5: How does module augmentation work?

<details>
<summary>Answer</summary>

Extend existing module types:

```typescript
// Extend external module
declare module 'express' {
  interface Request {
    userId?: string;
  }
}

// Now available:
app.use((req, res, next) => {
  req.userId = '123';  // ✓ TypeScript knows
});
```

**Requirements**:
- Must use `declare module 'exact-module-name'`
- File must be a module (has import/export)
- Works with interface merging

</details>

## Key Takeaways

1. **ESM vs CJS**: ESM is standard, better for bundlers, async loading
2. **Module resolution**: Use `NodeNext` for modern Node, `bundler` for webpack/vite
3. **Extensions**: Required in Node16+ ESM (.js for .ts files)
4. **esModuleInterop**: Essential for importing CJS with default syntax
5. **Dynamic imports**: Use for code splitting, lazy loading, conditional features
6. **Top-level await**: Only in ESM, blocks module loading
7. **Module augmentation**: Extend external module types

## Next Steps

In [Lesson 07: Declaration Files](lesson-07-declaration-files.md), you'll learn:
- Writing .d.ts files
- Ambient declarations
- Triple-slash directives
- DefinitelyTyped contributions
