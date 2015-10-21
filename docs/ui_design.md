# Helm UI Design Document

This document outlines a proposed architecture for a Helm UI MVP.

Architectural design goals:

- Outline a definited set of feature requirements that will ship with the MVP.
- Publish an agreed upon UI tech stack that is appropriate for the requirements of the MVP effort, ongoing feature development, and maintenance.
- Ensure consistency between UI consumption and API(s) delivery.

## Feature Requirements

We will confine our initial UI implementation to fulfill the following high level goals: 

- Introduction to **Helm: The Kubernetes Package Manager**, including concepts and terms 
- Simple *Getting Started* instructions for installation, usage, and best practices
- Simple *How to Contribute / Author Your Own Helm Charts* instructions

As the above requires no API interaction with a backend service, it is reasonable to suggest a simple static website implementation. However, we should consider the pros/cons of a static implementation vs. an API service-aware, easily extendable implementations.

### Static Website Implementation

#### Pros

- easiest bootstrapping approach: zero to presentable with minimal tech effort
- easiest to maintain, i.e. designers, CMS-capable marketing folks

#### CONS

- not easy to "upgrade" to a CRUD-aware front end solution
- higher likelihood of being separately maintained over time, as distinct from an eventual CRUD-aware front end, a net addition of overall tech maintenance burden
- does not create any tech equity toward the bootstraping of a future CRUD-aware front end

### Dynamic (single-page CRUD app) Website Implementation

#### Pros

- can fulfill MVP requirements and act as a static website
- enables an agile bridge for future CRUD-aware front end features

#### Cons
- not the easiest bootstrapping approach, requires thinking about future CRUD features and spending some time to choose the right tool for that future job
- requires engineering TLC over time: it's a dev project, not a CMS instance
