package authz.data

import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Filter user data based on roles
filtered_user_fields[field] if {
    field := "id"
}

filtered_user_fields[field] if {
    field := "username"
}

filtered_user_fields[field] if {
    field := "email"
    check_email_access
}

filtered_user_fields[field] if {
    field := "first_name"
}

filtered_user_fields[field] if {
    field := "last_name"
}

filtered_user_fields[field] if {
    field := "created_at"
    "admin" in input.user.roles
}

filtered_user_fields[field] if {
    field := "updated_at"
    "admin" in input.user.roles
}

filtered_user_fields[field] if {
    field := "is_active"
    "admin" in input.user.roles
}

# Check if user can see email addresses
check_email_access if {
    "admin" in input.user.roles
}

check_email_access if {
    input.user.id == input.resource.user_id
}

# Data transformation rules
transform_user(user) := transformed if {
    transformed := {
        field: user[field] |
        field in filtered_user_fields
    }
}

# Workflow visibility rules
visible_workflows[workflow] if {
    workflow := "user-onboarding"
    "admin" in input.user.roles
}

visible_workflows[workflow] if {
    workflow := "user-onboarding"
    "workflow_executor" in input.user.roles
}

# Audit log access
can_view_audit_logs if {
    "admin" in input.user.roles
}

can_view_audit_logs if {
    "auditor" in input.user.roles
}