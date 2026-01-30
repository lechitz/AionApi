# Diagram Syntax Guide

This folder uses sequencediagram.org text syntax.

## Example

title Example Flow

participant "User" as U
participant "Service" as S

U->S: Call endpoint
S->U: Response

## Notes

- Use `title` for the diagram title.
- Define participants with readable names.
- Keep flows short and focused.
