
# Job Postings API

The Job Posting API is an MVP for a backend service which enables employers to create, update or delete job postings, while allowing to employees to search and filter on the existent jobs. 


## Features

- CRUD Operations for Jobs: Create, Update and Delete require Authentication and Authorization
- Current Filters: `min_salary`, `max_salary`, `title`, `company`, `location`
- Sorting rule priority customization. 
  Default sorting:
    1. Posts created in the last 7 days rank first above older posts
    2. Posts with higher salaries rank first above posts with lower salaries
    3. Posts made by companies with more open job posts rank first than those for companies with less posts

  Current rules map:
    "0": sortByCompanyPostsCount
		"1": sortByHighSalary
		"2": sortByLastSevenDays

API Documentation available [here](/openapi.yaml)

## Tech Stack, Architecture and Design Decisions

### Tech Stack

- Go & Gin Framework
- Unit & Integration testing: go test
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

### Future work

This MVP is running on in-memory structures. This provides an easy interface to validate the sorting methods but is not production ready.
In order to evolve and be able to scale the following should be achieved:

- Pagination. There is no pagination in the current implementation. As the user would be potentially scrolling looking into the job postings, a cursor pagination would be a good fit for the list handler.

- Use a distributed database to store the job postings. Having data decoupled from the api will allow to scale the application horizontally and share the same data.

- The custom sorting feature as implemented runs over the data for every rule, creating a final list with the combination of filters by priority. This increases also the overhead while data increases. This can take some time of processing, which can block the request. An improvement would be to implement a worker to perform the sorting in async and provide the client with an entrypoint to check the status of the process, then read the results. This results can be read paginated also, which will also improve the performance.

- A cache to serve faster queries that have already been requested.


## Environment Variables

- JWT_SECRET: Secret for signing JWT. String (required).
- SEED_DATA: Enable the initial seed with random fake data. Bool (optional).
- SEED_FILE: Provide a JSON file with data to init the service. String (optional).
- NOW: Date with format `2025-05-11T00:00:00Z` to mock current day. String (optional).
- SORT_RULES_ORDER: Sorting rules al mapped. Default is ["0", "1", "2"] 

## Run Locally and Development

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
- `integration-test`: run integration tests
- `lint`: run code lint and openapi validation. [Vacuum](https://github.com/daveshanley/vacuum) is required to be locally installed.


### Dockerfile

Provided Dockerfile and docker-compose to run the api inside a container.


## License

[MIT](https://choosealicense.com/licenses/mit/)

