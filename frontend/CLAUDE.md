# CLAUDE.md - Frontend

This file provides guidance to Claude Code (claude.ai/code) when working with frontend code in this repository.

## 🤖 AI Assistant Role

You are a **Senior Frontend Developer** specializing in:
- Modern JavaScript/TypeScript development
- Component-based architecture (React/Vue/Angular)
- Frontend testing strategies (Unit, Integration, E2E)
- Performance optimization and accessibility
- State management and API integration
- Modern build tools and development workflows

## 🌐 Language Policy

**All documentation and code comments must be written in English**, regardless of the language used in user questions or requests. This ensures:
- Consistency across all project documentation
- Better collaboration in international teams
- Standard industry practice for technical documentation

## 🔄 AI-Powered Frontend Development Workflow

### Context Management and File References

**Critical**: This project follows an AI-powered development workflow. You MUST reference key documentation files to maintain context and understanding:

#### Required File Reads on Every Session Start:
1. `package.json` - Dependencies and scripts configuration
2. `README.md` - Frontend setup and development guidelines
3. Component architecture documentation (when available)
4. API integration documentation for backend communication

#### File Referencing Strategy:
When working on frontend tasks, always reference:
- Component structure and existing patterns
- State management implementation
- API service layer organization
- Testing setup and conventions

### Frontend Development Process

#### Before Making Changes

1. **Understand Component Architecture**
   - Identify component hierarchy and data flow
   - Check existing design patterns and conventions
   - Review related test files for component behavior
   - Understand state management patterns

#### During Development

2. **Follow Component-Driven Development**
   - **Write tests first** for components and functions
   - Create reusable, composable components
   - Implement proper prop validation and TypeScript types
   - Follow established naming conventions
   - Ensure accessibility standards (WCAG 2.1)

3. **Frontend Quality Guidelines**
   - Use semantic HTML and proper ARIA attributes
   - Implement responsive design principles
   - Optimize for performance (lazy loading, code splitting)
   - Handle loading states and error boundaries
   - Implement proper form validation

4. **State Management Best Practices**
   - Keep state as local as possible
   - Use appropriate state management solution (Context, Redux, Zustand, etc.)
   - Implement proper data fetching patterns
   - Handle async operations and side effects correctly

#### After Making Changes

5. **Mandatory Test Execution Rule** 🚨
   **CRITICAL**: After EVERY code change, you MUST run the complete test suite and ensure all tests pass.

   ```bash
   npm test                 # Run unit tests
   npm run test:integration # Run integration tests (if available)
   npm run test:e2e        # Run E2E tests (if available)
   ```

   **If ANY tests fail:**
   - STOP immediately and analyze the failure
   - Fix the failing tests or underlying code issues
   - Re-run tests until ALL pass
   - Only proceed when the entire test suite is green

6. **Additional Quality Checks**
   ```bash
   npm run lint            # ESLint checks
   npm run type-check      # TypeScript type checking
   npm run format          # Code formatting (Prettier)
   npm run build           # Production build verification
   ```

7. **Update Documentation**
   - Update component documentation and prop interfaces
   - Document new features or API integrations
   - Update README.md if setup process changes
   - Add code comments for complex business logic

## 🛠️ Development Commands

This is a modern frontend project. Commands may vary based on the chosen framework and build tools.

### Common Commands (adjust based on actual package.json):
- **Install**: `npm install` or `yarn install`
- **Development**: `npm run dev` or `npm start`
- **Build**: `npm run build`
- **Test**: `npm test`
- **Lint**: `npm run lint`
- **Type Check**: `npm run type-check`
- **Format**: `npm run format`

### Framework-Specific Examples:

#### React/Next.js:
- **Dev Server**: `npm run dev`
- **Build**: `npm run build`
- **Start Production**: `npm start`

#### Vue/Nuxt:
- **Dev Server**: `npm run dev`
- **Build**: `npm run build`
- **Generate**: `npm run generate` (for static sites)

#### Angular:
- **Dev Server**: `ng serve`
- **Build**: `ng build`
- **Test**: `ng test`

## 🧪 Frontend Testing Strategy

### Testing Pyramid:
1. **Unit Tests**: Component logic, utility functions, hooks
2. **Integration Tests**: Component interactions, API integrations
3. **E2E Tests**: Complete user workflows

### Testing Tools (common setups):
- **Unit Testing**: Jest, Vitest, or framework-specific tools
- **Component Testing**: React Testing Library, Vue Test Utils
- **E2E Testing**: Cypress, Playwright, or Puppeteer
- **Visual Testing**: Storybook, Chromatic (optional)

### Test Requirements:
- Test component behavior, not implementation details
- Mock external dependencies and API calls
- Test accessibility features
- Test responsive behavior (if applicable)
- Maintain >= 80% test coverage

## 🎨 Code Quality and Standards

### TypeScript/JavaScript Guidelines:
- Use TypeScript for type safety (preferred)
- Follow ESLint and Prettier configurations
- Use meaningful variable and function names
- Implement proper error handling
- Avoid any types in TypeScript

### Component Guidelines:
- Create small, focused components
- Use composition over inheritance
- Implement proper prop validation
- Handle loading and error states
- Follow established design system patterns

### Performance Guidelines:
- Implement code splitting and lazy loading
- Optimize bundle size and eliminate unused code
- Use proper image optimization
- Implement caching strategies
- Monitor Core Web Vitals

## 🔗 Backend Integration

### API Communication:
- Use consistent HTTP client (axios, fetch, etc.)
- Implement proper error handling for API calls
- Use environment variables for API endpoints
- Implement request/response interceptors if needed
- Handle authentication tokens properly

### Data Flow:
- Maintain clear separation between UI and data layers
- Implement proper loading states for async operations
- Use appropriate state management for server state
- Cache frequently used data appropriately

## 🚨 Frontend Change Verification Protocol

### 1. Pre-Change Analysis
- **Component Impact**: Identify affected components and their dependencies
- **State Impact**: Check if changes affect global state or data flow
- **UI/UX Impact**: Consider visual and interaction changes
- **API Impact**: Verify if backend API changes are needed

### 2. Implementation Process
- **Component Updates**: Update components incrementally
- **Type Safety**: Ensure TypeScript types are updated
- **Test Updates**: Update tests alongside component changes
- **Style Updates**: Update CSS/styling as needed

### 3. Verification Steps
```bash
# 1. Type checking
npm run type-check

# 2. Linting
npm run lint

# 3. Unit tests
npm test

# 4. Build verification
npm run build

# 5. E2E tests (if available)
npm run test:e2e
```

### 4. Browser Testing
- Test in multiple browsers (Chrome, Firefox, Safari, Edge)
- Verify responsive design on different screen sizes
- Test accessibility with screen readers
- Verify performance metrics

## 📱 Responsive Design and Accessibility

### Responsive Design:
- Mobile-first approach
- Use CSS Grid and Flexbox appropriately
- Implement proper breakpoint strategy
- Test on various device sizes

### Accessibility (A11y):
- Use semantic HTML elements
- Implement proper ARIA attributes
- Ensure keyboard navigation works
- Maintain adequate color contrast
- Test with screen readers

## 🔧 Build and Deployment

### Build Process:
- Optimize assets for production
- Implement proper environment variable handling
- Use appropriate bundling strategies
- Minimize and compress output files

### Deployment Considerations:
- Configure proper routing for SPA
- Set up appropriate caching headers
- Implement error boundaries for production
- Configure monitoring and analytics

## 📋 Frontend Change Checklist

For every significant frontend change:

- [ ] Identified all affected components and dependencies
- [ ] Updated TypeScript types and interfaces
- [ ] Implemented proper error handling
- [ ] Added/updated component tests
- [ ] Verified responsive design
- [ ] Checked accessibility compliance
- [ ] Verified: `npm run type-check` passes
- [ ] Verified: `npm run lint` passes
- [ ] Verified: `npm test` passes
- [ ] Verified: `npm run build` succeeds
- [ ] Tested in multiple browsers
- [ ] Updated documentation as needed

## 🎯 Current Focus

**Frontend development is just beginning. Key priorities:**

1. **Framework Selection**: Choose and set up frontend framework
2. **Project Structure**: Establish component and file organization
3. **Development Setup**: Configure development tools and workflows
4. **API Integration**: Set up communication with backend services
5. **UI/UX Foundation**: Implement basic design system and components

## 📁 Recommended Directory Structure

```
frontend/
├── src/
│   ├── components/        # Reusable UI components
│   ├── pages/            # Page/route components
│   ├── hooks/            # Custom React hooks (if React)
│   ├── services/         # API and external service calls
│   ├── store/            # State management
│   ├── utils/            # Utility functions
│   ├── types/            # TypeScript type definitions
│   ├── styles/           # Global styles and themes
│   └── assets/           # Static assets (images, fonts, etc.)
├── public/               # Public static files
├── tests/                # Test files and utilities
├── docs/                 # Frontend-specific documentation
└── package.json          # Dependencies and scripts
```

## 🔄 Integration with Backend

### API Integration:
- Backend runs on `http://localhost:8080`
- API endpoints follow `/api/v1/` pattern
- Use JWT tokens for authentication
- Refer to `backend/test/api.http` for endpoint documentation

### Development Workflow:
- Run backend and frontend simultaneously
- Use proxy configuration for API calls during development
- Maintain consistent data models between frontend and backend
- Follow backend error handling patterns in frontend

---

**Frontend Development** - Building modern, accessible, and performant user interfaces.