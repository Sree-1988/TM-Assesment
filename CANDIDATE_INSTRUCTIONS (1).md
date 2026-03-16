# Interview Exercise - Candidate Instructions

## Overview

Welcome to the technical assessment for the Software Engineer position.
This exercise simulates a real code review scenario you'd encounter on the team.

## The Task

A team member has submitted MR !1 to implement the Service Tags feature described in ticket SR-142 (below).
Your task is to review the merge request and provide feedback.

**Branch:** `tags` → `main`

## Ticket: SR-142

```
Title: Add tagging support to Service Registry API

Type: Feature
Priority: Medium
Story Points: 5

Description:
We need the ability to tag services with arbitrary labels for better organization
and filtering. This will allow teams to group services by environment, ownership,
technology stack, or any other custom categorization.

Acceptance Criteria:
- Services can have multiple tags (array of strings)
- Tags should be validated: non-empty, alphanumeric with hyphens, max 50 chars
- Add ability to filter services by tag via query parameter
- Update API documentation to reflect new functionality
- Include tests for tag validation and filtering
- Ensure backward compatibility with existing services (tags optional)

Technical Notes:
- Consider tag storage efficiency for future querying needs
- Follow existing patterns in the codebase
- Update relevant documentation (README, API specs)
```

## Getting Started

1. **Request access**: Send your GitLab.com username to your recruitment contact
   - Create an account at https://gitlab.com if needed

2. **Once you have access**: visit the repository at(https://gitlab.com/tm_techops/software-interview) and follow the README for setup

3. **Review the MR**: Compare the implementation(https://gitlab.com/tm_techops/software-interview/-/merge_requests/1) against the ticket requirements

## What We're Looking For

Review the MR as you would any team member's work:

- Does it fulfill the ticket requirements?
- Is the implementation sound?
- Are there bugs or potential issues?
- Is it production-ready?

## How to Submit Your Review

**You have write access to the repository.** Conduct your review as you would on your team:

- **Leave comments on the MR** - identify issues, ask questions, suggest improvements
- **Make direct changes if you choose** - push commits to the `tags` branch to fix critical issues
- **Think strategically** - beyond this MR, what future work would benefit the codebase?
  - Example: "We should consider adding request tracing to improve debuggability across services"
  - These aren't requirements for *this* MR, but valuable suggestions for the team's roadmap

## Important: Manage Your Time

**Prioritize breadth over depth.** Focus on:
1. Identifying and commenting on issues across the codebase
2. Reviewing against ticket requirements
3. Catching bugs, design problems, and missing functionality

You *can* make direct code changes, but extensive refactoring isn't expected.
**Comments are more valuable than fixes** - we want to see your ability to identify issues comprehensively.

Target: Complete your review within 1-2 hours.

## Timeline

Please submit your review within **48 hours** of receiving repository access.

## Questions?

If you have questions about the code being reviewed or the exercise, tag **@mrasnak3** (Michael Rasnake) or **@max-tm** (Max) directly on the MR.

If you have issues accessing the assigment, contact your recruitment representative.