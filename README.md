# URL-SHORTENER

## ABOUT

A sophisticated URL shortener service built with modern technologies and scalability in mind. While URL shorteners might seem simple, this project serves as a playground for implementing advanced infrastructure, clean architecture, and exploring new technologies.

## Tech Stack & Why I Chose It

-   **Golang:** After hearing praise about Go's simplicity and performance, I had to try it myself. Its concurrency model and standard library are fantastic for building web services.
-   **HTMX & GoTemplates:** After years of building JavaScript-heavy frontends, the simplicity of HTMX is refreshing. It's amazing how much you can achieve with so little code.
-   **MongoDB:** I chose Mongo for its scaling capabilities. Running it in a cluster configuration with replica sets gives us:
    -   A leader instance handling writes (which are less frequent in our case)
    -   Multiple read replicas to handle the more common read operations
    -   Easy horizontal scaling when needed
-   **Redis:** Using for:
    -   Caching shortened URLs to reduce database load
    -   Managing rate limiting (more on that below)
    -   Ensuring our application servers remain stateless for better scaling

## Architecture Overview

This project follows Clean Architecture principles and is designed with scalability and business requirements in mind. I may not have achieved a perfect implementation, but i will make some improvements in free time

### Infrastructure

The application itself can be scaled horizontally, with handling load balancing (for example nginx). This setup allows us to:

-   Add more application instances during high load
-   Perform zero-downtime deployments
-   Handle failover scenarios gracefully

MongoDB Cluster with replica-set ensures we can handle high read loads (which is the primary operation in a URL shortener) while maintaining data consistency.

### URL Shortening Process

1. **User Authentication & Security**

    - Google ReCaptcha v2 integration to prevent automated access
    - IP-based rate limiting (1 URL per 10 minutes)

2. **URL Generation Strategy**

    - Uses SnowflakeID converted to Base62 for generating short URLs
    - Chosen for:

        - Optimal length of generated URLs
        - High performance compared to cryptographic functions
        - Better scalability than auto-increment
        - More efficient than recursive random ID generation

3. **Data Flow**

    - New URLs are permanently stored in MongoDB
    - Generated short URLs are cached in Redis

### Redirection Flow

1. Check Redis cache first
2. On cache miss, query MongoDB
3. Cache the result in Redis
4. Perform a 301 redirect (enables browser-side caching)

## Running locally

1. **Clone the repo**

```bash
git clone https://github.com/SekulDev/url-shortener.git
```

2. **Set up your environment variables**

    - Copy `./config/.env.dev.example` to `./config/.env.dev`
    - Fill in your configuration values

3. **Run**

```bash
docker compose -f ./docker-compose.dev.yaml --env-file ./config/.env.dev up
```

## Project status

**⚠️ Note:** As the only one developer, I opted to commit directly to main, as implementing a more complex branching strategy wouldn't provide significant benefits for a project of this size.

In the project there are some areas for improvement:

-   Missing integrational / E2E tests, some unit tests are need improvement
-   Code refactoring in specific areas
-   Documentation improvements

Feedback is welcome through GitHub Issues!

## Learning Outcomes

This project served as an excellent learning experience for:

-   Working with Clean Architecture
-   Implementing scalable infrastructure
-   Building business-oriented solutions
-   Web services in Go

These technologies and methodologies have proven valuable and will definitely be utilized in future projects.
