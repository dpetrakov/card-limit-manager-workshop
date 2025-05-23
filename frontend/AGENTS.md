# Frontend Agent Rules

## Development Guidelines

- **Framework**: React 18 with Vite and TypeScript.
- **Styling**: CSS Modules or standard CSS. (Decision on Tailwind CSS/Emotion can be made later if needed).
- **State Management**: React Context API or local component state for initial tasks. (Decision on Redux/Zustand can be made later if needed).
- **API Calls**: Use `fetch` API or `axios`. Ensure proper error handling and request/response logging during development. All API calls should go through the `/api/v1/` path, which is proxied by Nginx to the backend.

## Code Quality

- **Linting**: ESLint. Ensure an ESLint configuration (e.g., `.eslintrc.json` or `.eslintrc.js`) is present in the `frontend/` directory, configured with recommended TypeScript and React rules.
- **Formatting**: Prettier. Ensure a Prettier configuration (e.g., `.prettierrc.json`) is present in the `frontend/` directory. Consider integrating with ESLint.
- **Comments**: Use comments for complex or non-obvious logic.

## Testing

- **Framework**: Jest and React Testing Library.
- **Configuration**: Ensure Jest is configured to work with Vite and TypeScript (e.g., via `jest.config.js` or in `package.json`).
- **Coverage**: Aim for >= 80% unit test coverage for new components and logic.
- **Test Execution**: Tests should be runnable with `npm test` or `yarn test`.

## Project Structure

- Organize components, services, hooks, and types into appropriate subdirectories within `frontend/src/`. For example:
  - `frontend/src/components/`
  - `frontend/src/pages/` (for top-level page components)
  - `frontend/src/services/` (for API call modules)
  - `frontend/src/hooks/`
  - `frontend/src/types/`

## CI/CD

- Linters and tests must pass before code is merged.
