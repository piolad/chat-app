# Real-Time Chat Application

## Overview

This real-time chat application enables users to communicate with each other in real-time. It's built on a microservices architecture using Docker, with services written in Python, C#, and C++.

<img src="./docs/general-design.png" alt="General design" width="70%"  align="center" />

## Features

- User authentication and registration.
- Real-time text messaging.
- Information about weather in the user's location.

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation

1. Clone the repository to your local machine:

```bash
git clone https://github.com/piolad/chat-app.git
cd real-time-chat-app
```

2. Run the postrgres database in docker:

```bash
docker network create auth-network
docker run --net auth-network --rm -p 5432:5432 --name auth-service-db -e POSTGRES_PASSWORD=mysecretpassword -d postgres
docker build -t auth-service:0.6 .\auth-service\
docker run --rm  --net auth-network --name auth-service auth-service:0.6
```