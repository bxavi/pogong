# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.16.0
import dataclasses
import datetime


@dataclasses.dataclass()
class Account:
    id: int
    email: str
    password: str
    created_at: datetime.datetime
