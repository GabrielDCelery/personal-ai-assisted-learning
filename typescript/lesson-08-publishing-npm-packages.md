# Lesson 08: Publishing npm Packages

Complete guide to building, versioning, and publishing TypeScript packages to npm.

## Package Setup

### Project Structure

```
my-package/
  src/
    index.ts
    utils.ts
  dist/           ← Generated
    index.js
    index.mjs
    index.d.ts
  package.json
  tsconfig.json
  .npmignore
  README.md
  LICENSE
```

### package.json Configuration

```json
{
  "name": "@scope/my-package",
  "version": "1.0.0",
  "description": "My awesome TypeScript package",
  "keywords": ["typescript", "utility"],
  "author": "Your Name <email@example.com>",
  "license": "MIT",
  "repository": {
    "type": "git",
    "url": "https://github.com/user/my-package.git"
  },
  "bugs": {
    "url": "https://github.com/user/my-package/issues"
  },
  "homepage": "https://github.com/user/my-package#readme",

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
    },
    "./package.json": "./package.json"
  },

  "files": [
    "dist",
    "README.md",
    "LICENSE"
  ],

  "scripts": {
    "build": "npm run build:cjs && npm run build:esm && npm run build:types",
    "build:cjs": "tsc --module commonjs --outDir dist",
    "build:esm": "tsc --module esnext --outDir dist/esm && node scripts/rename-esm.js",
    "build:types": "tsc --declaration --emitDeclarationOnly --outDir dist",
    "clean": "rm -rf dist",
    "prepublishOnly": "npm run clean && npm run build && npm test",
    "test": "jest"
  },

  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0",
    "jest": "^29.0.0"
  }
}
```

### tsconfig.json

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "lib": ["ES2020"],

    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,

    "outDir": "./dist",
    "rootDir": "./src",

    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,

    "moduleResolution": "node"
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
```

## Dual Package Publishing (ESM + CJS)

### Approach 1: Separate Output Dirs

```json
{
  "exports": {
    ".": {
      "import": "./dist/esm/index.js",
      "require": "./dist/cjs/index.js",
      "types": "./dist/types/index.d.ts"
    }
  }
}
```

**Build scripts**:
```bash
tsc --module commonjs --outDir dist/cjs
tsc --module esnext --outDir dist/esm
tsc --declaration --emitDeclarationOnly --outDir dist/types
```

### Approach 2: File Extension (.mjs)

```json
{
  "exports": {
    ".": {
      "import": "./dist/index.mjs",
      "require": "./dist/index.js",
      "types": "./dist/index.d.ts"
    }
  }
}
```

**Build script** (rename-esm.js):
```javascript
import fs from 'fs';
import path from 'path';

const esmDir = path.join(process.cwd(), 'dist/esm');
const files = fs.readdirSync(esmDir);

for (const file of files) {
  if (file.endsWith('.js')) {
    const oldPath = path.join(esmDir, file);
    const newPath = path.join(process.cwd(), 'dist', file.replace('.js', '.mjs'));
    fs.renameSync(oldPath, newPath);
  }
}
```

### Approach 3: Use Build Tools

**Using tsup**:
```bash
npm install -D tsup
```

```javascript
// tsup.config.ts
import { defineConfig } from 'tsup';

export default defineConfig({
  entry: ['src/index.ts'],
  format: ['cjs', 'esm'],
  dts: true,
  splitting: false,
  sourcemap: true,
  clean: true,
});
```

```json
{
  "scripts": {
    "build": "tsup"
  }
}
```

**Using @microsoft/api-extractor** (for API documentation):
```bash
npm install -D @microsoft/api-extractor
```

## Version Management

### Semantic Versioning (semver)

```
MAJOR.MINOR.PATCH

1.2.3
│ │ └─ Patch: Bug fixes
│ └─── Minor: New features (backward compatible)
└───── Major: Breaking changes
```

### npm version Command

```bash
# Patch (1.0.0 → 1.0.1)
npm version patch

# Minor (1.0.0 → 1.1.0)
npm version minor

# Major (1.0.0 → 2.0.0)
npm version major

# Prerelease (1.0.0 → 1.0.1-0)
npm version prerelease

# Specific version
npm version 2.0.0

# Create git tag
npm version patch -m "Release %s"
```

### Version Hooks

```json
{
  "scripts": {
    "preversion": "npm test",
    "version": "npm run build && git add -A dist",
    "postversion": "git push && git push --tags"
  }
}
```

**Workflow**:
1. `npm version patch`
2. Runs `preversion` (tests)
3. Updates version in package.json
4. Runs `version` (builds, stages files)
5. Creates git commit and tag
6. Runs `postversion` (pushes to git)

### Pre-release Versions

```bash
# Alpha
npm version prerelease --preid=alpha
# 1.0.0 → 1.0.1-alpha.0

# Beta
npm version prerelease --preid=beta
# 1.0.0 → 1.0.1-beta.0

# RC (Release Candidate)
npm version prerelease --preid=rc
# 1.0.0 → 1.0.1-rc.0
```

### Publishing Pre-releases

```bash
# Publish with tag (not 'latest')
npm publish --tag beta

# Install specific tag
npm install my-package@beta
```

## Publishing Workflow

### Initial Setup

```bash
# Login to npm
npm login

# For scoped packages (@scope/package)
npm login --scope=@scope
```

### First Publish

```bash
# Public package
npm publish

# Scoped public package (default is private)
npm publish --access public

# Check what would be published
npm pack --dry-run
```

### .npmignore

Controls what gets published (similar to .gitignore).

```
# .npmignore
src/
tests/
*.test.ts
.github/
.vscode/
tsconfig.json
jest.config.js
.eslintrc.js
```

**Note**: If .npmignore doesn't exist, npm uses .gitignore.

**Always include**: dist/, README.md, LICENSE

### Publishing Checklist

```bash
# 1. Run tests
npm test

# 2. Build
npm run build

# 3. Check package contents
npm pack --dry-run

# 4. Test in local project
npm pack
# Creates my-package-1.0.0.tgz
cd ../test-project
npm install ../my-package/my-package-1.0.0.tgz

# 5. Check types work
# In test project, import and verify types

# 6. Bump version
npm version patch

# 7. Publish
npm publish

# 8. Verify on npm
npm view my-package
```

## Automated Publishing

### GitHub Actions

```yaml
# .github/workflows/publish.yml
name: Publish to npm

on:
  release:
    types: [created]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          registry-url: 'https://registry.npmjs.org'

      - name: Install dependencies
        run: npm ci

      - name: Run tests
        run: npm test

      - name: Build
        run: npm run build

      - name: Publish
        run: npm publish --access public
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
```

### Semantic Release

Automates versioning and publishing based on commit messages.

```bash
npm install -D semantic-release @semantic-release/git @semantic-release/changelog
```

```javascript
// .releaserc.js
module.exports = {
  branches: ['main'],
  plugins: [
    '@semantic-release/commit-analyzer',
    '@semantic-release/release-notes-generator',
    '@semantic-release/changelog',
    '@semantic-release/npm',
    '@semantic-release/git',
    '@semantic-release/github',
  ],
};
```

**Commit conventions**:
```bash
feat: add new feature      # → Minor release
fix: resolve bug           # → Patch release
BREAKING CHANGE: ...       # → Major release
```

## Package Distribution

### CDN Distribution

Packages automatically available on CDNs:

**unpkg**:
```html
<script src="https://unpkg.com/my-package@1.0.0/dist/index.js"></script>
```

**jsDelivr**:
```html
<script src="https://cdn.jsdelivr.net/npm/my-package@1.0.0/dist/index.js"></script>
```

### Private Packages

```bash
# Publish to private registry
npm publish --registry https://npm.pkg.github.com

# Or configure in package.json
{
  "publishConfig": {
    "registry": "https://npm.pkg.github.com"
  }
}
```

### Scoped Packages

```json
{
  "name": "@myorg/my-package",
  "publishConfig": {
    "access": "public"  // or "restricted" for private
  }
}
```

## Package Maintenance

### Deprecating Versions

```bash
# Deprecate specific version
npm deprecate my-package@1.0.0 "Critical bug, use 1.0.1+"

# Deprecate all versions
npm deprecate my-package "Package no longer maintained"
```

### Unpublishing

```bash
# Can only unpublish within 72 hours
npm unpublish my-package@1.0.0

# Unpublish entire package (not recommended)
npm unpublish my-package --force
```

### Package Transfer

```bash
# Add collaborator
npm owner add <username> my-package

# Remove collaborator
npm owner rm <username> my-package

# List owners
npm owner ls my-package
```

## Best Practices

### 1. Version Constraints for Dependencies

```json
{
  "dependencies": {
    "lodash": "^4.17.21"      // ✓ Allow minor/patch
  },
  "devDependencies": {
    "typescript": "^5.0.0"    // ✓ Allow minor/patch
  },
  "peerDependencies": {
    "react": "^18.0.0"        // ✓ Wide range for compatibility
  }
}
```

### 2. Include Essential Files Only

```json
{
  "files": [
    "dist",
    "README.md",
    "LICENSE"
  ]
}
```

**Don't include**: src/, tests/, config files

### 3. Provide Examples

```markdown
# README.md

## Installation

\`\`\`bash
npm install my-package
\`\`\`

## Usage

\`\`\`typescript
import { myFunction } from 'my-package';

const result = myFunction('input');
\`\`\`
```

### 4. Semantic Versioning

- Patch: Bug fixes only
- Minor: New features, backward compatible
- Major: Breaking changes

### 5. Test Before Publishing

```bash
# Always run:
npm run build
npm test
npm pack --dry-run
```

## Hands-On Exercise: Publish a Package

Create and publish a simple utility package:

**Requirements**:
1. Create a package with one utility function
2. Support both ESM and CommonJS
3. Include TypeScript types
4. Write README with examples
5. Publish to npm (or test with `npm pack`)

<details>
<summary>Solution</summary>

```typescript
// src/index.ts
/**
 * Capitalizes the first letter of a string
 * @param str - Input string
 * @returns Capitalized string
 */
export function capitalize(str: string): string {
  if (!str) return str;
  return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Debounces a function
 * @param fn - Function to debounce
 * @param delay - Delay in milliseconds
 * @returns Debounced function
 */
export function debounce<T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout>;

  return function (...args: Parameters<T>) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn(...args), delay);
  };
}
```

```json
// package.json
{
  "name": "@yourusername/utils",
  "version": "1.0.0",
  "description": "Simple utility functions",
  "main": "./dist/index.js",
  "module": "./dist/index.mjs",
  "types": "./dist/index.d.ts",
  "exports": {
    ".": {
      "types": "./dist/index.d.ts",
      "import": "./dist/index.mjs",
      "require": "./dist/index.js"
    }
  },
  "files": ["dist"],
  "scripts": {
    "build": "tsc && tsc --module esnext --outDir dist/esm",
    "prepublishOnly": "npm run build"
  },
  "keywords": ["utils", "typescript"],
  "license": "MIT",
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}
```

```markdown
# @yourusername/utils

Simple utility functions in TypeScript.

## Installation

\`\`\`bash
npm install @yourusername/utils
\`\`\`

## Usage

\`\`\`typescript
import { capitalize, debounce } from '@yourusername/utils';

console.log(capitalize('hello')); // 'Hello'

const search = debounce((query) => {
  console.log('Searching:', query);
}, 300);
\`\`\`

## License

MIT
```

```bash
# Build and test locally
npm run build
npm pack

# Test in another project
cd ../test-project
npm install ../utils/yourusername-utils-1.0.0.tgz

# Publish
npm publish --access public
```

</details>

## Interview Questions

### Q1: How to publish dual ESM/CJS package?

<details>
<summary>Answer</summary>

**Method 1: Different extensions**
```json
{
  "main": "./dist/index.js",      // CJS
  "module": "./dist/index.mjs",   // ESM
  "types": "./dist/index.d.ts",
  "exports": {
    ".": {
      "import": "./dist/index.mjs",
      "require": "./dist/index.js",
      "types": "./dist/index.d.ts"
    }
  }
}
```

**Build**: Compile twice, rename ESM files to .mjs

**Method 2: Separate directories**
```json
{
  "exports": {
    ".": {
      "import": "./dist/esm/index.js",
      "require": "./dist/cjs/index.js",
      "types": "./dist/types/index.d.ts"
    }
  }
}
```

</details>

### Q2: What's the difference between dependencies and peerDependencies?

<details>
<summary>Answer</summary>

**dependencies**: Installed with your package
```json
{
  "dependencies": {
    "lodash": "^4.17.21"  // Always installed
  }
}
```

**peerDependencies**: User must install (not bundled)
```json
{
  "peerDependencies": {
    "react": "^18.0.0"  // User's project must have React
  }
}
```

**Use peerDependencies** for:
- Plugins/extensions (React components, Babel plugins)
- Avoid version conflicts
- Prevent bundling multiple copies

</details>

### Q3: How does npm version work?

<details>
<summary>Answer</summary>

```bash
npm version patch  # 1.0.0 → 1.0.1
```

**Process**:
1. Runs `preversion` script
2. Updates package.json version
3. Runs `version` script
4. Commits changes
5. Creates git tag
6. Runs `postversion` script

**Automation**:
```json
{
  "scripts": {
    "preversion": "npm test",
    "version": "npm run build && git add dist",
    "postversion": "git push && git push --tags && npm publish"
  }
}
```

</details>

### Q4: What files should be published?

<details>
<summary>Answer</summary>

**Include**:
- `dist/` (compiled code)
- `README.md`
- `LICENSE`
- `CHANGELOG.md`

**Exclude** (via .npmignore):
- `src/` (source TypeScript)
- `tests/`
- Config files (tsconfig.json, .eslintrc, etc.)
- `.github/`, `.vscode/`

```json
{
  "files": ["dist", "README.md", "LICENSE"]
}
```

Check before publishing:
```bash
npm pack --dry-run
```

</details>

## Key Takeaways

1. **Dual Publishing**: Support both ESM and CJS with exports field
2. **Types**: Always include .d.ts files (declaration: true)
3. **Versioning**: Follow semver, use npm version commands
4. **Files**: Only publish dist/, exclude source and config
5. **Testing**: Test locally with npm pack before publishing
6. **Automation**: Use GitHub Actions or semantic-release
7. **Documentation**: Clear README with installation and usage examples

## Summary

You now have a comprehensive understanding of:
- Package configuration for publishing
- Dual ESM/CJS distribution
- Versioning and semver
- Publishing workflow and automation
- Package maintenance and best practices

This completes the TypeScript interview prep series. You're now equipped with the knowledge to confidently discuss TypeScript ecosystem, tooling, type system, and publishing practices in interviews.
