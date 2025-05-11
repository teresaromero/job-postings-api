
# Job Postings API

The Job Posting API is a backend service which enables employers to create, update or delete job postings, while allowing to employees to search and filter on the existent jobs. 


## Features

- CRUD Operations for Jobs: Create, Update and Delete require Authentication and Authorization
- Current Filters: `min_salary`, `max_salary`, `title`, `company`, `location`
- Sorting rule priority customization. 
  Default sorting:
    1. Posts created in the last 7 days rank first above older posts
    2. Posts with higher salaries rank first above posts with lower salaries
    3. Posts made by companies with more open job posts rank first than those for companies with less posts

API Documentation available [here](/openapi.yaml)

## Tech Stack, Architecture and Design Decisions

### Tech Stack

- Go & Gin Framework
- Unit testing: go test
- Lint: go vet

### Decisions

- Configuration Loading: All runtime settings (port, JWT secret, sort-rule order, seed flag) come from a central config.Load() call, decoupling config from code.
- Modularity through clear package boundaries (config, storage, repository, handlers, middleware, sorting).
- Extensibility via injection of Sorter rules and pluggable storage.
- Maintainability using versioned routing, centralized config, custom validators, and a repository layer that makes unit testing and future storage swaps trivial.

- Server default middleware settings:
  - Logging & Recovery middleware for robust request tracing and crash protection.
  - Trusted-proxy disabled by default, letting you control deployment topology explicitly. 
  - Custom validation binding for JobType Enum

- Sorter Component: Constructs a Sorter using the configured rule order (e.g. “recent first, then salary, then company count”), so you can reorder or extend ranking rules without touching handler logic.

- Storage Layer: encapsulates in-memory data operations (CRUD, filtering, sorting) behind a clean interface—ready to swap in a persistent store later.

- Repository Pattern: Wraps storage in a repository that exposes domain-specific methods, isolating business logic from raw data mechanics.

- API Versioning: All routes under /api/v1, making future non-breaking or breaking upgrades straightforward.

- Authentication & Authorization:
  - JWT-based middleware on protected routes (/posts POST/PUT/DELETE).
  - Role checks (“Employer” only)
  - This implementation is just for DEMO purposes. Needs to be implemented with proper password verification and user handling.

- Healthcheck: A lightweight /ping route for uptime monitoring.

- Token Issuance: POST /auth/token separates auth concerns from business logic, centralizing JWT secret use.

## Environment Variables

- JWT_SECRET: Secret for signing JWT. String (required).
- SEED_DATA: Enable the initial seed of the data. Bool (optional).

## Run Locally

Clone the project

```bash
  git clone https://github.com/teresaromero/job-postings-api.git
```

Use the Makefile commands to run locally:

```bash
  make <command>
```

- `build`: builds the binary
- `dev`: runs the `main.go` file
- `run`: builds the binary and runs it
- `clean`: remove the binary
- `test`: run unit tests
- `lint`: run code lint and openapi validation. [Vacuum](https://github.com/daveshanley/vacuum) is required to be locally installed.


## License

[MIT](https://choosealicense.com/licenses/mit/)

