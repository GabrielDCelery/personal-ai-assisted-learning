# Lesson 01: package.json Deep Dive

Critical knowledge about package.json for TypeScript developers - dependency management, versioning, and publishing.

## Dependency Types

| Type | When to Use | Installed | Bundled |
|------|-------------|-----------|---------|
| `dependencies` | Runtime requirements | Always | Yes (if published) |
| `devDependencies` | Build/test tools only | Only in dev | No |
| `peerDependencies` | Plugin/extension deps | Manual (shows warning) | No |
| `optionalDependencies` | Nice-to-have, fallback if fails | Tries, continues on fail | Yes |
| `bundledDependencies` | Force specific versions | Always | Yes (literally bundled) |

### Common Mistakes

```json
{
  "dependencies": {
    "typescript": "^5.0.0"  // ❌ Wrong - TypeScript is a build tool
  },
  "devDependencies": {
    "typescript": "^5.0.0", // ✓ Correct
    "express": "^4.18.0"    // ❌ Wrong - express is needed at runtime
  }
}
```

### Peer Dependencies Deep Dive

**Use case**: Your library extends/plugins into another library.

```json
{
  "name": "my-react-component",
  "peerDependencies": {
    "react": "^18.0.0",      // User must have React 18+
    "react-dom": "^18.0.0"
  },
  "peerDependenciesMeta": {
    "react-dom": {
      "optional": true        // Won't warn if missing
    }
  },
  "devDependencies": {
    "react": "^18.0.0",       // Still needed for development
    "react-dom": "^18.0.0"
  }
}
```

**Why?** Prevents version conflicts. If your library bundled React, users would have two React copies.

## Semver (Semantic Versioning)

Format: `MAJOR.MINOR.PATCH` (e.g., `1.4.2`)

| Version | Meaning | Allows |
|---------|---------|--------|
| `1.4.2` | Exact version | Only 1.4.2 |
| `^1.4.2` | Compatible (caret) | >=1.4.2 <2.0.0 |
| `~1.4.2` | Approximately (tilde) | >=1.4.2 <1.5.0 |
| `>=1.4.2 <2.0.0` | Range | 1.4.2 to <2.0.0 |
| `*` or `latest` | Latest | Any version |

### Interview Question: ^ vs ~

```bash
# Current version: 1.4.2
^1.4.2  →  Can install 1.4.3, 1.5.0, 1.99.0  (NOT 2.0.0)
~1.4.2  →  Can install 1.4.3, 1.4.99         (NOT 1.5.0)

# Current version: 0.4.2 (pre-1.0)
^0.4.2  →  Can install 0.4.3, 0.4.99         (NOT 0.5.0) ⚠️ Different behavior!
~0.4.2  →  Can install 0.4.3, 0.4.99         (NOT 0.5.0)
```

**Gotcha**: For `0.x.y` versions, `^` behaves like `~` because any minor change can be breaking.

### Lock Files

| File | Manager | Purpose |
|------|---------|---------|
| `package-lock.json` | npm | Locks exact dependency tree |
| `yarn.lock` | yarn | Locks exact dependency tree |
| `pnpm-lock.yaml` | pnpm | Locks exact dependency tree |

**Critical**: Always commit lock files. They ensure reproducible builds.

## Entry Points for Publishing

Modern packages need multiple entry points for different environments.

```json
{
  "name": "my-library",
  "version": "1.0.0",

  // Legacy Node.js (CommonJS)
  "main": "./dist/index.js",

  // Legacy bundlers (ESM)
  "module": "./dist/index.mjs",

  // TypeScript types
  "types": "./dist/index.d.ts",

  // Modern - most important for Node 16+
  "exports": {
    ".": {
      "types": "./dist/index.d.ts",
      "import": "./dist/index.mjs",    // ESM
      "require": "./dist/index.js",    // CommonJS
      "default": "./dist/index.js"
    },
    "./package.json": "./package.json",
    "./utils": {
      "types": "./dist/utils.d.ts",
      "import": "./dist/utils.mjs",
      "require": "./dist/utils.js"
    }
  }
}
```

### Export Conditions Order Matters

```json
{
  "exports": {
    ".": {
      "types": "./dist/index.d.ts",    // ✓ Types first (for TypeScript)
      "import": "./dist/index.mjs",    // ✓ ESM
      "require": "./dist/index.js",    // ✓ CJS
      "default": "./dist/index.js"     // ✓ Fallback
    }
  }
}
```

### Subpath Exports

```json
{
  "exports": {
    ".": "./dist/index.js",
    "./utils": "./dist/utils.js",
    "./internal/*": null              // ❌ Block access to internals
  }
}
```

```typescript
// Users can import:
import lib from 'my-library';          // ✓ Works
import utils from 'my-library/utils';  // ✓ Works
import foo from 'my-library/internal/foo'; // ❌ Error!
```

## Publishing-Critical Fields

```json
{
  "name": "@scope/package-name",  // Scoped package
  "version": "1.0.0",
  "description": "Shows in npm search",
  "keywords": ["searchable", "terms"],
  "author": "Your Name <email@example.com>",
  "license": "MIT",
  "repository": {
    "type": "git",
    "url": "https://github.com/user/repo.git"
  },
  "homepage": "https://github.com/user/repo#readme",
  "bugs": {
    "url": "https://github.com/user/repo/issues"
  },

  // What gets published
  "files": [
    "dist",
    "README.md",
    "LICENSE"
  ],

  // Prevent accidental publishing
  "private": true,  // Set to false when ready to publish

  // npm version compatibility
  "engines": {
    "node": ">=16.0.0",
    "npm": ">=8.0.0"
  },

  // Type of module system
  "type": "module"  // or "commonjs" (default)
}
```

### Files Field

```json
{
  "files": ["dist"]  // Only publishes dist/ directory
}
```

**Always included** (can't exclude):
- package.json
- README
- LICENSE
- CHANGELOG

**Always excluded** (can't include):
- node_modules
- .git

**Use `.npmignore`** to exclude additional files (similar to .gitignore).

## Scripts & Lifecycle Hooks

### Common Scripts

```json
{
  "scripts": {
    "build": "tsc",
    "test": "jest",
    "lint": "eslint src",
    "format": "prettier --write src",

    // Composite scripts
    "clean": "rm -rf dist",
    "prebuild": "npm run clean",      // Runs BEFORE build
    "postbuild": "npm run test",      // Runs AFTER build

    // Development
    "dev": "tsc --watch",
    "start": "node dist/index.js"
  }
}
```

### npm Lifecycle Hooks

Automatically run in this order:

```
npm install:
  preinstall → install → postinstall

npm publish:
  prepublishOnly → prepare → prepublish (deprecated) → publish → postpublish
```

### Publishing Workflow

```json
{
  "scripts": {
    "clean": "rm -rf dist",
    "build": "tsc",
    "test": "jest",

    // Runs before npm publish (one-time publish only)
    "prepublishOnly": "npm run build && npm test",

    // Runs on both npm install and npm publish
    "prepare": "npm run build"
  }
}
```

**Key difference**:
- `prepublishOnly`: Only before `npm publish`
- `prepare`: Before `npm publish` AND after `npm install` (useful for git dependencies)

## Package Types

```json
{
  "type": "module"  // Default: "commonjs"
}
```

| Type | .js files are | .mjs | .cjs |
|------|---------------|------|------|
| `"module"` | ESM | ESM | CommonJS |
| `"commonjs"` | CommonJS | ESM | CommonJS |

```javascript
// package.json: { "type": "module" }
// index.js
export const foo = 'bar';  // ✓ Valid ESM

// package.json: { "type": "commonjs" } or missing
// index.js
module.exports = { foo: 'bar' };  // ✓ Valid CJS
```

## Hands-On Exercise 1: Dependency Audit

Analyze a package.json and identify issues:

```json
{
  "dependencies": {
    "react": "^18.0.0",
    "typescript": "^5.0.0",
    "axios": "^1.0.0"
  },
  "devDependencies": {
    "@types/react": "^18.0.0",
    "jest": "^29.0.0",
    "dotenv": "^16.0.0"
  }
}
```

<details>
<summary>Solution</summary>

**Issues**:
1. ❌ `typescript` should be in `devDependencies` (build tool)
2. ⚠️ `dotenv` might be needed in `dependencies` if loading env vars at runtime
3. ⚠️ If this is a library (not an app), `react` should be a `peerDependency`

**Fixed**:
```json
{
  "dependencies": {
    "axios": "^1.0.0",
    "dotenv": "^16.0.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "@types/react": "^18.0.0",
    "jest": "^29.0.0",
    "react": "^18.0.0"
  },
  "peerDependencies": {
    "react": "^18.0.0"
  }
}
```

</details>

## Hands-On Exercise 2: Configure Dual Package

Create package.json for a library that supports both ESM and CommonJS.

**Requirements**:
- Build outputs to `dist/`
- Support TypeScript
- Provide both ESM (.mjs) and CJS (.js) builds

<details>
<summary>Solution</summary>

```json
{
  "name": "my-dual-package",
  "version": "1.0.0",
  "type": "module",
  "main": "./dist/index.js",
  "module": "./dist/index.mjs",
  "types": "./dist/index.d.ts",
  "exports": {
    ".": {
      "types": "./dist/index.d.ts",
      "import": "./dist/index.mjs",
      "require": "./dist/index.js",
      "default": "./dist/index.js"
    }
  },
  "files": [
    "dist"
  ],
  "scripts": {
    "build": "tsc && tsc --module esnext --outDir dist/esm && mv dist/esm/index.js dist/index.mjs",
    "prepublishOnly": "npm run build"
  },
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}
```

</details>

## Interview Questions

### Q1: What's the difference between ^ and ~ in package.json?

<details>
<summary>Answer</summary>

- `^1.4.2`: Allows updates that don't change the leftmost non-zero digit
  - Allows: 1.4.3, 1.5.0, 1.99.0
  - Blocks: 2.0.0
- `~1.4.2`: Allows patch updates only
  - Allows: 1.4.3, 1.4.99
  - Blocks: 1.5.0

**Exception**: For 0.x.y, `^0.4.2` behaves like `~0.4.2` (only allows patch updates) because pre-1.0 versions can have breaking changes in minor versions.

</details>

### Q2: Why use peerDependencies instead of dependencies?

<details>
<summary>Answer</summary>

Prevents version conflicts for plugins/extensions. Example:

If your React component library used regular dependencies:
```
App depends on React 18.0.0
Your library bundles React 18.2.0
→ Two React copies in bundle! (Breaks React)
```

With peerDependencies:
```
App depends on React 18.0.0
Your library requires React ^18.0.0 (peer)
→ Single React copy shared
```

</details>

### Q3: What's the purpose of the "exports" field?

<details>
<summary>Answer</summary>

1. **Encapsulation**: Block access to internal modules
2. **Multiple entry points**: Different files for different imports
3. **Conditional exports**: Serve different files for ESM vs CommonJS
4. **Future-proof**: Standard way to define package structure

```json
{
  "exports": {
    ".": "./index.js",              // import 'pkg'
    "./utils": "./utils.js",        // import 'pkg/utils'
    "./internal/*": null            // Block 'pkg/internal/*'
  }
}
```

</details>

### Q4: When does "prepare" script run vs "prepublishOnly"?

<details>
<summary>Answer</summary>

- `prepare`: Runs on:
  - `npm install` (after installing from git)
  - `npm publish`
  - Local `npm install`

- `prepublishOnly`: Runs only on:
  - `npm publish`

**Use case**:
- `prepublishOnly`: Build + test before publishing
- `prepare`: Build (for git dependencies that need compilation)

</details>

## Key Takeaways

1. **Dependencies**: Runtime vs dev vs peer - know the difference
2. **Semver**: `^` for minor updates, `~` for patches, exact for critical deps
3. **Exports**: Modern way to define package entry points (supports dual ESM/CJS)
4. **Files**: Control what gets published to avoid bloat
5. **Lock files**: Always commit them for reproducible builds
6. **Scripts**: Use lifecycle hooks (prepublishOnly, prepare) for automation
7. **Type field**: Controls .js file interpretation (module vs commonjs)

## Next Steps

In [Lesson 02: tsconfig.json Mastery](lesson-02-tsconfig-mastery.md), you'll learn:
- Critical compiler options and their implications
- Module resolution strategies
- Project references for monorepos
- Performance optimization
