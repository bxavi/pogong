# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.16.0
import dataclasses


@dataclasses.dataclass()
class Accounts:
    id: int
    email: str
    password: str