package authz

import future.keywords.contains
import future.keywords.if
import future.keywords.in

default allow := false

# Allow health check endpoint for everyone
allow if {
    input.path == "/health"
    input.method == "GET"
}

# Allow user registration without authentication
allow if {
    input.path == "/api/v1/auth/register"
    input.method == "POST"
}

# Allow login without authentication
allow if {
    input.path == "/api/v1/auth/login"
    input.method == "POST"
}

# Authenticated users can access their own profile
allow if {
    input.method == "GET"
    input.path == sprintf("/api/v1/users/%s", [input.user.id])
    input.user.id != ""
}

# Authenticated users can update their own profile
allow if {
    input.method == "PUT"
    input.path == sprintf("/api/v1/users/%s", [input.user.id])
    input.user.id != ""
}

# Admin users can access all user endpoints
allow if {
    input.path_prefix == "/api/v1/users"
    "admin" in input.user.roles
}

# Admin users can trigger workflows
allow if {
    input.path_prefix == "/api/v1/workflows"
    "admin" in input.user.roles
}

# Users with workflow_executor role can trigger specific workflows
allow if {
    input.path == "/api/v1/workflows/user-onboarding"
    input.method == "POST"
    "workflow_executor" in input.user.roles
}

# Rate limiting rules
rate_limit := 100 if {
    "premium" in input.user.roles
} else := 10

# Resource access rules
resource_access[resource] if {
    some role in input.user.roles
    some resource in data.roles[role].resources
}

# Define role permissions
roles := {
    "admin": {
        "resources": ["users", "workflows", "policies", "logs"],
        "actions": ["create", "read", "update", "delete"]
    },
    "user": {
        "resources": ["profile"],
        "actions": ["read", "update"]
    },
    "workflow_executor": {
        "resources": ["workflows"],
        "actions": ["execute"]
    },
    "premium": {
        "resources": ["profile", "analytics"],
        "actions": ["read", "update"]
    }
}