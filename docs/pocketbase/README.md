PocketBase Site
======================================================================

This is the PocketBase Site (aka. https://pocketbase.io), built with SvelteKit.


## Development and contribution

Download the repo and run the appropriate console commands:

```sh
# install dependencies
npm install

# start a dev server with hot reload at localhost:5173
npm run dev

# or generate production ready bundle
npm run build
```

# PocketBase Documentation

This directory contains a local copy of the PocketBase documentation from https://github.com/pocketbase/site, as of April 2025.

## Purpose
This documentation is maintained as a reference to ensure that:
1. Your implementations using PocketBase are consistent with official documentation
2. You have offline access to PocketBase API and SDK reference material
3. You can check for best practices and examples when implementing PocketBase features

## Structure
The main documentation is located in `src/routes/(app)/docs/`. The content is organized by topic and implementation language:

- Go implementation: folders prefixed with `go-`
- JavaScript implementation: folders prefixed with `js-`
- API documentation: folders prefixed with `api-`
- General topics: other folders

## Additional Resources
For more detailed guidance on using this documentation, see the `.github/copilot-instructions.md` file in the root of this repository.

## Updating
To update this documentation in the future, you can pull the latest version from the official repository:

```
cd /path/to/docs/pocketbase
git pull origin main
```

Then check for any API changes that might affect your implementation.
